package internal

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/blugnu/test"
)

type fakeParser test.Fake[any]

func (fake fakeParser) parse(r io.Reader, tr *testrun) error {
	return fake.Err
}

func Test_osFileMode(t *testing.T) {
	// there are no meaningful tests for this function;
	// we exercise the code for coverage, for which we need
	// a valid file descriptor for which we will use the README
	// file from the testdata folder

	f, err := os.Open("./testdata/README.md")
	test.That(t, err).IsNil()

	fi, err := f.Stat()
	test.That(t, err).IsNil()

	// we don't care about the return value
	_ = osFileMode(fi)
}

func Test_mdExport(t *testing.T) {
	// there are no meaningful tests for this function;
	// we exercise the code for coverage, for which we need
	// a valid markdown object and a valid writer

	md := &markdown{testrun: &testrun{}}
	w := io.Discard

	// we don't care about the return value
	_ = mdExport(md, w)
}

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
				defer test.Using(&os.Stdin, nil)()

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

func TestGenerateReportRun(t *testing.T) {
	// ARRANGE
	exitCode := 0
	defer test.Using(&osExit, func(code int) {
		exitCode = code
	})()

	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "not piped",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Stdin, nil)()

				gen := &generateReport{}

				// ACT
				_ = gen.Run(&Options{})

				// ASSERT
				test.That(t, exitCode).Equals(-2)
			},
		},
		{scenario: "parser error",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Stdin, os.Stdout)() // we just need a non-nil *os.File
				defer test.Using(&osFileMode, func(file fs.FileInfo) fs.FileMode {
					return os.ModeNamedPipe
				})()

				sut := &generateReport{
					parser: fakeParser{Err: errors.New("test error")},
				}

				// ACT
				_ = sut.Run(&Options{})

				// ASSERT
				test.That(t, exitCode).Equals(-2)
			},
		},
		{scenario: "file creation error",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Stdin, os.Stdout)() // we just need a non-nil *os.File
				defer test.Using(&osFileMode, func(file fs.FileInfo) fs.FileMode {
					return os.ModeNamedPipe
				})()
				defer test.Using(&osCreate, func(name string) (*os.File, error) {
					return nil, errors.New("file creation error")
				})()

				sut := &generateReport{
					parser: fakeParser{},
				}

				// ACT
				_ = sut.Run(&Options{})

				// ASSERT
				test.That(t, exitCode).Equals(-2)
			},
		},
		{scenario: "markdown export error",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Stdin, os.Stdout)() // we just need a non-nil *os.File
				defer test.Using(&osFileMode, func(file fs.FileInfo) fs.FileMode {
					return os.ModeNamedPipe
				})()
				defer test.Using(&osCreate, func(name string) (*os.File, error) {
					return &os.File{}, nil
				})()
				defer test.Using(&mdExport, func(md *markdown, w io.Writer) error {
					return errors.New("markdown export error")
				})()

				sut := &generateReport{
					parser: fakeParser{},
				}

				// ACT
				_ = sut.Run(&Options{})

				// ASSERT
				test.That(t, exitCode).Equals(-2)
			},
		},
		{scenario: "success/all tests passed",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Stdin, os.Stdout)() // we just need a non-nil *os.File
				defer test.Using(&osFileMode, func(file fs.FileInfo) fs.FileMode {
					return os.ModeNamedPipe
				})()
				defer test.Using(&osCreate, func(name string) (*os.File, error) {
					return &os.File{}, nil
				})()
				defer test.Using(&mdExport, func(md *markdown, w io.Writer) error {
					return nil
				})()

				sut := &generateReport{
					parser: fakeParser{},
				}

				// ACT
				result := sut.Run(&Options{})

				// ASSERT
				test.That(t, result).Equals(0)
			},
		},
		{scenario: "success/some tests failed",
			exec: func(t *testing.T) {
				// ARRANGE
				defer test.Using(&os.Stdin, os.Stdout)() // we just need a non-nil *os.File
				defer test.Using(&osFileMode, func(file fs.FileInfo) fs.FileMode {
					return os.ModeNamedPipe
				})()
				defer test.Using(&osCreate, func(name string) (*os.File, error) {
					return &os.File{}, nil
				})()
				defer test.Using(&mdExport, func(md *markdown, w io.Writer) error {
					md.testrun.numFailed = 1
					return nil
				})()

				sut := &generateReport{
					parser: fakeParser{},
				}

				// ACT
				result := sut.Run(&Options{})

				// ASSERT
				test.That(t, result).Equals(-1)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			// ARRANGE
			exitCode = 0

			// ACT
			tc.exec(t)
		})
	}
}
