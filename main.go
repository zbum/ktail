package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	namespace   string
	podName     string
	tailLines   = 100
	multiSelect bool
	container   string
	noColor     bool
)

var rootCmd = &cobra.Command{
	Use:   "ktail",
	Short: "A Kubernetes log tail utility with interactive namespace and pod selection",
	Long: `ktail is a tool that provides tail-like functionality for Kubernetes pod logs.
It allows you to interactively select namespaces and pods using fzf for a better user experience.

Features:
- Single or multi-namespace selection
- Single or multi-pod selection
- Real-time log streaming from multiple pods across multiple namespaces
- Interactive fuzzy search for namespaces and pods
- Support for all pods in selected namespace(s)

Examples:
  ktail                                    # Interactive selection (all pods)
  ktail -n my-namespace                    # All pods in my-namespace
  ktail -m                                 # Multi-select pods
  ktail -n my-ns -p my-pod                 # Specific pod
  ktail -1000f                             # Follow recent 1000 lines
  ktail -500f -n my-ns                     # Follow recent 500 lines from my-ns namespace`,
	Run: runKtail,
}

func init() {
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (if not provided, will be selected interactively)")
	rootCmd.Flags().StringVarP(&podName, "pod", "p", "", "Pod name (if not provided, will select all pods in namespace)")
	rootCmd.Flags().IntVarP(&tailLines, "tail", "t", 10, "Number of lines to show from the end of logs")
	rootCmd.Flags().BoolVarP(&multiSelect, "multi", "m", true, "Enable multi-selection for pods")
	rootCmd.Flags().StringVarP(&container, "container", "c", "", "Container name")
	rootCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}

func main() {
	// Parse custom flags like -1000f before cobra processing
	parseCustomFlags()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runKtail(cmd *cobra.Command, args []string) {
	// Create Kubernetes client
	clientset, err := createK8sClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Determine target namespaces
	var targetNamespace string
	// Single namespace selection
	targetNamespace, err = selectNamespace(clientset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to select namespace: %v\n", err)
		os.Exit(1)
	}

	// Collect pods from all target namespaces
	var allPods []PodInfo

	// Original logic for single namespace or when not using multi-select across namespaces

	var podNames []string
	if podName != "" {
		// Single pod specified
		podNames = []string{podName}
	} else if multiSelect {
		// Multi-select pods in single namespace
		podNames, err = selectPodsMulti(clientset, targetNamespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to select pods in namespace %s: %v\n", targetNamespace, err)
			os.Exit(1)
		}
	} else {
		// Default: Select all pods in the namespace
		podNames, err = getAllPods(clientset, targetNamespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get all pods in namespace %s: %v\n", targetNamespace, err)
			os.Exit(1)
		}
	}

	// Skip if no pods found in this namespace
	if len(podNames) == 0 {
		return
	}

	// Get container names for all selected pods in this namespace
	for _, pod := range podNames {
		var containerName string
		if container == "" {
			containerName, err = getFirstContainer(clientset, targetNamespace, pod)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get container name for pod %s in namespace %s: %v\n", pod, targetNamespace, err)
				os.Exit(1)
			}
		} else {
			containerName = container
		}
		allPods = append(allPods, PodInfo{
			Namespace: targetNamespace,
			Name:      pod,
			Container: containerName,
		})
	}

	if len(allPods) == 0 {
		fmt.Fprintf(os.Stderr, "No pods selected\n")
		os.Exit(1)
	}

	// Display selected pods
	fmt.Printf("Tailing logs for %d pod(s) across %d namespace(s)\n", len(allPods), len(targetNamespace))
	for _, pod := range allPods {
		fmt.Printf("  - %s/%s (container: %s)\n",
			colorizeNamespace(pod.Namespace),
			colorizePod(pod.Name),
			colorizeContainer(pod.Container))
	}
	fmt.Println("Press Ctrl+C to stop...")

	// Stream logs for all selected pods
	err = streamLogsWithWatch(clientset, allPods, targetNamespace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to stream logs: %v\n", err)
		os.Exit(1)
	}
}
