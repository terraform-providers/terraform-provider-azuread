package validate

import (
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func MailNickname(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Expected a string value",
			AttributePath: path,
		})
		return
	}

	if regexp.MustCompile(`[@()\\\[\]";:.<>, ]`).MatchString(v) {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Value cannot contain these characters: @()\\[]\";:.<>,SPACE",
			AttributePath: path,
		})
	}

	return
}
