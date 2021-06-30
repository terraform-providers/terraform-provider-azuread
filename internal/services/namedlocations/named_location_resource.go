package namedlocations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func namedLocationResource() *schema.Resource {
	eactlyOneOf := []string{"ip", "country"}
	return &schema.Resource{
		CreateContext: namedLocationResourceCreate,
		ReadContext:   namedLocationResourceRead,
		UpdateContext: namedLocationResourceUpdate,
		DeleteContext: namedLocationResourceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			if _, err := uuid.ParseUUID(id); err != nil {
				return fmt.Errorf("specified ID (%q) is not valid: %s", id, err)
			}
			return nil
		}),

		Schema: map[string]*schema.Schema{

			"display_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"ip": {
				Type:         schema.TypeList,
				ExactlyOneOf: eactlyOneOf,
				Optional:     true,
				MaxItems:     1,
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_ranges": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"trusted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"country": {
				Type:         schema.TypeList,
				ExactlyOneOf: eactlyOneOf,
				Optional:     true,
				MaxItems:     1,
				ForceNew:     true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"countries_and_regions": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"include_unknown_countries_and_regions": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func namedLocationResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).NamedLocations.MsClient

	displayName := d.Get("display_name").(string)

	base := msgraph.BaseNamedLocation{
		DisplayName: utils.String(displayName),
	}

	if v, ok := d.GetOk("ip"); ok {
		properties := expandIPNamedLocation(v.([]interface{}))
		properties.BaseNamedLocation = &base
		location, _, err := client.CreateIP(ctx, *properties)
		if err != nil {
			return tf.ErrorDiagF(err, "Could not create named location")
		}

		if location.ID == nil || *location.ID == "" {
			return tf.ErrorDiagF(errors.New("Bad API response"), "Object ID returned for named location is nil/empty")
		}
		d.SetId(*location.ID)
		return namedLocationResourceRead(ctx, d, meta)

	}

	if v, ok := d.GetOk("country"); ok {
		properties := expandCountryNamedLocation(v.([]interface{}))
		properties.BaseNamedLocation = &base
		location, _, err := client.CreateCountry(ctx, *properties)
		if err != nil {
			return tf.ErrorDiagF(err, "Could not create named location")
		}

		if location.ID == nil || *location.ID == "" {
			return tf.ErrorDiagF(errors.New("Bad API response"), "Object ID returned for named location is nil/empty")
		}
		d.SetId(*location.ID)
		return namedLocationResourceRead(ctx, d, meta)

	}

	return tf.ErrorDiagF(errors.New("Could not match named location"), "The named location object provided couldn't be matched to a country/ip object")

}

func namedLocationResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).NamedLocations.MsClient

	base := msgraph.BaseNamedLocation{
		ID: utils.String(d.Id()),
	}

	if d.HasChange("display_name") {
		displayName := d.Get("display_name").(string)
		base.DisplayName = &displayName
	}

	if v, ok := d.GetOk("ip"); ok {
		properties := expandIPNamedLocation(v.([]interface{}))
		properties.BaseNamedLocation = &base

		if _, err := client.UpdateIP(ctx, *properties); err != nil {
			return tf.ErrorDiagF(err, "Could not update named location with ID: %q", d.Id())
		}

	}
	if v, ok := d.GetOk("country"); ok {
		properties := expandCountryNamedLocation(v.([]interface{}))
		properties.BaseNamedLocation = &base

		if _, err := client.UpdateCountry(ctx, *properties); err != nil {
			return tf.ErrorDiagF(err, "Could not update named location with ID: %q", d.Id())
		}

	}

	return nil
}

func namedLocationResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).NamedLocations.MsClient

	location, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Named Location with Object ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
	}

	if ipnl, ok := (*location).(msgraph.IPNamedLocation); ok {
		d.SetId(*ipnl.ID)
		tf.Set(d, "display_name", ipnl.DisplayName)
		tf.Set(d, "ip", flattenIPNamedLocation(&ipnl))
	}

	if cnl, ok := (*location).(msgraph.CountryNamedLocation); ok {
		d.SetId(*cnl.ID)
		tf.Set(d, "display_name", cnl.DisplayName)
		tf.Set(d, "country", flattenCountryNamedLocation(&cnl))
	}

	return nil
}

func namedLocationResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).NamedLocations.MsClient

	if _, ok := d.GetOk("ip"); ok {
		_, status, err := client.GetIP(ctx, d.Id())
		if err != nil {
			if status == http.StatusNotFound {
				log.Printf("[DEBUG] Named Location with ID %q already deleted", d.Id())
				return nil
			}

			return tf.ErrorDiagPathF(err, "id", "Retrieving named location with ID %q", d.Id())
		}
	}

	if _, ok := d.GetOk("country"); ok {
		_, status, err := client.GetCountry(ctx, d.Id())
		if err != nil {
			if status == http.StatusNotFound {
				log.Printf("[DEBUG] Named Location with ID %q already deleted", d.Id())
				return nil
			}

			return tf.ErrorDiagPathF(err, "id", "Retrieving named location with ID %q", d.Id())
		}
	}

	status, err := client.Delete(ctx, d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Deleting named location with ID %q, got status %d", d.Id(), status)
	}

	return nil
}

func expandIPNamedLocation(in []interface{}) *msgraph.IPNamedLocation {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.IPNamedLocation{}
	config := in[0].(map[string]interface{})

	ipRanges := config["ip_ranges"].([]interface{})
	trusted := config["trusted"]

	result.IPRanges = expandIPNamedLocationIPRange(ipRanges)
	result.IsTrusted = utils.Bool(trusted.(bool))

	return &result
}

func expandIPNamedLocationIPRange(in []interface{}) *[]msgraph.IPNamedLocationIPRange {

	if len(in) == 0 {
		return nil
	}

	result := make([]msgraph.IPNamedLocationIPRange, 0)
	for _, cidr := range in {
		result = append(result, msgraph.IPNamedLocationIPRange{
			CIDRAddress: utils.String(cidr.(string)),
		})
	}

	return &result

}

func expandCountryNamedLocation(in []interface{}) *msgraph.CountryNamedLocation {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.CountryNamedLocation{}
	config := in[0].(map[string]interface{})

	countriesAndRegions := config["countries_and_regions"].([]interface{})
	includeUnknown := config["include_unknown_countries_and_regions"]

	result.CountriesAndRegions = tf.ExpandStringSlicePtr(countriesAndRegions)
	result.IncludeUnknownCountriesAndRegions = utils.Bool(includeUnknown.(bool))

	return &result
}

func flattenIPNamedLocation(in *msgraph.IPNamedLocation) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"ip_ranges": flattenIPNamedLocationIPRange(in.IPRanges),
			"trusted":   in.IsTrusted,
		},
	}
}

func flattenIPNamedLocationIPRange(in *[]msgraph.IPNamedLocationIPRange) []interface{} {
	if len(*in) == 0 {
		return []interface{}{}
	}

	result := make([]string, 0)
	for _, cidr := range *in {
		result = append(result, *cidr.CIDRAddress)
	}

	return tf.FlattenStringSlicePtr(&result)
}

func flattenCountryNamedLocation(in *msgraph.CountryNamedLocation) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"countries_and_regions":                 tf.FlattenStringSlicePtr(in.CountriesAndRegions),
			"include_unknown_countries_and_regions": in.IncludeUnknownCountriesAndRegions,
		},
	}
}
