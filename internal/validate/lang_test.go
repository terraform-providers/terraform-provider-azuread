package validate

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestISO639Language(t *testing.T) {
	cases := []struct {
		Value    string
		TestName string
		ErrCount int
	}{
		{
			Value:    "fr",
			TestName: "Language",
			ErrCount: 0,
		},
		{
			Value:    "es-MX",
			TestName: "LanguageLocale",
			ErrCount: 0,
		},
		{
			Value:    "cr-belter",
			TestName: "UnrecognisedButValid",
			ErrCount: 0,
		},
		{
			Value:    "en-en-GB",
			TestName: "Malformed",
			ErrCount: 1,
		},
		{
			Value:    "192",
			TestName: "NumericInvalid",
			ErrCount: 1,
		},
		{
			Value:    "@#$%^&*",
			TestName: "Gibberish",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			diags := ISO639Language(tc.Value, cty.Path{})

			if len(diags) != tc.ErrCount {
				t.Fatalf("Expected ISO639Language to have %d not %d errors for %q", tc.ErrCount, len(diags), tc.TestName)
			}
		})
	}
}
