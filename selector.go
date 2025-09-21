package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// selectNamespace allows interactive selection of a single namespace
func selectNamespace(clientset *kubernetes.Clientset) (string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaceList []string
	for _, ns := range namespaces.Items {
		namespaceList = append(namespaceList, ns.Name)
	}

	return runFuzzyFinder(namespaceList, "Select namespace:")
}

// selectNamespacesMulti allows interactive multi-selection of namespaces
func selectNamespacesMulti(clientset *kubernetes.Clientset) ([]string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %v", err)
	}

	var namespaceList []string
	for _, ns := range namespaces.Items {
		namespaceList = append(namespaceList, ns.Name)
	}

	selected, err := runFuzzyFinderMulti(namespaceList, "Select namespaces (use Tab to select multiple):")
	if err != nil {
		return nil, err
	}

	return selected, nil
}

// selectPod allows interactive selection of a single pod in a namespace
func selectPod(clientset *kubernetes.Clientset, namespace string) (string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pods: %v", err)
	}

	var podList []string
	for _, pod := range pods.Items {
		status := string(pod.Status.Phase)
		if pod.Status.Phase == "Running" {
			status = "🟢 Running"
		} else if pod.Status.Phase == "Pending" {
			status = "🟡 Pending"
		} else if pod.Status.Phase == "Failed" {
			status = "🔴 Failed"
		} else {
			status = "⚪ " + status
		}
		podList = append(podList, fmt.Sprintf("%s\t%s", pod.Name, status))
	}

	// Check if there are any pods in the namespace
	if len(podList) == 0 {
		return "", fmt.Errorf("no pods found in namespace %s", namespace)
	}

	selected, err := runFuzzyFinder(podList, "Select pod:")
	if err != nil {
		return "", err
	}

	// Extract pod name from the selected line
	parts := strings.Split(selected, "\t")
	return parts[0], nil
}

// selectPodsMulti allows interactive multi-selection of pods in a namespace
func selectPodsMulti(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	var podList []string
	for _, pod := range pods.Items {
		status := string(pod.Status.Phase)
		if pod.Status.Phase == "Running" {
			status = "🟢 Running"
		} else if pod.Status.Phase == "Pending" {
			status = "🟡 Pending"
		} else if pod.Status.Phase == "Failed" {
			status = "🔴 Failed"
		} else {
			status = "⚪ " + status
		}
		podList = append(podList, fmt.Sprintf("%s\t%s", pod.Name, status))
	}

	// Check if there are any pods in the namespace
	if len(podList) == 0 {
		fmt.Printf("No pods found in namespace %s, skipping...\n", namespace)
		return []string{}, nil
	}

	selected, err := runFuzzyFinderMulti(podList, "Select pods (use Tab to select multiple):")
	if err != nil {
		return nil, err
	}

	// Extract pod names from the selected lines
	var podNames []string
	for _, line := range selected {
		parts := strings.Split(line, "\t")
		podNames = append(podNames, parts[0])
	}

	return podNames, nil
}

// selectPodsMultiAcrossNamespaces allows interactive multi-selection of pods across multiple namespaces
func selectPodsMultiAcrossNamespaces(clientset *kubernetes.Clientset, namespaces []string) ([]PodInfo, error) {
	var allPods []PodInfo

	// Collect all pods from all namespaces
	for _, ns := range namespaces {
		pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list pods in namespace %s: %v", ns, err)
		}

		for _, pod := range pods.Items {
			status := string(pod.Status.Phase)
			if pod.Status.Phase == "Running" {
				status = "🟢 Running"
			} else if pod.Status.Phase == "Pending" {
				status = "🟡 Pending"
			} else if pod.Status.Phase == "Failed" {
				status = "🔴 Failed"
			} else {
				status = "⚪ " + status
			}

			allPods = append(allPods, PodInfo{
				Namespace: ns,
				Name:      pod.Name,
				Status:    status,
			})
		}
	}

	// Check if there are any pods
	if len(allPods) == 0 {
		return nil, fmt.Errorf("no pods found in any of the selected namespaces")
	}

	// Create display strings for fuzzy finder
	var podList []string
	for _, pod := range allPods {
		podList = append(podList, fmt.Sprintf("%s/%s\t%s", pod.Namespace, pod.Name, pod.Status))
	}

	selected, err := runFuzzyFinderMulti(podList, "Select pods across namespaces (use Tab to select multiple):")
	if err != nil {
		return nil, err
	}

	// Extract pod info from the selected lines
	var selectedPods []PodInfo
	for _, line := range selected {
		parts := strings.Split(line, "\t")
		namespacePod := strings.Split(parts[0], "/")
		if len(namespacePod) == 2 {
			selectedPods = append(selectedPods, PodInfo{
				Namespace: namespacePod[0],
				Name:      namespacePod[1],
			})
		}
	}

	return selectedPods, nil
}

// runFuzzyFinder runs a single-selection fuzzy finder
func runFuzzyFinder(options []string, prompt string) (string, error) {
	idx, err := fuzzyfinder.Find(
		options,
		func(i int) string {
			return options[i]
		},
		fuzzyfinder.WithPromptString(prompt),
	)
	if err != nil {
		return "", fmt.Errorf("fuzzy finder selection cancelled or failed: %v", err)
	}

	return options[idx], nil
}

// runFuzzyFinderMulti runs a multi-selection fuzzy finder
func runFuzzyFinderMulti(options []string, prompt string) ([]string, error) {
	indices, err := fuzzyfinder.FindMulti(
		options,
		func(i int) string {
			return options[i]
		},
		fuzzyfinder.WithPromptString(prompt),
		fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop),
	)
	if err != nil {
		return nil, fmt.Errorf("fuzzy finder multi-selection cancelled or failed: %v", err)
	}

	var selected []string
	for _, idx := range indices {
		selected = append(selected, options[idx])
	}

	return selected, nil
}
