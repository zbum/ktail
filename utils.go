package main

import (
	"os"
	"strconv"
)

// int64Ptr returns a pointer to an int64 value
func int64Ptr(i int64) *int64 { return &i }

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
