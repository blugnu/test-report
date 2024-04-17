package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/blugnu/test"
)

func TestGetReportIcon(t *testing.T) {
	// ARRANGE
	md := &markdown{
		testrun: &testrun{},
	}

	t.Run("skipped tests, none failed", func(t *testing.T) {
		// ARRANGE
		md.testrun.numSkipped = 1
		defer func() { md.testrun.numSkipped = 0 }()

		// ACT
		result := md.getReportIcon()

		// ASSERT
		test.That(t, result).Equals(icon.yellowBook)
	})

	testcases := []struct {
		from   int
		to     int
		want   string
		assert func(*testing.T, string, string)
	}{
		{
			from: 0, to: 84,
			want: icon.redBook,
			assert: func(t *testing.T, got, wanted string) {
				test.That(t, got).Equals(wanted)
			},
		},
		{
			from: 85, to: 94,
			want: icon.orangeBook,
			assert: func(t *testing.T, got, wanted string) {
				test.That(t, got).Equals(wanted)
			},
		},
		{
			from: 95, to: 99,
			want: icon.yellowBook,
			assert: func(t *testing.T, got, wanted string) {
				test.That(t, got).Equals(wanted)
			},
		},
		{
			from: 100, to: 100,
			want: icon.greenBook,
			assert: func(t *testing.T, got, wanted string) {
				test.That(t, got).Equals(wanted)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%d-%d%%", tc.from, tc.to), func(t *testing.T) {
			for p := tc.from; p <= tc.to; p++ {
				// ARRANGE
				md.testrun.percentPassed = p

				// ACT
				result := md.getReportIcon()

				// ASSERT
				tc.assert(t, result, tc.want)
			}
		})
	}
}

