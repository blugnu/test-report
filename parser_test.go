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
		fmt.Printf("chdir 'testdata': %s\n", err)
		return
	}

	testPackages := func(dest string, pkgs ...string) {
		cmd := exec.Command("go", append([]string{"test", "-json"}, pkgs...)...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%s: %s\n", cmd, err)
		}

		file, err := os.Create(dest)
		if err != nil {
			fmt.Printf("create testdata/%s: %s\n", dest, err)
		}
		defer file.Close()
		if err := cmd.Run(); err != nil {
			fmt.Printf("run '%s': %s\n", cmd, err)
		}

		if _, err := io.Copy(file, bytes.NewReader(output)); err != nil {
			fmt.Printf("save output to 'testdata/%s': %s\n", dest, err)
		}
	}

	testPackages("packages.json", "./pkga", "./pkgb")
	testPackages("no-test-files.json", "./no-test-files")
	testPackages("no-code.json", "./no-code")
}

func TestParse(t *testing.T) {
	// ARRANGE
	generate()

	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "packages.json",
			exec: func(t *testing.T) {
				report := &testrun{}
				input, err := os.Open("./testdata/packages.json")
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
					"pkgb_test.go:11": {
						"this test fails",
						"with four",
						"lines of output",
						"  the last is indented",
						"raw output is not indented (unlike test failure output)",
					},
				})
			},
		},
		{scenario: "no-test-files.json",
			exec: func(t *testing.T) {
				report := &testrun{}
				input, err := os.Open("./testdata/no-test-files.json")
				if err != nil {
					t.Fatalf("error loading test data: %s", err)
				}
				defer input.Close()
				p := parser{}

				// ACT
				err = p.parse(input, report)

				// ASSERT
				test.Error(t, err).IsNil()
				test.That(t, len(report.packages)).Equals(1, "number of packages")
				test.That(t, report.numTests).Equals(0, "number of tests")
				test.That(t, report.numPassed).Equals(0, "tests passed")
				test.That(t, report.numFailed).Equals(0, "tests failed")
				test.That(t, report.numSkipped).Equals(0, "tests skipped")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
