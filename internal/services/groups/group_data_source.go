package groups

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
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func groupDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: groupDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Description:      "The display name for the group",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"display_name", "object_id"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"object_id": {
				Description:      "The object ID of the group",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"display_name", "object_id"},
				ValidateDiagFunc: validate.UUID,
			},

			"mail_enabled": {
				Description: "Whether the group is mail-enabled",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},

			"security_enabled": {
				Description: "Whether the group is a security group",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},

			"description": {
				Description: "The optional description of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"members": {
				Description: "The object IDs of the group members",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"owners": {
				Description: "The object IDs of the group owners",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"types": {
				Description: "A list of group types configured for the group. The only supported type is `Unified`, which specifies a Microsoft 365 group",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func groupDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient

	var group msgraph.Group
	var displayName string

	if v, ok := d.GetOk("display_name"); ok {
		displayName = v.(string)
	}

	var mailEnabled, securityEnabled *bool
	if v, exists := d.GetOk("mail_enabled"); exists {
		mailEnabled = utils.Bool(v.(bool))
	}
	if v, exists := d.GetOk("security_enabled"); exists {
		securityEnabled = utils.Bool(v.(bool))
	}

	if displayName != "" {
		filter := fmt.Sprintf("displayName eq '%s'", displayName)
		if mailEnabled != nil {
			filter = fmt.Sprintf("%s and mailEnabled eq %t", filter, *mailEnabled)
		}
		if securityEnabled != nil {
			filter = fmt.Sprintf("%s and securityEnabled eq %t", filter, *securityEnabled)
		}

		groups, _, err := client.List(ctx, odata.Query{Filter: filter})
		if err != nil {
			return tf.ErrorDiagPathF(err, "display_name", "No group found matching specified filter (%s)", filter)
		}

		count := len(*groups)
		if count > 1 {
			return tf.ErrorDiagPathF(err, "display_name", "More than one group found matching specified filter (%s)", filter)
		} else if count == 0 {
			return tf.ErrorDiagPathF(err, "display_name", "No group found matching specified filter (%s)", filter)
		}

		group = (*groups)[0]
	} else if objectId, ok := d.Get("object_id").(string); ok && objectId != "" {
		g, status, err := client.Get(ctx, objectId)
		if err != nil {
			if status == http.StatusNotFound {
				return tf.ErrorDiagPathF(nil, "object_id", "No group found with object ID: %q", objectId)
			}
			return tf.ErrorDiagF(err, "Retrieving group with object ID: %q", objectId)
		}
		if g == nil {
			return tf.ErrorDiagPathF(nil, "object_id", "Group not found with object ID: %q", objectId)
		}

		if mailEnabled != nil && (g.MailEnabled == nil || *g.MailEnabled != *mailEnabled) {
			var actual string
			if g.MailEnabled == nil {
				actual = "nil"
			} else {
				actual = fmt.Sprintf("%t", *g.MailEnabled)
			}
			return tf.ErrorDiagPathF(nil, "mail_enabled", "Group with object ID %q does not have the specified mail_enabled setting (expected: %t, actual: %s)", objectId, *mailEnabled, actual)
		}

		if securityEnabled != nil && (g.SecurityEnabled == nil || *g.SecurityEnabled != *securityEnabled) {
			var actual string
			if g.SecurityEnabled == nil {
				actual = "nil"
			} else {
				actual = fmt.Sprintf("%t", *g.SecurityEnabled)
			}
			return tf.ErrorDiagPathF(nil, "security_enabled", "Group with object ID %q does not have the specified security_enabled setting (expected: %t, actual: %s)", objectId, *securityEnabled, actual)
		}

		group = *g
	}

	if group.ID == nil {
		return tf.ErrorDiagF(errors.New("API returned group with nil object ID"), "Bad API Response")
	}

	d.SetId(*group.ID)

	tf.Set(d, "description", group.Description)
	tf.Set(d, "display_name", group.DisplayName)
	tf.Set(d, "mail_enabled", group.MailEnabled)
	tf.Set(d, "object_id", group.ID)
	tf.Set(d, "security_enabled", group.SecurityEnabled)
	tf.Set(d, "types", group.GroupTypes)

	members, _, err := client.ListMembers(ctx, d.Id())
	if err != nil {
		return tf.ErrorDiagF(err, "Could not retrieve group members for group with object ID: %q", d.Id())
	}
	tf.Set(d, "members", members)

	owners, _, err := client.ListOwners(ctx, d.Id())
	if err != nil {
		return tf.ErrorDiagF(err, "Could not retrieve group owners for group with object ID: %q", d.Id())
	}
	tf.Set(d, "owners", owners)

	return nil
}
