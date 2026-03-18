package main

import (
	"os"
	"os/exec"
	"testing"
)

// binaryPath holds the path to the compiled binary built once for the entire test suite.
var binaryPath string

// TestMain builds the binary once before all tests run and cleans it up afterward.
// This avoids the overhead of recompiling via "go run" on every test invocation.
func TestMain(m *testing.M) {
	// Create a temporary file to hold the compiled binary
	tmpFile, err := os.CreateTemp("", "hello-*")
	if err != nil {
		panic("failed to create temp file for binary: " + err.Error())
	}
	if err = tmpFile.Close(); err != nil {
		panic("failed to close temp file: " + err.Error())
	}
	binaryPath = tmpFile.Name()

	// Build the binary once for the entire test suite
	build := exec.Command("go", "build", "-o", binaryPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build binary: " + err.Error() + "\n" + string(out))
	}

	// Run all tests, then clean up the binary
	code := m.Run()
	_ = os.Remove(binaryPath)
	os.Exit(code)
}

func TestCLIExecutable(t *testing.T) {
	// Run the pre-built binary instead of recompiling with "go run"
	cmd := exec.Command(binaryPath, "Gopher")

	// Capture both standard output and standard error
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Subprocess execution failed: %v", err)
	}

	// Convert the byte slice to a string and verify the output
	actualOutput := string(outputBytes)
	expectedOutput := "Hello, Gopher!\n"

	if actualOutput != expectedOutput {
		t.Errorf("Expected CLI output %q, got %q", expectedOutput, actualOutput)
	}
}
