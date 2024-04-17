package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/blugnu/test"
)

func generate() {
	pwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(pwd) }()

	if err := os.Chdir("testdata"); err != nil {
		fmt.Printf("cd testdata: %s\n", err)
	}

	cmd := exec.Command("go", "test", "./...", "-json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s: %s\n", cmd, err)
	}

	file, err := os.Create(".json")
	if err != nil {
		fmt.Printf("create testdata/.json: %s\n", err)
	}
	defer file.Close()
	if err := cmd.Run(); err != nil {
		fmt.Printf("run '%s': %s\n", cmd, err)
	}

	if _, err := io.Copy(file, bytes.NewReader(output)); err != nil {
		fmt.Printf("save output to 'testdata/.json': %s\n", err)
	}
}

func TestParse(t *testing.T) {
	// ARRANGE
	generate()

	report := &testrun{}
	input, err := os.Open("./testdata/.json")
	if err != nil {
		t.Fatalf("error loading test data: %s", err)
	}
	defer input.Close()
	p := parser{}

	// ACT
	err = p.parse(input, report)

	// ASSERT
	test.Error(t, err).IsNil()
	test.That(t, len(report.packages)).Equals(2, "number of packages")
	test.That(t, report.numTests).Equals(9, "number of tests")
	test.That(t, report.numPassed).Equals(3, "tests passed")
	test.That(t, report.numFailed).Equals(4, "tests failed")
	test.That(t, report.numSkipped).Equals(2, "tests skipped")

	test.Map(t, report.packages[1].tests[1].output).Equals(map[string][]string{
		"pkgb_test.go:8": {
			"this test fails",
			"with four",
			"lines of output",
			"  the last is indented",
		},
	})
	// test.That(t, report.packages[1].tests[1].output["pkgb_test.go:12"][0]).Equals("TestFails")
	// test.That(t, report.packages["github.com/blugnu/test-report/testdata/pkgb"]["TestFails"]).IsNotNil()
}
