package main

import (
	"testing"

	"github.com/blugnu/test"
)

func TestShowVersion(t *testing.T) {
	// ARRANGE
	cmd := &showVersion{}

	// ACT
	stdout, _ := test.CaptureOutput(t, func() {
		cmd.run(nil)
	})

	// ASSERT
	stdout.Equals([]string{"test-report v0.1.0"})
}
