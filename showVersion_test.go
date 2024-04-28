package main

import (
	"testing"

	"github.com/blugnu/test"
)

func using[T any](v *T, tv T) func() {
	og := *v
	*v = tv
	return func() { *v = og }
}

func TestShowVersion(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "version not set",
			exec: func(t *testing.T) {
				// ARRANGE
				cmd := &showVersion{}
				defer using(&version, "")()

				// ACT
				stdout, _ := test.CaptureOutput(t, func() {
					cmd.run(nil)
				})

				// ASSERT
				stdout.Equals([]string{"test-report <unknown version>", "(built from source)"})
			},
		},
		{scenario: "version set",
			exec: func(t *testing.T) {
				// ARRANGE
				cmd := &showVersion{}
				defer using(&version, "1.0.0")()
				defer using(&commit, "commit")()
				defer using(&date, "date")()
				defer using(&builtBy, "builtBy")()
				defer using(&goVersion, "goVersion")()

				// ACT
				stdout, _ := test.CaptureOutput(t, func() {
					cmd.run(nil)
				})

				// ASSERT
				stdout.Equals([]string{
					"test-report v1.0.0",
					"commit: commit",
					"date: date",
					"built by: builtBy",
					"go version: goVersion",
				})
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
