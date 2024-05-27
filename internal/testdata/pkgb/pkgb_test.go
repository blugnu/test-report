package pkgb

import (
	"fmt"
	"testing"
)

func TestPasses(t *testing.T) {}

func TestFails(t *testing.T) {
	t.Error("this test fails\nwith four\nlines of output\n  the last is indented")
	fmt.Print("raw output is not indented (unlike test failure output)")
}

func TestSkipped(t *testing.T) {
	t.Skip("this test is skipped")
}

func TestSubtest(t *testing.T) {
	t.Run("subtest", func(t *testing.T) {
		t.Run("fails", func(t *testing.T) {
			t.Error("this test fails")
		})

		t.Run("passes", func(t *testing.T) {
			t.Log("this test passes")
			t.Log("with two lines of output")
		})

		t.Run("skipped", func(t *testing.T) {
			t.Skip("this test is skipped")
		})
	})
}
