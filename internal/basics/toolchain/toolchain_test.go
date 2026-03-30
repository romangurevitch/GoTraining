package main

import "testing"

func TestDemo(t *testing.T) {
	// A simple test to demonstrate 'go test'
	t.Log("Testing toolchain demo...")
	if false {
		t.Error("This should not happen")
	}
}
