package main

import (
	"errors"
	"flag"
	"os"
	"testing"

	"github.com/blugnu/test"
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

				oga := os.Args
				os.Args = []string{"test-report"}
				defer func() { os.Args = oga }()

				defer test.ExpectPanic(nil).Assert(t)

				// ACT
				main()
			},
		},
		{scenario: "invalid invocation",
			exec: func(t *testing.T) {
				// ARRANGE
				oga := os.Args
				defer func() { os.Args = oga }()

				ogpf := parseFlags
				defer func() { parseFlags = ogpf }()

				os.Args = []string{"test-report", "-invalid-flag"}
				parseFlags = func(*flag.FlagSet, []string) error {
					return errors.New("test error")
				}

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
