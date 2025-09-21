package main

// PodInfo represents information about a Kubernetes pod
type PodInfo struct {
	Namespace string
	Name      string
	Container string
}

// LogLine represents a log line with associated pod information
type LogLine struct {
	PodInfo PodInfo
	Line    string
}
