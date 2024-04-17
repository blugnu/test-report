package main

import "fmt"

// version of the program
const version = "v0.1.0" //FUTURE: this should be set by the build system.

// showVersion is a command that prints the version.
type showVersion struct{}

// run prints the version.
func (showVersion) run(*opts) {
	fmt.Printf("test-report %s\n", version)
}
