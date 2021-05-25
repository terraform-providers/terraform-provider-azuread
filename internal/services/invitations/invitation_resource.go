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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_email_address": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.StringIsEmailAddress,
			},
			"redirect_url": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.IsHTTPOrHTTPSURL,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_display_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},
			"user_message_info": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				MaxItems:     1,
				RequiredWith: []string{"send_invitation_message"},
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
						},
						"message_language": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validate.NoEmptyStrings,
						},
					},
				},
			},
			"send_invitation_message": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"redeem_url": {
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
