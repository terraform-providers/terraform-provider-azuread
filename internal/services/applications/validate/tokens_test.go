package validate

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestRoleScopeClaimValue(t *testing.T) {
	cases := []struct {
		Value    string
		TestName string
		ErrCount int
	}{
		{
			Value:    "administer",
			TestName: "Valid_Alpha",
			ErrCount: 0,
		},
		{
			Value:    "administer123",
			TestName: "Valid_Alphanumeric",
			ErrCount: 0,
		},
		{
			Value:    "administer123!#$%&'()*+,-./:;<=>?@[]^+_`{|}~",
			TestName: "Valid_AllChars",
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(120),
			TestName: "Valid_MaxLength",
			ErrCount: 0,
		},
		{
			Value:    "",
			TestName: "Valid_Empty",
			ErrCount: 0,
		},
		{
			Value:    acctest.RandString(121),
			TestName: "Invalid_MaxLength",
			ErrCount: 1,
		},
		{
			Value:    ".administer",
			TestName: "Invalid_StartsWithPeriod",
			ErrCount: 1,
		},
		{
			Value:    "administer123£¥€¢",
			TestName: "Invalid_BadChars",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			diags := RoleScopeClaimValue(tc.Value, cty.Path{})

			if len(diags) != tc.ErrCount {
				t.Fatalf("Expected RoleScopeClaimValue to have %d not %d errors for %q", tc.ErrCount, len(diags), tc.TestName)
			}
		})
	}
}
