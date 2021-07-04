package serviceprincipals

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/helpers"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func servicePrincipalData() *schema.Resource {
	return &schema.Resource{
		ReadContext: servicePrincipalDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Description:      "The object ID of the service principal",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"application_id", "display_name", "object_id"},
				ValidateDiagFunc: validate.UUID,
			},

			"display_name": {
				Description:      "The display name of the application associated with this service principal",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"application_id", "display_name", "object_id"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"application_id": {
				Description:      "The application ID (client ID) of the application associated with this service principal",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"application_id", "display_name", "object_id"},
				ValidateDiagFunc: validate.UUID,
			},

			"app_roles": schemaAppRolesComputed(),

			"oauth2_permission_scopes": schemaOauth2PermissionScopesComputed(),
		},
	}
}

func servicePrincipalDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	var servicePrincipal *msgraph.ServicePrincipal

	if v, ok := d.GetOk("object_id"); ok {
		objectId := v.(string)
		sp, status, err := client.Get(ctx, objectId)
		if err != nil {
			if status == http.StatusNotFound {
				return tf.ErrorDiagPathF(nil, "object_id", "Service principal with object ID %q was not found", objectId)
			}

			return tf.ErrorDiagPathF(err, "object_id", "Retrieving service principal with object ID %q", objectId)
		}

		servicePrincipal = sp
	} else if _, ok := d.GetOk("display_name"); ok {
		displayName := d.Get("display_name").(string)
		query := odata.Query{
			Filter: fmt.Sprintf("displayName eq '%s'", displayName),
		}

		result, _, err := client.List(ctx, query)
		if err != nil {
			return tf.ErrorDiagF(err, "Listing service principals for filter %q", query.Filter)
		}
		if result == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}

		for _, sp := range *result {
			if sp.DisplayName == nil {
				continue
			}

			if *sp.DisplayName == displayName {
				servicePrincipal = &sp
				break
			}
		}

		if servicePrincipal == nil {
			return tf.ErrorDiagF(nil, "No service principal found matching display name: %q", displayName)
		}
	} else {
		applicationId := d.Get("application_id").(string)
		query := odata.Query{
			Filter: fmt.Sprintf("appId eq '%s'", applicationId),
		}

		result, _, err := client.List(ctx, query)
		if err != nil {
			return tf.ErrorDiagF(err, "Listing service principals for filter %q", query.Filter)
		}
		if result == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}

		for _, sp := range *result {
			if sp.AppId == nil {
				continue
			}

			if *sp.AppId == applicationId {
				servicePrincipal = &sp
				break
			}
		}

		if servicePrincipal == nil {
			return tf.ErrorDiagF(nil, "No service principal found for application ID: %q", applicationId)
		}
	}

	if servicePrincipal.ID == nil {
		return tf.ErrorDiagF(errors.New("API returned service principal with nil object ID"), "Bad API Response")
	}

	d.SetId(*servicePrincipal.ID)

	tf.Set(d, "app_roles", helpers.ApplicationFlattenAppRoles(servicePrincipal.AppRoles))
	tf.Set(d, "application_id", servicePrincipal.AppId)
	tf.Set(d, "display_name", servicePrincipal.DisplayName)
	tf.Set(d, "oauth2_permission_scopes", helpers.ApplicationFlattenOAuth2PermissionScopes(servicePrincipal.PublishedPermissionScopes))
	tf.Set(d, "object_id", servicePrincipal.ID)

	return nil
}
