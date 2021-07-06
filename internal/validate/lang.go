package validate

import (
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"golang.org/x/text/language"
)

func ISO639Language(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Expected a string value",
			AttributePath: path,
		})
		return
	}

	if _, err := language.Parse(v); err != nil && strings.Contains(err.Error(), "not well-formed") {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Language value is not well-formed",
			AttributePath: path,
		})
		return
	}

	return
}
