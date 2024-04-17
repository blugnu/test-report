package main

import "fmt"

// showUsage is a command that prints the usage message.
type showUsage struct{}

// run prints the usage message.
func (showUsage) run(*opts) {
	fmt.Println("Usage: test-report [options]")
	fmt.Println("Options:")
	fmt.Println("    -f, -full      complete test report (includes passed tests)")
	fmt.Println("    -s, -summary   summary only")
	fmt.Println("    -t, -title     report title (default: 'Test Report')")
	fmt.Println()
	fmt.Println("    -o, -output    output filename (default: 'test-report.md')")
	fmt.Println()
	fmt.Println("    -h, -help      show this help message")
}
