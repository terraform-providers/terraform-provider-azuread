package validate

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func UUID(i interface{}, path cty.Path) (ret diag.Diagnostics) {
	v, ok := i.(string)
	if !ok {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Expected a string value",
			AttributePath: path,
		})
		return
	}

	if _, err := uuid.ParseUUID(v); err != nil {
		ret = append(ret, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Value must be a valid UUID",
			AttributePath: path,
		})
	}

	return
}
