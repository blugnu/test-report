package main

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/blugnu/test"

	"github.com/blugnu/test-report/internal"
)

func TestMainFunc(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "valid invocation",
			exec: func(t *testing.T) {
				// ARRANGE
				ogin := os.Stdin
				os.Stdin, _ = os.CreateTemp(".", "stdin-test-*")
				defer func() {
					os.Remove(os.Stdin.Name())
					os.Stdin = ogin
				}()
				defer test.Using(&os.Args, []string{"test-report"})()
				defer test.Using(&osExit, func(int) {})()

				// ARRANGE ASSERT
				defer test.ExpectPanic(nil).Assert(t)

				// ACT
				main()
			},
		},
		{scenario: "invalid invocation",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Args, []string{"test-report", "-invalid-flag"})()
				defer test.Using(&internal.ParseFlags, func(*flag.FlagSet, []string) error {
					return errors.New("test error")
				})()

				// ARRANGE ASSERT
				defer test.ExpectPanic(nil).Assert(t)

				// ACT
				stdout, _ := test.CaptureOutput(t, main)

				// ASSERT
				stdout.Contains("test error")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
