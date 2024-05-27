package internal

import (
	"time"
)

// testResult is an enumeration of the possible results of a test.
//
//	trFailed    // the test failed
//	trPassed    // the test passed
//	trSkipped   // the test was skipped
//
// The zero value is trFailed.
type testResult int

const (
	trFailed  testResult = iota // the test failed
	trPassed                    // the test passed
	trSkipped                   // the test was skipped
)

// testinfo contains information about a single test.
type testinfo struct {
	path        string        // the path to (name of) the test
	result      testResult    // the result of the test
	elapsed     time.Duration // the time taken to run the test (if recorded)
	packageName string        // the name of the package containing the test

	// the output of the test; during parsing all output is added to
	// a "raw" item in the map.  Once all output has been parsed, the "raw"
	// output is parsed to identify the source locations for any output with
	// separate items created for each unique location and the output associated
	// with it.  The "raw" item is then removed from the map.
	output map[string][]string
}

// packageinfo contains information about a single package, including a
// slice of testinfo items for each test in the package.
type packageinfo struct {
	name    string        // the name of the package
	passed  bool          // true if all tests in the package passed
	elapsed time.Duration // the time taken to run all tests in the package (if recorded)
	tests   []*testinfo   // the tests in the package
}

// testrun contains information about a test run, including a slice of
// packageinfo items for each package in the test run.
type testrun struct {
	elapsed       time.Duration  // the time taken to run all tests (if recorded)
	packages      []*packageinfo // the packages in the test run
	numFailed     int            // the number of failed tests
	numPassed     int            // the number of passed tests
	numTests      int            // the total number of tests
	numSkipped    int            // the number of skipped tests
	percentPassed int            // the percentage of tests that passed
}
