package main

import "fmt"

// version information: set by the build process
var (
	commit    string
	date      string
	version   string
	goVersion string
	builtBy   string
)

// showVersion is a command that prints the version.
type showVersion struct{}

// run prints the version.
func (showVersion) run(*opts) {
	if version == "" {
		fmt.Println("test-report <unknown version>")
		fmt.Println("(built from source)")
		return
	}
	fmt.Printf("test-report v%s\n", version)
	fmt.Println("commit:", commit)
	fmt.Println("date:", date)
	fmt.Println("built by:", builtBy)
	fmt.Println("go version:", goVersion)
}
