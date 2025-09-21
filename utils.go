package main

import (
	"os"
	"strconv"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
)

// int64Ptr returns a pointer to an int64 value
func int64Ptr(i int64) *int64 { return &i }

// colorize returns colored text if colors are enabled, otherwise returns plain text
func colorize(text, color string) string {
	if noColor {
		return text
	}
	return color + text + ColorReset
}

// colorizeNamespace returns colored namespace text
func colorizeNamespace(namespace string) string {
	return colorize(namespace, ColorGreen)
}

// colorizePod returns colored pod name text
func colorizePod(pod string) string {
	return colorize(pod, ColorGreen)
}

// colorizeContainer returns colored container name text
func colorizeContainer(container string) string {
	return colorize(container, ColorCyan)
}

// parseCustomFlags parses custom flags like -1000f, -500f, etc.
func parseCustomFlags() {
	args := os.Args[1:]
	for i, arg := range args {
		// Check for patterns like -1000f, -500f, etc.
		if len(arg) > 2 && arg[0] == '-' && arg[len(arg)-1] == 'f' {
			// Extract the number part
			numStr := arg[1 : len(arg)-1]
			if num, err := strconv.Atoi(numStr); err == nil && num > 0 {
				// Set tailLines
				tailLines = num
				// Remove this argument from os.Args
				os.Args = append(os.Args[:i+1], os.Args[i+2:]...)
				break // Only process the first matching flag
			}
		}
	}
}
