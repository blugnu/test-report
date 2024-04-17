This folder contains two packages implementing tests that exercise all outcomes
of a test run: test failures, passing tests and skipped tests.

The `generate()` function in `parser_test.go` is called to perform `go test -json`
for this testdata folder, to automatically generate the test data (.json) which
is then used by the tests implmented in `parser_test.go` itself.