func TestMarkdown(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		// summary tests
		{scenario: "summary/no tests",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:         rmSummaryOnly,
					testrun:      &testrun{},
					IndentWriter: &IndentWriter{output: buf},
				}

				// ACT
				md.writeSummary()

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<table>",
					"  <tr>",
					"    <td><b>packages</b></td>",
					"    <td>0</td>",
					"    <td>0s</td>",
					"    <td><b>tests</b></td>",
					"    <td align='right'>0</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>📕</td>",
					"    <td>passed</td>",
					"    <td align='right'>0%</td>",
					"  </tr>",
					"</table>",
					"",
				})
			},
		},
		{scenario: "summary/1 package, 2 tests, 1 failed",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:         rmSummaryOnly,
					testrun:      &testrun{},
					IndentWriter: &IndentWriter{output: buf},
				}
				md.testrun.packages = []*packageinfo{{}}
				md.testrun.numTests = 2
				md.testrun.numFailed = 1
				md.testrun.percentPassed = 50

				// ACT
				md.writeSummary()

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<table>",
					"  <tr>",
					"    <td><b>packages</b></td>",
					"    <td>1</td>",
					"    <td>0s</td>",
					"    <td><b>tests</b></td>",
					"    <td align='right'>2</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔴</td>",
					"    <td>failed</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>📕</td>",
					"    <td>passed</td>",
					"    <td align='right'>50%</td>",
					"  </tr>",
					"</table>",
					"",
				})
			},
		},
		{scenario: "summary/2 packages, 3 tests, 1 failed, 1 skipped, 1 passed",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:         rmSummaryOnly,
					testrun:      &testrun{},
					IndentWriter: &IndentWriter{output: buf},
				}
				md.testrun.packages = []*packageinfo{{}, {}} // 2 packages
				md.testrun.numTests = 3
				md.testrun.numFailed = 1
				md.testrun.numSkipped = 1
				md.testrun.percentPassed = 33

				// ACT
				md.writeSummary()

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<table>",
					"  <tr>",
					"    <td><b>packages</b></td>",
					"    <td>2</td>",
					"    <td>0s</td>",
					"    <td><b>tests</b></td>",
					"    <td align='right'>3</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔴</td>",
					"    <td>failed</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔕</td>",
					"    <td>skipped</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>📕</td>",
					"    <td>passed</td>",
					"    <td align='right'>33%</td>",
					"  </tr>",
					"</table>",
					"",
				})
			},
		},

		// output tests
		{scenario: "output/1 source, 1 line of output",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					IndentWriter: &IndentWriter{output: buf},
				}
				output := map[string][]string{
					"filename_test.go:12": {"output"},
				}

				// ACT
				md.writeOutput(output)

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<div><i>filename_test.go:12</i></div>",
					"<pre>output</pre>",
					"",
				})
			},
		},
		{scenario: "output/1 source, 2 lines of output",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					IndentWriter: &IndentWriter{output: buf},
				}
				output := map[string][]string{
					"filename_test.go:12": {
						"output line 1",
						"output line 2",
					},
				}

				// ACT
				md.writeOutput(output)

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<div><i>filename_test.go:12</i></div>",
					"<pre>output&nbsp;line&nbsp;1",
					"output&nbsp;line&nbsp;2</pre>",
					"",
				})
			},
		},
		{scenario: "output/2 sources, 1 with 3 lines of output",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					IndentWriter: &IndentWriter{output: buf},
				}
				output := map[string][]string{
					"filename_test.go:12": {
						"first output line 1",
						"first output line 2",
						"first output line 3",
					},
					"filename_test.go:14": {
						"second output",
					},
				}

				// ACT
				md.writeOutput(output)

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<div><i>filename_test.go:12</i></div>",
					"<pre>first&nbsp;output&nbsp;line&nbsp;1",
					"first&nbsp;output&nbsp;line&nbsp;2",
					"first&nbsp;output&nbsp;line&nbsp;3</pre>",
					"<div><i>filename_test.go:14</i></div>",
					"<pre>second&nbsp;output</pre>",
					"",
				})
			},
		},

		// package tests
		{scenario: "tests/failed tests mode, 3 tests, 1 failed, 1 skipped, 1 passed",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:         rmFailedTests,
					IndentWriter: &IndentWriter{output: buf},
				}
				pkg := &packageinfo{
					name:    "github.com/foo/package",
					elapsed: 6 * time.Millisecond,
					tests: []*testinfo{
						{path: "Test1", result: trFailed, elapsed: 1 * time.Millisecond},
						{path: "Test2", result: trSkipped, elapsed: 2 * time.Millisecond},
						{path: "Test3", result: trPassed, elapsed: 3 * time.Millisecond},
					},
				}

				// ACT
				md.writePackage(pkg)

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<tr>",
					"  <td>🔴</td>",
					"  <td colspan='2'><b>github.com/foo/package</b></td>",
					"  <td align='right'>6ms</td>",
					"</tr>",
					"<tr valign='top'>",
					"  <td></td>",
					"  <td>🔴</td>",
					"  <td>",
					"    <b>Test1</b>",
					"  </td>",
					"  <td align='right'>1ms</td>",
					"</tr>",
					"",
				})
			},
		},

		// detail tests
		{scenario: "detail/1 package, 1 failed test, 1 passed (failed tests mode)",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:         rmFailedTests,
					testrun:      &testrun{},
					IndentWriter: &IndentWriter{output: buf},
				}
				md.testrun.packages = []*packageinfo{{}}
				md.testrun.numTests = 2
				md.testrun.numFailed = 1
				md.testrun.numPassed = 1
				md.testrun.percentPassed = 50
				md.testrun.packages = []*packageinfo{{
					name:    "github.com/foo/package",
					elapsed: 3 * time.Millisecond,
					tests: []*testinfo{
						{path: "Test1", result: trFailed, elapsed: 1 * time.Millisecond},
						{path: "Test2", result: trPassed, elapsed: 2 * time.Millisecond},
					},
				}}

				// ACT
				md.writeDetail()

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<table>",
					"  <tr>",
					"    <td>🔴</td>",
					"    <td colspan='2'><b>github.com/foo/package</b></td>",
					"    <td align='right'>3ms</td>",
					"  </tr>",
					"  <tr valign='top'>",
					"    <td></td>",
					"    <td>🔴</td>",
					"    <td>",
					"      <b>Test1</b>",
					"    </td>",
					"    <td align='right'>1ms</td>",
					"  </tr>",
					"  <tr>",
					"    <td>✅</td>",
					"    <td colspan=3><b>1 test passed</b></td>",
					"  </tr>",
					"</table>",
					"",
				})
			},
		},
		{scenario: "detail/1 package, 1 failed test, 2 passed (failed tests mode)",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:         rmFailedTests,
					testrun:      &testrun{},
					IndentWriter: &IndentWriter{output: buf},
				}
				md.testrun.packages = []*packageinfo{{}}
				md.testrun.numTests = 3
				md.testrun.numFailed = 1
				md.testrun.numPassed = 2
				md.testrun.percentPassed = 66
				md.testrun.packages = []*packageinfo{{
					name:    "github.com/foo/package",
					elapsed: 6 * time.Millisecond,
					tests: []*testinfo{
						{path: "Test1", result: trFailed, elapsed: 1 * time.Millisecond},
						{path: "Test2", result: trPassed, elapsed: 2 * time.Millisecond},
						{path: "Test3", result: trPassed, elapsed: 3 * time.Millisecond},
					},
				}}

				// ACT
				md.writeDetail()

				// ASSERT
				test.Strings(t, buf.Bytes()).Equals([]string{
					"<table>",
					"  <tr>",
					"    <td>🔴</td>",
					"    <td colspan='2'><b>github.com/foo/package</b></td>",
					"    <td align='right'>6ms</td>",
					"  </tr>",
					"  <tr valign='top'>",
					"    <td></td>",
					"    <td>🔴</td>",
					"    <td>",
					"      <b>Test1</b>",
					"    </td>",
					"    <td align='right'>1ms</td>",
					"  </tr>",
					"  <tr>",
					"    <td>✅</td>",
					"    <td colspan=3><b>2 tests passed</b></td>",
					"  </tr>",
					"</table>",
					"",
				})
			},
		},

		// export tests
		{scenario: "export/1 package, 1 failed test, 1 skipped, 1 passed (summary only)",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:    rmSummaryOnly,
					title:   "Test Report",
					testrun: &testrun{},
				}
				md.testrun.elapsed = 6 * time.Millisecond
				md.testrun.packages = []*packageinfo{{}}
				md.testrun.numTests = 1
				md.testrun.numFailed = 1
				md.testrun.numSkipped = 1
				md.testrun.percentPassed = 33

				// ACT
				err := md.export(buf)

				// ASSERT
				test.Error(t, err).IsNil()
				test.Strings(t, buf.Bytes()).Equals([]string{
					"## 📕&nbsp;&nbsp;Test Report",
					"",
					"<table>",
					"  <tr>",
					"    <td><b>packages</b></td>",
					"    <td>1</td>",
					"    <td>6ms</td>",
					"    <td><b>tests</b></td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔴</td>",
					"    <td>failed</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔕</td>",
					"    <td>skipped</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>📕</td>",
					"    <td>passed</td>",
					"    <td align='right'>33%</td>",
					"  </tr>",
					"</table>",
					"",
					"<hr>",
					"",
					"_markdown test report generated by https://github.com/blugnu/test-report_",
					"",
				})
			},
		},
		{scenario: "export/1 package, 1 failed test, 1 skipped, 1 passed (failed tests mode)",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:    rmFailedTests,
					title:   "Test Report",
					testrun: &testrun{},
				}
				md.testrun.elapsed = 6 * time.Millisecond
				md.testrun.numTests = 3
				md.testrun.numFailed = 1
				md.testrun.numSkipped = 1
				md.testrun.numPassed = 1
				md.testrun.percentPassed = 33
				md.testrun.packages = []*packageinfo{{
					name:    "github.com/foo/package",
					elapsed: 6 * time.Millisecond,
					tests: []*testinfo{
						{path: "Test1", result: trFailed, elapsed: 1 * time.Millisecond},
						{path: "Test2", result: trSkipped, elapsed: 2 * time.Millisecond},
						{path: "Test3", result: trPassed, elapsed: 3 * time.Millisecond},
					},
				}}

				// ACT
				err := md.export(buf)

				// ASSERT
				test.Error(t, err).IsNil()
				test.Strings(t, buf.Bytes()).Equals([]string{
					"## 📕&nbsp;&nbsp;Test Report",
					"",
					"<table>",
					"  <tr>",
					"    <td><b>packages</b></td>",
					"    <td>1</td>",
					"    <td>6ms</td>",
					"    <td><b>tests</b></td>",
					"    <td align='right'>3</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔴</td>",
					"    <td>failed</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔕</td>",
					"    <td>skipped</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>📕</td>",
					"    <td>passed</td>",
					"    <td align='right'>33%</td>",
					"  </tr>",
					"</table>",
					"<table>",
					"  <tr>",
					"    <td>🔴</td>",
					"    <td colspan='2'><b>github.com/foo/package</b></td>",
					"    <td align='right'>6ms</td>",
					"  </tr>",
					"  <tr valign='top'>",
					"    <td></td>",
					"    <td>🔴</td>",
					"    <td>",
					"      <b>Test1</b>",
					"    </td>",
					"    <td align='right'>1ms</td>",
					"  </tr>",
					"  <tr>",
					"    <td>🔕</td>",
					"    <td colspan=3><b>1 test was skipped</b></td>",
					"  </tr>",
					"  <tr>",
					"    <td>✅</td>",
					"    <td colspan=3><b>1 test passed</b></td>",
					"  </tr>",
					"</table>",
					"",
					"<hr>",
					"",
					"_markdown test report generated by https://github.com/blugnu/test-report_",
					"",
				})
			},
		},
		{scenario: "export/1 package, 1 failed test, 1 skipped, 1 passed (all tests mode)",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				md := &markdown{
					mode:    rmAllTests,
					title:   "Test Report",
					testrun: &testrun{},
				}
				md.testrun.elapsed = 6 * time.Millisecond
				md.testrun.numTests = 3
				md.testrun.numFailed = 1
				md.testrun.numSkipped = 1
				md.testrun.numPassed = 1
				md.testrun.percentPassed = 33
				md.testrun.packages = []*packageinfo{{
					name:    "github.com/foo/package",
					elapsed: 6 * time.Millisecond,
					tests: []*testinfo{
						{path: "Test1", result: trFailed, elapsed: 1 * time.Millisecond},
						{path: "Test2", result: trSkipped, elapsed: 2 * time.Millisecond},
						{path: "Test3", result: trPassed, elapsed: 3 * time.Millisecond},
					},
				}}

				// ACT
				err := md.export(buf)

				// ASSERT
				test.Error(t, err).IsNil()
				test.Strings(t, buf.Bytes()).Equals([]string{
					"## 📕&nbsp;&nbsp;Test Report",
					"",
					"<table>",
					"  <tr>",
					"    <td><b>packages</b></td>",
					"    <td>1</td>",
					"    <td>6ms</td>",
					"    <td><b>tests</b></td>",
					"    <td align='right'>3</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔴</td>",
					"    <td>failed</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>🔕</td>",
					"    <td>skipped</td>",
					"    <td align='right'>1</td>",
					"  </tr>",
					"  <tr>",
					"    <td colspan=3 align='right'>📕</td>",
					"    <td>passed</td>",
					"    <td align='right'>33%</td>",
					"  </tr>",
					"</table>",
					"<table>",
					"  <tr>",
					"    <td>🔴</td>",
					"    <td colspan='2'><b>github.com/foo/package</b></td>",
					"    <td align='right'>6ms</td>",
					"  </tr>",
					"  <tr valign='top'>",
					"    <td></td>",
					"    <td>🔴</td>",
					"    <td>",
					"      <b>Test1</b>",
					"    </td>",
					"    <td align='right'>1ms</td>",
					"  </tr>",
					"  <tr valign='top'>",
					"    <td></td>",
					"    <td>🔕</td>",
					"    <td>",
					"      <b>Test2</b>",
					"    </td>",
					"    <td align='right'>2ms</td>",
					"  </tr>",
					"  <tr valign='top'>",
					"    <td></td>",
					"    <td>✅</td>",
					"    <td>",
					"      <b>Test3</b>",
					"    </td>",
					"    <td align='right'>3ms</td>",
					"  </tr>",
					"</table>",
					"",
					"<hr>",
					"",
					"_markdown test report generated by https://github.com/blugnu/test-report_",
					"",
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
