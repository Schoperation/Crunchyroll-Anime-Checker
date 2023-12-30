package core

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLocaleFromId(t *testing.T) {
	testCases := []struct {
		name           string
		input          int
		expectedOutput Locale
		expectedError  error
	}{
		{
			name:          "zero_locale_id_returns_error",
			input:         0,
			expectedError: fmt.Errorf("could not parse locale id %d", 0),
		},
		{
			name:          "big_invalid_locale_id_returns_error",
			input:         9999,
			expectedError: fmt.Errorf("could not parse locale id %d", 9999),
		},
		{
			name:           "japanese_locale_id_returns_success",
			input:          LocaleJaJP,
			expectedOutput: Locale(LocaleJaJP),
		},
		{
			name:           "english_locale_id_returns_success",
			input:          LocaleEnUS,
			expectedOutput: Locale(LocaleEnUS),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			locale, err := NewLocaleFromId(tc.input)

			require.ErrorIs(t, errors.Unwrap(err), errors.Unwrap(tc.expectedError))
			require.Equal(t, tc.expectedOutput, locale)
		})
	}
}

func TestNewLocaleFromString(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput Locale
		expectedError  error
	}{
		{
			name:          "blank_locale_string_returns_error",
			input:         "",
			expectedError: fmt.Errorf("could not parse locale %s", ""),
		},
		{
			name:          "invalid_locale_string_returns_error",
			input:         "AHHHHHHHHHHHHHHHHHHHHHHHHH",
			expectedError: fmt.Errorf("could not parse locale %s", "AHHHHHHHHHHHHHHHHHHHHHHHHH"),
		},
		{
			name:           "japanese_locale_returns_success",
			input:          "ja-JP",
			expectedOutput: Locale(LocaleJaJP),
		},
		{
			name:           "japanese_locale_all_lowercase_returns_success",
			input:          "ja-jp",
			expectedOutput: Locale(LocaleJaJP),
		},
		{
			name:           "japanese_locale_all_caps_returns_success",
			input:          "JA-JP",
			expectedOutput: Locale(LocaleJaJP),
		},
		{
			name:           "japanese_locale_with_mixed_case_returns_success",
			input:          "Ja-jP",
			expectedOutput: Locale(LocaleJaJP),
		},
		{
			name:           "english_locale_string_returns_success",
			input:          "en-us",
			expectedOutput: Locale(LocaleEnUS),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			locale, err := NewLocaleFromString(tc.input)

			require.ErrorIs(t, errors.Unwrap(err), errors.Unwrap(tc.expectedError))
			require.Equal(t, tc.expectedOutput, locale)
		})
	}
}
