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
		input          string
		expectedOutput Locale
		expectedError  error
	}{
		{
			name:          "invalid_locale_returns_error",
			input:         "asdfghkjdfsajhgfdkjhfghasdk",
			expectedError: fmt.Errorf("could not parse locale %s", "asdfghkjdfsajhgfdkjhfghasdk"),
		},
		{
			name:          "locale_with_underscore_returns_error",
			input:         "en_US",
			expectedError: fmt.Errorf("could not parse locale %s", "en_US"),
		},
		{
			name:           "valid_locale_in_lowercase_returns_success",
			input:          "en-us",
			expectedOutput: Locale{id: 4, name: "en-US"},
		},
		{
			name:           "valid_locale_in_uppercase_returns_success",
			input:          "JA-JP",
			expectedOutput: Locale{id: 1, name: "ja-JP"},
		},
		{
			name:           "valid_locale_in_mixed_case_returns_success",
			input:          "ko-KR",
			expectedOutput: Locale{id: 2, name: "ko-KR"},
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
