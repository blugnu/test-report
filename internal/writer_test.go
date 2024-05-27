package internal

import (
	"bytes"
	"errors"
	"testing"

	"github.com/blugnu/test"
)

func TestWriter(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "in error state",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				w := &IndentWriter{
					output: buf,
					error:  errors.New("writer error "),
				}

				// ACT
				w.Write("this should not be written")

				// ASSERT
				test.That(t, buf.Bytes()).IsNil()
			},
		},
		{scenario: "WriteXMLElement bare",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				w := &IndentWriter{output: buf}

				// ACT
				w.WriteXMLElement(func() { w.WriteLn("content") }, "tag")

				// ASSERT
				test.That(t, buf.String()).Equals("<tag>\n  content\n</tag>\n")
			},
		},
		{scenario: "WriteXMLElement with attributes",
			exec: func(t *testing.T) {
				// ARRANGE
				buf := bytes.NewBuffer(nil)
				w := &IndentWriter{output: buf}

				// ACT
				w.WriteXMLElement(func() { w.WriteLn("content") }, "tag", "attr1='1'", "attr2")

				// ASSERT
				test.That(t, buf.String()).Equals("<tag attr1='1' attr2>\n  content\n</tag>\n")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
