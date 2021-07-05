package tf

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ErrorDiagF(err error, format string, a ...interface{}) diag.Diagnostics {
	return ErrorDiagPathF(err, "", format, a...)
}

func ErrorDiagPathF(err error, attr string, summary string, a ...interface{}) diag.Diagnostics {
	d := diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf(summary, a...),
	}
	if err != nil {
		d.Detail = err.Error()
	}
	if attr != "" {
		d.AttributePath = cty.Path{cty.GetAttrStep{Name: attr}}
	}
	return diag.Diagnostics{d}
}

func ImportAsDuplicateError(resourceName, id, name string) error {
	d := ImportAsDuplicateDiag(resourceName, id, name)
	if len(d) > 0 {
		return fmt.Errorf("%s. %s", d[0].Summary, d[0].Detail)
	}
	return nil
}

func ImportAsDuplicateDiag(resourceName, id, name string) diag.Diagnostics {
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       fmt.Sprintf("An existing %q with name %q (ID: %q) was found and `prevent_duplicate_names` was specified", resourceName, name, id),
		Detail:        fmt.Sprintf("To be managed via Terraform, this resource needs to be imported into the State. Please see the resource documentation for %q for more information.", resourceName),
		AttributePath: cty.Path{cty.GetAttrStep{Name: "id"}},
	}}
}

func ImportAsExistsDiag(resourceName, id string) diag.Diagnostics {
	return diag.Diagnostics{diag.Diagnostic{
		Severity:      diag.Error,
		Summary:       fmt.Sprintf("A resource with the ID %q already exists", id),
		Detail:        fmt.Sprintf("To be managed via Terraform, this resource needs to be imported into the State. Please see the resource documentation for %q for more information.", resourceName),
		AttributePath: cty.Path{cty.GetAttrStep{Name: "id"}},
	}}
}
