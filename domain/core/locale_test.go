package core

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLocale(t *testing.T) {
	testCases := []struct {
		name           string
		input          int
		expectedOutput Locale
		expectedError  error
	}{
		{
			name:          "invalid_locale_returns_error",
			input:         0,
			expectedError: fmt.Errorf("could not parse locale id %d", 0),
		},
		{
			name:           "valid_locale_returns_success",
			input:          4,
			expectedOutput: Locale{id: 4, name: "en-US"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			locale, err := NewLocale(tc.input)

			require.ErrorIs(t, errors.Unwrap(err), errors.Unwrap(tc.expectedError))
			require.Equal(t, tc.expectedOutput, locale)
		})
	}
}
