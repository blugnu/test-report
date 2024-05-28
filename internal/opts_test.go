package internal

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

				defer test.Using(&os.Args, []string{"test-report", "-invalid-flag"})()
				defer test.Using(&ParseFlags, func(f *flag.FlagSet, args []string) error {
					return flagserr
				})()

				opts := &Options{}

				// ACT
				result, err := opts.Parse()

				// ASSERT
				test.Error(t, err).Is(flagserr)
				test.That(t, result).IsNil()
			},
		},
		{scenario: "parse",
			exec: func(t *testing.T) {
				testcases := []struct {
					args   []string
					result interface{ Run(*Options) int }
				}{
					{args: []string{"version"}, result: showVersion{}},
					{args: []string{"-h"}, result: showUsage{}},
					{args: []string{"-help"}, result: showUsage{}},
					{args: []string{},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-o", "report.md"},
						result: generateReport{
							filename: "report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-output", "report.md"},
						result: generateReport{
							filename: "report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-t", "My Title"},
						result: generateReport{
							filename: "test-report.md",
							title:    "My Title",
							mode:     rmFailedTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-title", "My Title"},
						result: generateReport{
							filename: "test-report.md",
							title:    "My Title",
							mode:     rmFailedTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-f"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmAllTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-full"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmAllTests,
							parser:   &parser{},
						},
					},
					{args: []string{"-s"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmSummaryOnly,
							parser:   &parser{},
						},
					},
					{args: []string{"-summary"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmSummaryOnly,
							parser:   &parser{},
						},
					},
					{args: []string{"-v"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
							parser:   &parser{verbose: true},
						},
					},
					{args: []string{"--verbose"},
						result: generateReport{
							filename: "test-report.md",
							title:    "Test Report",
							mode:     rmFailedTests,
							parser:   &parser{verbose: true},
						},
					},
				}
				for _, tc := range testcases {
					t.Run(fmt.Sprintf("%s", tc.args), func(t *testing.T) {
						og := os.Args
						defer func() { os.Args = og }()

						os.Args = append([]string{"test-report"}, tc.args...)
						sut := &Options{}

						// ACT
						result, err := sut.Parse()

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
