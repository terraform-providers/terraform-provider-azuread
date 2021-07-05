package validate

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestMailNickname(t *testing.T) {
	cases := []struct {
		Value    string
		TestName string
		ErrCount int
	}{
		{
			Value:    "initech",
			TestName: "Alpha",
			ErrCount: 0,
		},
		{
			Value:    "12345",
			TestName: "Numeric",
			ErrCount: 0,
		},
		{
			Value:    "floor3",
			TestName: "Alphanumeric",
			ErrCount: 0,
		},
		{
			Value:    "alias@",
			TestName: "At-sign",
			ErrCount: 1,
		},
		{
			Value:    "al\\ias",
			TestName: "Backslash",
			ErrCount: 1,
		},
		{
			Value:    "bob,bob",
			TestName: "Comma",
			ErrCount: 1,
		},
		{
			Value:    "group[1]",
			TestName: "Brackets",
			ErrCount: 1,
		},
		{
			Value:    "case(mondays)",
			TestName: "Parentheses",
			ErrCount: 1,
		},
		{
			Value:    "b0bby;.Tables\";",
			TestName: "QuotesColons",
			ErrCount: 1,
		},
		{
			Value:    "email me at this address",
			TestName: "Spaces",
			ErrCount: 1,
		},
		{
			Value:    "Bill<tps>",
			TestName: "LtGt",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			diags := MailNickname(tc.Value, cty.Path{})

			if len(diags) != tc.ErrCount {
				t.Fatalf("Expected MailNickname to have %d not %d errors for %q", tc.ErrCount, len(diags), tc.TestName)
			}
		})
	}
}
