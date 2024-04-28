package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"time"
)

type line struct {
	Time    string   `json:"Time"`
	Action  string   `json:"Action"`
	Package string   `json:"Package"`
	Test    *string  `json:"Test,omitempty"`
	Elapsed *float64 `json:"Elapsed,omitempty"`
	Output  *string  `json:"Output,omitempty"`
}

func (line line) elapsedDur() time.Duration {
	return time.Duration(math.Round(*line.Elapsed*1000)) * time.Millisecond
}

type parser struct {
	pkgs   map[string]*packageinfo
	tests  map[string]map[string]*testinfo
	srcref *regexp.Regexp
}

func (p *parser) parse(r io.Reader, rpt *testrun) error {
	p.pkgs = map[string]*packageinfo{}
	p.tests = map[string]map[string]*testinfo{}
	p.srcref, _ = regexp.Compile(`(.*\.go:[0-9]*): (.*)\n`)

	*rpt = testrun{}

	decoder := json.NewDecoder(r)
	for {
		l := &line{}
		if err := decoder.Decode(l); err != nil {
			break
		}

		s, _ := json.Marshal(l)
		fmt.Println(string(s))

		if l.Test == nil && l.Elapsed != nil {
			rpt.elapsed = l.elapsedDur()
		}

		map[string]func(*line, *testrun){
			"start":  p.addPackage,
			"run":    p.addTest,
			"output": p.recordOutput,
			"pass":   p.recordPass,
			"fail":   p.recordFailure,
			"skip":   p.recordSkip,
		}[l.Action](l, rpt)
	}
	p.processOutput()

	if rpt.numTests > 0 {
		rpt.percentPassed = (rpt.numPassed * 100) / rpt.numTests
	}

	return nil
}

// addPackage adds a package to the testrun using the package name from
// the line, setting the initial state of the package passed flag to true.
func (p *parser) addPackage(line *line, rpt *testrun) {
	pi := &packageinfo{
		name:   line.Package,
		passed: true,
		tests:  []*testinfo{},
	}
	rpt.packages = append(rpt.packages, pi)
	p.pkgs[line.Package] = pi
	p.tests[line.Package] = map[string]*testinfo{}
}

// addTest adds a test to the testrun using the package name and test name
// from the line, setting the initial state of the test result to failed.
func (p *parser) addTest(line *line, rpt *testrun) {
	ti := &testinfo{
		path:        *line.Test,
		output:      map[string][]string{},
		packageName: line.Package,
	}
	p.tests[line.Package][*line.Test] = ti
	pkg := p.pkgs[line.Package]
	pkg.tests = append(pkg.tests, ti)
	rpt.numTests++
}

// recordOutput records the output of a test, adding it to the testinfo output
// map "raw" item.  If the output is a test result, the test result is updated
// and the output is not recorded.
func (p *parser) recordOutput(line *line, rpt *testrun) {
	if line.Test == nil || strings.HasPrefix(*line.Output, "=== RUN") {
		return
	}

	test := p.tests[line.Package][*line.Test]
	if strings.HasPrefix(*line.Output, "--- FAIL") {
		test.result = trFailed
		//FUTURE: extract the elapsed time from the output (formatted in the output string)
		return
	}
	if strings.HasPrefix(*line.Output, "--- PASS") {
		test.result = trPassed
		//FUTURE: extract the elapsed time from the output (formatted in the output string)
		return
	}

	test.output["raw"] = append(test.output["raw"], *line.Output)
}

// recordPass records a test pass, updating the test result and incrementing
// the number of passed tests in the testrun.  If the line has an elapsed time
// with no associated test, the elapsed time is updated in the package info.
func (p *parser) recordPass(line *line, rpt *testrun) {
	switch {
	case line.Test != nil:
		rpt.numPassed++
		p.tests[line.Package][*line.Test].result = trPassed
	case line.Elapsed != nil:
		p.pkgs[line.Package].elapsed = line.elapsedDur()
	}
}

// recordFailure records a test failure, updating the test result and incrementing
// the number of failed tests in the testrun.  If the line has an elapsed time
// with no associated test, the elapsed time is updated in the package info.
func (p *parser) recordFailure(line *line, rpt *testrun) {
	p.pkgs[line.Package].passed = false
	switch {
	case line.Test != nil:
		rpt.numFailed++
	case line.Elapsed != nil:
		p.pkgs[line.Package].elapsed = line.elapsedDur()
	}
}

// recordSkip records a test skip, updating the test result and incrementing
// the number of skipped tests in the testrun.
func (p *parser) recordSkip(line *line, rpt *testrun) {
	if line.Test == nil {
		p.pkgs[line.Package].passed = false
		return
	}
	p.tests[line.Package][*line.Test].result = trSkipped
	rpt.numSkipped++
}

// processOutput calls processTestOutput for each test that has "raw" output.
func (p *parser) processOutput() {
	for _, tests := range p.tests {
		for _, test := range tests {
			if len(test.output["raw"]) > 0 {
				p.processTestOutput(test)
				delete(test.output, "raw")
			}
		}
	}
}

// processTestOutput processes the raw output of a test, identifying
// the output emitted from each source location and storing it in
// the testinfo output map.
//
// The raw output is expected to be in the format:
//
//	|    <filename>:<line #>: <output line 1>
//	|        <output line 2>
//	|        ...
//	|        <output line N>
//
// The output is stored in the testinfo output map keyed by the
// "<filename>:<line #>""; the value is a slice of strings containing the
// output lines.  Indentation of output lines is preserved, after removing
// indentation introduced by the test runner (8 spaces for all lines apart
// from the initial line with, source reference, presented in-line
// with the source reference with no additional indentation).
func (p *parser) processTestOutput(test *testinfo) {
	skiplog, _ := regexp.Compile(fmt.Sprintf(`: %s \([0-9]+.[0-9]+s\)`, test.path))
	ref := ""
	for _, s := range test.output["raw"] {
		if s := p.srcref.FindAllStringSubmatch(s, -1); len(s) > 0 {
			ref = strings.TrimSpace(s[0][1])
			test.output[ref] = append([]string{}, s[0][2])
			continue
		}
		if test.result != trSkipped || !skiplog.MatchString(s) {
			if strings.HasPrefix(s, "        ") {
				s = s[8 : len(s)-1]
			}
			test.output[ref] = append(test.output[ref], s)
		}
	}
}
