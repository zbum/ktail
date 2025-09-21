package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// createK8sClient creates a Kubernetes client using in-cluster config or kubeconfig
func createK8sClient() (*kubernetes.Clientset, error) {
	// Try to use in-cluster config first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubeconfig: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	return clientset, nil
}

// getAllPods retrieves all pod names in a given namespace
func getAllPods(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}

	// Check if there are any pods in the namespace
	if len(podNames) == 0 {
		fmt.Printf("No pods found in namespace %s, skipping...\n", namespace)
	}

	return podNames, nil
}

// getFirstContainer retrieves the first container name from a pod
func getFirstContainer(clientset *kubernetes.Clientset, namespace, podName string) (string, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get pod: %v", err)
	}

	if len(pod.Spec.Containers) == 0 {
		return "", fmt.Errorf("no containers found in pod")
	}

	return pod.Spec.Containers[0].Name, nil
}

// getContainerName extracts the container name from a pod object
func getContainerName(pod *corev1.Pod) string {
	if len(pod.Spec.Containers) == 0 {
		return ""
	}
	return pod.Spec.Containers[0].Name
}
