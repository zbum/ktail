package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// streamLogsMultiNamespace streams logs from multiple pods across multiple namespaces
func streamLogsMultiNamespace(clientset *kubernetes.Clientset, pods []PodInfo) error {
	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Create channels for each pod's log stream
	logChan := make(chan LogLine, 100)

	// Start goroutines for each pod
	for _, pod := range pods {
		go func(pod PodInfo) {
			// Send header information for this pod
			logChan <- LogLine{
				PodInfo: pod,
				Line:    fmt.Sprintf("=== Starting logs for %s/%s (container: %s) ===", pod.Namespace, pod.Name, pod.Container),
			}

			req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
				Container: pod.Container,
				Follow:    true,
				TailLines: int64Ptr(int64(tailLines)),
			})

			stream, err := req.Stream(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create log stream for %s/%s: %v\n", pod.Namespace, pod.Name, err)
				return
			}
			defer stream.Close()

			scanner := bufio.NewScanner(stream)
			for scanner.Scan() {
				select {
				case <-ctx.Done():
					return
				default:
					logChan <- LogLine{
						PodInfo: pod,
						Line:    scanner.Text(),
					}
				}
			}

			if err := scanner.Err(); err != nil && err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading log stream for %s/%s: %v\n", pod.Namespace, pod.Name, err)
			}
		}(pod)
	}

	// Process log lines from all pods
	for {
		select {
		case <-ctx.Done():
			return nil
		case logLine := <-logChan:
			// Format: [namespace/pod] log line
			fmt.Printf("[%s/%s] %s\n", logLine.PodInfo.Namespace, logLine.PodInfo.Name, logLine.Line)
		}
	}
}

// streamLogsWithWatch streams logs with pod watching capability
func streamLogsWithWatch(clientset *kubernetes.Clientset, initialPods []PodInfo, namespaces []string) error {
	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Create channels for log streaming
	logChan := make(chan LogLine, 100)

	// Track which pods are already being streamed to avoid duplicates
	streamingPods := make(map[string]bool)
	var streamingMutex sync.Mutex

	// Start streaming logs for initial pods
	for _, pod := range initialPods {
		podKey := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
		streamingMutex.Lock()
		streamingPods[podKey] = true
		streamingMutex.Unlock()
		go streamPodLogs(clientset, pod, logChan, ctx)
	}

	// Start pod watchers for each namespace
	for _, ns := range namespaces {
		go watchPodsWithTracking(clientset, ns, logChan, ctx, &streamingPods, &streamingMutex)
	}

	// Process log lines from all pods
	for {
		select {
		case <-ctx.Done():
			return nil
		case logLine := <-logChan:
			// Format: [namespace/pod] log line
			fmt.Printf("[%s/%s] %s\n", logLine.PodInfo.Namespace, logLine.PodInfo.Name, logLine.Line)
		}
	}
}

// watchPodsWithTracking watches for pod changes and tracks streaming status
func watchPodsWithTracking(clientset *kubernetes.Clientset, namespace string, logChan chan<- LogLine, ctx context.Context, streamingPods *map[string]bool, streamingMutex *sync.Mutex) {
	watcher, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create pod watcher for namespace %s: %v\n", namespace, err)
		return
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				fmt.Fprintf(os.Stderr, "Pod watcher channel closed for namespace %s\n", namespace)
				return
			}

			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			podKey := fmt.Sprintf("%s/%s", namespace, pod.Name)

			switch event.Type {
			case "ADDED":
				// New pod created, wait for it to be ready and start streaming its logs
				fmt.Printf("New pod detected: %s/%s, waiting for container to be ready...\n", namespace, pod.Name)
				go waitForPodAndStreamLogsWithTracking(clientset, PodInfo{
					Namespace: namespace,
					Name:      pod.Name,
					Container: getContainerName(pod),
				}, logChan, ctx, streamingPods, streamingMutex)
			case "MODIFIED":
				// Pod status changed, check if it's now ready
				if pod.Status.Phase == corev1.PodRunning {
					// Check if the specific container is ready
					for _, containerStatus := range pod.Status.ContainerStatuses {
						if containerStatus.Name == getContainerName(pod) && containerStatus.Ready {
							streamingMutex.Lock()
							if !(*streamingPods)[podKey] {
								(*streamingPods)[podKey] = true
								streamingMutex.Unlock()
								fmt.Printf("Pod %s/%s is now ready, starting log stream...\n", namespace, pod.Name)
								go streamPodLogs(clientset, PodInfo{
									Namespace: namespace,
									Name:      pod.Name,
									Container: getContainerName(pod),
								}, logChan, ctx)
							} else {
								streamingMutex.Unlock()
							}
							break
						}
					}
				}
			case "DELETED":
				fmt.Printf("Pod deleted: %s/%s, stopping log stream...\n", namespace, pod.Name)
				streamingMutex.Lock()
				delete(*streamingPods, podKey)
				streamingMutex.Unlock()
			}
		}
	}
}

