package serviceprincipals

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
	"github.com/hashicorp/terraform-provider-azuread/internal/helpers"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

const servicePrincipalResourceName = "azuread_service_principal"

func servicePrincipalResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: servicePrincipalResourceCreate,
		ReadContext:   servicePrincipalResourceRead,
		UpdateContext: servicePrincipalResourceUpdate,
		DeleteContext: servicePrincipalResourceDelete,

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
			"application_id": {
				Description:      "The application ID (client ID) of the application for which to create a service principal",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"app_role_assignment_required": {
				Description: "Whether this service principal requires an app role assignment to a user or group before Azure AD will issue a user or access token to the application",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			"display_name": {
				Description: "The display name of the application associated with this service principal",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"object_id": {
				Description: "The object ID of the service principal",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"app_roles": schemaAppRolesComputed(),

			"oauth2_permission_scopes": schemaOauth2PermissionScopesComputed(),

			"tags": {
				Description: "A set of tags to apply to the service principal",
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func servicePrincipalResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	properties := msgraph.ServicePrincipal{
		AccountEnabled:            utils.Bool(true),
		AppId:                     utils.String(d.Get("application_id").(string)),
		AppRoleAssignmentRequired: utils.Bool(d.Get("app_role_assignment_required").(bool)),
		Tags:                      tf.ExpandStringSlicePtr(d.Get("tags").(*schema.Set).List()),
	}

	servicePrincipal, _, err := client.Create(ctx, properties)
	if err != nil {
		return tf.ErrorDiagF(err, "Could not create service principal")
	}
	if servicePrincipal.ID == nil || *servicePrincipal.ID == "" {
		return tf.ErrorDiagF(errors.New("Object ID returned for service principal is nil"), "Bad API response")
	}
	d.SetId(*servicePrincipal.ID)

	return servicePrincipalResourceRead(ctx, d, meta)
}

func servicePrincipalResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	properties := msgraph.ServicePrincipal{
		ID:                        utils.String(d.Id()),
		AppRoleAssignmentRequired: utils.Bool(d.Get("app_role_assignment_required").(bool)),
		Tags:                      tf.ExpandStringSlicePtr(d.Get("tags").(*schema.Set).List()),
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Updating service principal with object ID: %q", d.Id())
	}

	return servicePrincipalResourceRead(ctx, d, meta)
}

func servicePrincipalResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient
	objectId := d.Id()

	servicePrincipal, status, err := client.Get(ctx, objectId)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Service Principal with Object ID %q was not found - removing from state!", objectId)
			d.SetId("")
			return nil
		}

		return tf.ErrorDiagF(err, "retrieving service principal with object ID: %q", d.Id())
	}

	tf.Set(d, "app_role_assignment_required", servicePrincipal.AppRoleAssignmentRequired)
	tf.Set(d, "app_roles", helpers.ApplicationFlattenAppRoles(servicePrincipal.AppRoles))
	tf.Set(d, "application_id", servicePrincipal.AppId)
	tf.Set(d, "display_name", servicePrincipal.DisplayName)
	tf.Set(d, "oauth2_permission_scopes", helpers.ApplicationFlattenOAuth2PermissionScopes(servicePrincipal.PublishedPermissionScopes))
	tf.Set(d, "object_id", servicePrincipal.ID)
	tf.Set(d, "tags", servicePrincipal.Tags)

	return nil
}

func servicePrincipalResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).ServicePrincipals.ServicePrincipalsClient

	_, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(fmt.Errorf("Service Principal was not found"), "id", "Retrieving service principal with object ID %q", d.Id())
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving service principal with object ID %q", d.Id())
	}

	status, err = client.Delete(ctx, d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Deleting service principal with object ID %q, got status %d", d.Id(), status)
	}

	return nil
}
