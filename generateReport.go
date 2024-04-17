package main

import (
	"fmt"
	"io/fs"
	"os"
)

// reportMode is the mode of the report to generate.
//
//	rmAllTests       // all tests
//	rmFailedTests    // failed tests only
//	rmSummaryOnly    // summary only
//
// The zero-value is rmFailedTests.
type reportMode int

const (
	rmFailedTests reportMode = iota //
	rmAllTests
	rmSummaryOnly
)

// function variables to facilitate testing
var (
	osExit     = os.Exit
	osFileMode = func(file fs.FileInfo) fs.FileMode {
		return file.Mode()
	}
)

// generateReport is a command that generates a report.
type generateReport struct {
	title    string
	mode     reportMode
	filename string
}

// checkError is a method that checks for an error.  If an error is found,
// is is printed to the console and the program terminated.
func (gen generateReport) checkError(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		osExit(1)
	}
}

// run is a method that generates a report.
func (gen generateReport) run(opts *opts) {
	gen.checkError(gen.checkPipe())

	td := &testrun{}
	p := &parser{}
	gen.checkError(p.parse(os.Stdin, td))

	output, err := os.Create(gen.filename)
	gen.checkError(err)
	defer output.Close()

	md := &markdown{
		title:   gen.title,
		mode:    gen.mode,
		testrun: td,
	}
	gen.checkError(md.export(output))
}

// checkPipe is a method that checks if the program is being piped input.
func (gen generateReport) checkPipe() error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if (osFileMode(stat) & os.ModeCharDevice) != 0 {
		return ErrNotPiped
	}

	return nil
}