// waitForPodAndStreamLogsWithTracking waits for a pod to be ready and starts streaming its logs
func waitForPodAndStreamLogsWithTracking(clientset *kubernetes.Clientset, pod PodInfo, logChan chan<- LogLine, ctx context.Context, streamingPods *map[string]bool, streamingMutex *sync.Mutex) {
	// Wait for pod to be ready
	maxRetries := 30 // Wait up to 5 minutes (30 * 10 seconds)
	retryCount := 0

	podKey := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)

	for retryCount < maxRetries {
		select {
		case <-ctx.Done():
			return
		default:
			// Check if already streaming
			streamingMutex.Lock()
			if (*streamingPods)[podKey] {
				streamingMutex.Unlock()
				return
			}
			streamingMutex.Unlock()

			// Check pod status
			podObj, err := clientset.CoreV1().Pods(pod.Namespace).Get(ctx, pod.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get pod %s/%s: %v\n", pod.Namespace, pod.Name, err)
				return
			}

			// Check if pod is running and container is ready
			if podObj.Status.Phase == corev1.PodRunning {
				// Check if the specific container is ready
				for _, containerStatus := range podObj.Status.ContainerStatuses {
					if containerStatus.Name == pod.Container && containerStatus.Ready {
						streamingMutex.Lock()
						if !(*streamingPods)[podKey] {
							(*streamingPods)[podKey] = true
							streamingMutex.Unlock()
							fmt.Printf("Pod %s/%s is ready, starting log stream...\n", pod.Namespace, pod.Name)
							streamPodLogs(clientset, pod, logChan, ctx)
						} else {
							streamingMutex.Unlock()
						}
						return
					}
				}
			} else if podObj.Status.Phase == corev1.PodFailed || podObj.Status.Phase == corev1.PodSucceeded {
				fmt.Printf("Pod %s/%s is in %s state, skipping log stream\n", pod.Namespace, pod.Name, podObj.Status.Phase)
				return
			}

			// Wait 10 seconds before retrying
			time.Sleep(10 * time.Second)
			retryCount++
		}
	}

	fmt.Printf("Pod %s/%s did not become ready within timeout, skipping log stream\n", pod.Namespace, pod.Name)
}

// streamPodLogs streams logs from a single pod
func streamPodLogs(clientset *kubernetes.Clientset, pod PodInfo, logChan chan<- LogLine, ctx context.Context) {
	// Send header information for this pod
	logChan <- LogLine{
		PodInfo: pod,
		Line:    fmt.Sprintf("=== Starting logs for %s/%s (container: %s) ===", pod.Namespace, pod.Name, pod.Container),
	}

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		Container: pod.Container,
		Follow:    true,
		TailLines: int64Ptr(int64(tailLines)),
	})

	stream, err := req.Stream(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log stream for %s/%s: %v\n", pod.Namespace, pod.Name, err)
		return
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			logChan <- LogLine{
				PodInfo: pod,
				Line:    scanner.Text(),
			}
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Fprintf(os.Stderr, "Error reading log stream for %s/%s: %v\n", pod.Namespace, pod.Name, err)
	}
}
