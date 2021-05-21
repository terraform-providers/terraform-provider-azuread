package invitations

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func invitationResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: invitationResourceCreate,
		ReadContext:   invitationResourceRead,
		DeleteContext: invitationResourceDelete,

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			if _, err := uuid.ParseUUID(id); err != nil {
				return fmt.Errorf("specified ID (%q) is not valid: %s", id, err)
			}
			return nil
		}),

		Schema: map[string]*schema.Schema{
			"invited_user_display_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},
			"invited_user_email_address": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.StringIsEmailAddress,
			},
			"invited_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"invited_user_message_info": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cc_recipients": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"customized_message_body": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validate.NoEmptyStrings,
							ConflictsWith:    []string{"invited_user_message_info.message_language"},
						},
						"message_language": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validate.NoEmptyStrings,
							ConflictsWith:    []string{"invited_user_message_info.customized_message_body"},
						},
					},
				},
			},
			"invite_redirect_url": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.IsHTTPOrHTTPSURL,
			},
			"invite_redeem_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func invitationResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if useMsGraph := meta.(*clients.Client).EnableMsGraphBeta; useMsGraph {
		return invitationResourceCreateMsGraph(ctx, d, meta)
	}
	return tf.ErrorDiagF(nil, "azuread_invitation does not support AAD Graph. Please enable `EnableMsGraphBeta`")
}

func invitationResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if useMsGraph := meta.(*clients.Client).EnableMsGraphBeta; useMsGraph {
		return invitationResourceReadMsGraph(ctx, d, meta)
	}
	return tf.ErrorDiagF(nil, "azuread_invitation does not support AAD Graph. Please enable `EnableMsGraphBeta`")
}

func invitationResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if useMsGraph := meta.(*clients.Client).EnableMsGraphBeta; useMsGraph {
		return invitationResourceDeleteMsGraph(ctx, d, meta)
	}
	return tf.ErrorDiagF(nil, "azuread_invitation does not support AAD Graph. Please enable `EnableMsGraphBeta`")
}
