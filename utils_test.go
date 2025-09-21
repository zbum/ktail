package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseCustomFlags(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedTail int
		expectedArgs []string
		description  string
	}{
		{
			name:         "Valid single digit flag",
			args:         []string{"-5f", "other", "args"},
			expectedTail: 5,
			expectedArgs: []string{"other", "args"},
			description:  "Should parse -5f and set tailLines to 5",
		},
		{
			name:         "Valid multi-digit flag",
			args:         []string{"-1000f", "other", "args"},
			expectedTail: 1000,
			expectedArgs: []string{"other", "args"},
			description:  "Should parse -1000f and set tailLines to 1000",
		},
		{
			name:         "Valid flag at end",
			args:         []string{"other", "args", "-500f"},
			expectedTail: 500,
			expectedArgs: []string{"other", "args"},
			description:  "Should parse -500f at the end and set tailLines to 500",
		},
		{
			name:         "Valid flag in middle",
			args:         []string{"first", "-250f", "last"},
			expectedTail: 250,
			expectedArgs: []string{"first", "last"},
			description:  "Should parse -250f in the middle and set tailLines to 250",
		},
		{
			name:         "Invalid flag with zero",
			args:         []string{"-0f", "other", "args"},
			expectedTail: 100, // default value
			expectedArgs: []string{"-0f", "other", "args"},
			description:  "Should not parse -0f (zero is not positive)",
		},
		{
			name:         "Invalid flag with negative number",
			args:         []string{"--100f", "other", "args"},
			expectedTail: 100, // default value
			expectedArgs: []string{"--100f", "other", "args"},
			description:  "Should not parse --100f (double dash)",
		},
		{
			name:         "Invalid flag with non-numeric",
			args:         []string{"-abcf", "other", "args"},
			expectedTail: 100, // default value
			expectedArgs: []string{"-abcf", "other", "args"},
			description:  "Should not parse -abcf (non-numeric)",
		},
		{
			name:         "Invalid flag too short",
			args:         []string{"-f", "other", "args"},
			expectedTail: 100, // default value
			expectedArgs: []string{"-f", "other", "args"},
			description:  "Should not parse -f (too short)",
		},
		{
			name:         "No custom flags",
			args:         []string{"normal", "args", "here"},
			expectedTail: 100, // default value
			expectedArgs: []string{"normal", "args", "here"},
			description:  "Should not modify args when no custom flags present",
		},
		{
			name:         "Empty args",
			args:         []string{},
			expectedTail: 100, // default value
			expectedArgs: []string{},
			description:  "Should handle empty args",
		},
		{
			name:         "Multiple valid flags - first one wins",
			args:         []string{"-100f", "-200f", "other"},
			expectedTail: 100,
			expectedArgs: []string{"-200f", "other"},
			description:  "Should parse first valid flag and remove it from args",
		},
		{
			name:         "Flag with leading zeros",
			args:         []string{"-007f", "other", "args"},
			expectedTail: 7,
			expectedArgs: []string{"other", "args"},
			description:  "Should parse -007f and set tailLines to 7",
		},
		{
			name:         "Flag with very large number",
			args:         []string{"-999999f", "other", "args"},
			expectedTail: 999999,
			expectedArgs: []string{"other", "args"},
			description:  "Should parse very large numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset tailLines to default value
			tailLines = 100

			// Backup original os.Args
			originalArgs := os.Args
			defer func() {
				os.Args = originalArgs
			}()

			// Set up test args
			os.Args = append([]string{"ktail"}, tt.args...)

			// Call the function
			parseCustomFlags()

			// Check tailLines
			if tailLines != tt.expectedTail {
				t.Errorf("parseCustomFlags() tailLines = %d, want %d", tailLines, tt.expectedTail)
			}

			// Check that os.Args was modified correctly
			// os.Args[0] should always be the program name, so we check from index 1
			actualArgs := os.Args[1:]
			if !reflect.DeepEqual(actualArgs, tt.expectedArgs) {
				t.Errorf("parseCustomFlags() modified os.Args = %v, want %v", actualArgs, tt.expectedArgs)
			}
		})
	}
}

func TestParseCustomFlagsEdgeCases(t *testing.T) {
	t.Run("Only program name", func(t *testing.T) {
		originalArgs := os.Args
		defer func() {
			os.Args = originalArgs
		}()

		os.Args = []string{"ktail"}
		tailLines = 100

		parseCustomFlags()

		if tailLines != 100 {
			t.Errorf("Expected tailLines to remain 100, got %d", tailLines)
		}

		if len(os.Args) != 1 || os.Args[0] != "ktail" {
			t.Errorf("Expected os.Args to remain unchanged, got %v", os.Args)
		}
	})

	t.Run("Flag with only f", func(t *testing.T) {
		originalArgs := os.Args
		defer func() {
			os.Args = originalArgs
		}()

		os.Args = []string{"ktail", "-f"}
		tailLines = 100

		parseCustomFlags()

		if tailLines != 100 {
			t.Errorf("Expected tailLines to remain 100, got %d", tailLines)
		}

		expectedArgs := []string{"-f"}
		actualArgs := os.Args[1:]
		if !reflect.DeepEqual(actualArgs, expectedArgs) {
			t.Errorf("Expected os.Args[1:] to be %v, got %v", expectedArgs, actualArgs)
		}
	})

	t.Run("Flag with empty number", func(t *testing.T) {
		originalArgs := os.Args
		defer func() {
			os.Args = originalArgs
		}()

		os.Args = []string{"ktail", "-f", "other"}
		tailLines = 100

		parseCustomFlags()

		if tailLines != 100 {
			t.Errorf("Expected tailLines to remain 100, got %d", tailLines)
		}

		expectedArgs := []string{"-f", "other"}
		actualArgs := os.Args[1:]
		if !reflect.DeepEqual(actualArgs, expectedArgs) {
			t.Errorf("Expected os.Args[1:] to be %v, got %v", expectedArgs, actualArgs)
		}
	})
}

// Benchmark test for performance
func BenchmarkParseCustomFlags(b *testing.B) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = []string{"ktail", "-1000f", "other", "args", "here"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tailLines = 100
		parseCustomFlags()
	}
}
