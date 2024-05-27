package internal

import (
	"fmt"
	"io"
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
	osCreate   = os.Create
	osExit     = os.Exit
	osFileMode = func(file fs.FileInfo) fs.FileMode {
		return file.Mode()
	}
	mdExport = func(md *markdown, w io.Writer) error {
		return md.export(w)
	}
)

// generateReport is a command that generates a report.
type generateReport struct {
	title    string
	mode     reportMode
	filename string
	parser   interface {
		parse(io.Reader, *testrun) error
	}
}

// checkError is a method that checks for an error.  If an error is found,
// is is printed to the console and the program terminated.
func (generateReport) checkError(err error) bool {
	if err == nil {
		return true
	}
	fmt.Println("ERROR:", err)
	osExit(-2)
	return false // in testing osExit is mocked so we need a valid return
}

// Run is a method that generates a report.
func (cmd generateReport) Run(opts *Options) int {
	if !cmd.checkError(cmd.checkPipe()) {
		return 1
	}

	td := &testrun{}
	if !cmd.checkError(cmd.parser.parse(os.Stdin, td)) {
		return 1
	}

	output, err := osCreate(cmd.filename)
	if !cmd.checkError(err) {
		return 1
	}
	defer output.Close()

	md := &markdown{
		title:   cmd.title,
		mode:    cmd.mode,
		testrun: td,
	}
	if !cmd.checkError(mdExport(md, output)) {
		return 1
	}

	return map[bool]int{
		true:  -1,
		false: 0,
	}[td.numFailed > 0]
}

// checkPipe is a method that checks if the program is being piped input.
func (generateReport) checkPipe() error {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if (osFileMode(stat) & os.ModeCharDevice) != 0 {
		return ErrNotPiped
	}

	return nil
}
