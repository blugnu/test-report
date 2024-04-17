package main

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/blugnu/test"
)

func TestCheckError(t *testing.T) {
	// ARRANGE
	callsOSExit := false
	gen := &generateReport{}
	og := osExit
	defer func() { osExit = og }()
	osExit = func(c int) { callsOSExit = true }

	// ACT
	gen.checkError(nil)

	// ASSERT
	test.IsFalse(t, callsOSExit)

	// ACT
	gen.checkError(errors.New("test error"))

	// ASSERT
	test.IsTrue(t, callsOSExit)
}

func TestCheckPipe(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "stdin is nil",
			exec: func(t *testing.T) {
				// ARRANGE
				og := os.Stdin
				os.Stdin = nil
				defer func() { os.Stdin = og }()

				gen := &generateReport{}

				// ACT
				err := gen.checkPipe()

				// ASSERT
				test.Error(t, err).Is(os.ErrInvalid)

			},
		},
		{scenario: "stdin is not nil but not a valid pipe",
			exec: func(t *testing.T) {
				// ARRANGE
				{
					og := os.Stdin
					defer func() { os.Stdin = og }()
					os.Stdin = os.Stdout // we just need a non-nil *os.File
				}
				{
					og := osFileMode
					defer func() { osFileMode = og }()
					osFileMode = func(file fs.FileInfo) fs.FileMode {
						return os.ModeCharDevice
					}
				}

				gen := &generateReport{}

				// ACT
				err := gen.checkPipe()

				// ASSERT
				test.Error(t, err).Is(ErrNotPiped)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
