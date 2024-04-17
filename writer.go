package main

import (
	"fmt"
	"io"
	"strings"
)

// IndentWriter is a wrapper around an io.Writer that will consistently indent
// output and providing methods to simplify the writing of text output
// containing indented sections and the handling of errors.
//
// The IndentWriter is not thread-safe and should not be used concurrently.
//
// If the underlying writer returns an error, the error is stored in the
// IndentWriter and no further output is written by any IndentWriter methods.
type IndentWriter struct {
	output   io.Writer // the underlying writer
	indent   string    // the current indent string
	indented bool      // indicates if the indent has been written on the current line
	error              // any error that occurred during writing
}

// writeIndent writes the current indent string to the output writer (if not
// already written and not in an error state) and clears the indented flag.
func (w *IndentWriter) writeIndent() {
	if w.indented || w.error != nil {
		return
	}
	_, w.error = io.WriteString(w.output, w.indent)
	w.indented = w.error == nil
}

// WriteIndented calls the specified function with the indent string increased
// by two spaces. The indent is restored to its original value after the
// function returns.
func (w *IndentWriter) WriteIndented(fn func()) {
	og := w.indent
	defer func() { w.indent = og }()

	w.indent += "  "
	fn()
}

// Write writes the specified arguments to the output writer. If the writer is
// in an error state, no output is written.
//
// One or more arguments may be provided. If more than one argument is provided,
// the first argument is treated as a string to be used as a format string with
// the remaining arguments supplied as values.
//
// If not yet written, the current indent string is written to the output
// writer before the specified arguments.
//
// No newline is appended to the output.
func (w *IndentWriter) Write(s string, args ...any) {
	if w.error != nil {
		return
	}

	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	w.writeIndent()
	_, w.error = io.WriteString(w.output, s)
}

// WriteLn writes the specified arguments with a newline appended.
// Different numbers of arguments are treated differently:
//
//	no arguments          only a newline is written.
//
//	one argument          the argument is written and a new-line appended.
//
//	multiple arguments    the first argument must be a string
//	                      to be used as a format string with the
//	                      remaining arguments supplied as values.
//
// If the writer is in an error state, no output is written.
//
// If not yet written, the current indent string is written to the output
// writer before the specified arguments.
func (w *IndentWriter) WriteLn(args ...any) {
	defer func() { w.indented = false }()
	var (
		s string
		a []any
	)
	switch len(args) {
	case 0:
		w.Write("\n")
		return

	case 1:
		s = fmt.Sprintf("%s", args[0])

	default:
		s = args[0].(string)
		a = args[1:]
	}

	if len(a) > 0 {
		s = fmt.Sprintf(s, a...)
	}

	lines := strings.Split(s, "\n")
	for _, line := range lines {
		w.writeIndent()
		if len(line) > 0 {
			w.Write(line)
		}
		w.Write("\n")
	}
}

// WriteXMLElement writes an XML element with optional attributes. The content
// is written by calling the specified function with the indent increased by
// two spaces. The element is closed with the appropriate end tag.
func (w *IndentWriter) WriteXMLElement(fn func(), el string, attrs ...string) {
	w.Write("<%s", el)
	for _, attr := range attrs {
		w.Write(" ")
		w.Write(attr)
	}
	w.WriteLn(">")
	w.WriteIndented(fn)
	w.WriteLn("</%s>", el)
}
