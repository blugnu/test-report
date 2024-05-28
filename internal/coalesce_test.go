package internal

import (
	"testing"

	"github.com/blugnu/test"
)

func TestCoalesce(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		{scenario: "(z,nz)",
			exec: func(t *testing.T) {
				// ACT
				result := coalesce(0, 1)

				// ASSERT
				test.That(t, result).Equals(1)
			},
		},
		{scenario: "(nz,nz)",
			exec: func(t *testing.T) {
				// ACT
				result := coalesce(1, 2)

				// ASSERT
				test.That(t, result).Equals(1)
			},
		},
		{scenario: "no args",
			exec: func(t *testing.T) {
				// ACT
				result := coalesce[int]()

				// ASSERT
				test.That(t, result).Equals(0)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
