package main

import "testing"

func TestDemo(t *testing.T) {
	// A simple test to demonstrate 'go test'
	t.Log("Testing toolchain demo...")

	got := 1
	want := 1
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
