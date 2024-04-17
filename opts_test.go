package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/blugnu/test"
)

func TestOpts(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "parse/invalid flags",
			exec: func(t *testing.T) {
				// ARRANGE
				flagserr := fmt.Errorf("flags error")

				oga := os.Args
				defer func() { os.Args = oga }()
				os.Args = []string{"test-report", "-invalid-flag"}

				ogpf := parseFlags
				defer func() { parseFlags = ogpf }()
				parseFlags = func(f *flag.FlagSet, args []string) error {
					return flagserr
				}
				opts := &opts{}

				// ACT
				result, err := opts.parse()

				// ASSERT
				test.Error(t, err).Is(flagserr)
				test.That(t, result).IsNil()
			},
		},
		{scenario: "parse",
			exec: func(t *testing.T) {
				testcases := []struct {
					args   []string
					result interface{ run(*opts) }
				}{
					{args: []string{"version"}, result: showVersion{}},
					{args: []string{"-h"}, result: showUsage{}},
					{args: []string{"-help"}, result: showUsage{}},
					{args: []string{},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
						},
					},
					{args: []string{"-o", "report.md"},
						result: generateReport{
							filename: "report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
						},
					},
					{args: []string{"-output", "report.md"},
						result: generateReport{
							filename: "report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
						},
					},
					{args: []string{"-t", "My Title"},
						result: generateReport{
							filename: "test-report.md",
							title:    "My Title",
							mode:     rmFailedTests,
						},
					},
					{args: []string{"-title", "My Title"},
						result: generateReport{
							filename: "test-report.md",
							title:    "My Title",
							mode:     rmFailedTests,
						},
					},
					{args: []string{"-f"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmAllTests,
						},
					},
					{args: []string{"-full"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmAllTests,
						},
					},
					{args: []string{"-s"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmSummaryOnly,
						},
					},
					{args: []string{"-summary"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmSummaryOnly,
						},
					},
				}
				for _, tc := range testcases {
					t.Run(fmt.Sprintf("%s", tc.args), func(t *testing.T) {
						og := os.Args
						defer func() { os.Args = og }()

						os.Args = append([]string{"test-report"}, tc.args...)
						sut := &opts{}

						// ACT
						result, err := sut.parse()

						// ASSERT
						test.Error(t, err).IsNil()
						test.That(t, result).Equals(tc.result)
					})
				}
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ACT
			tc.exec(t)
		})
	}
}
