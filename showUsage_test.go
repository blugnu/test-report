package main

import (
	"testing"

	"github.com/blugnu/test"
)

func TestShowUsage(t *testing.T) {
	// ARRANGE
	cmd := &showUsage{}

	// ACT
	stdout, _ := test.CaptureOutput(t, func() {
		cmd.run(nil)
	})

	// ASSERT
	stdout.Equals([]string{
		"Usage: test-report [options]",
		"Options:",
		"    -f, -full      complete test report (includes passed tests)",
		"    -s, -summary   summary only",
		"    -t, -title     report title (default: 'Test Report')",
		"",
		"    -o, -output    output filename (default: 'test-report.md')",
		"",
		"    -h, -help      show this help message",
	})
}
