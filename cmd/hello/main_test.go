package main

import (
	"os/exec"
	"testing"
)

func TestCLIExecutable(t *testing.T) {
	// Construct an exec.Command to invoke the Go toolchain
	// This executes the main.go file as an isolated subprocess
	cmd := exec.Command("go", "run", "main.go", "Gopher")

	// Capture both standard output and standard error
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Subprocess execution failed: %v", err)
	}

	// Convert the byte slice to a string and verify the output
	actualOutput := string(outputBytes)
	expectedOutput := "Hello, Gopher!\n" // Note the trailing newline from fmt.Println

	if actualOutput != expectedOutput {
		t.Errorf("Expected CLI output %q, got %q", expectedOutput, actualOutput)
	}
}
