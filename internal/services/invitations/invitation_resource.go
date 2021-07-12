package invitations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
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
						"customised_message_body": {
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
	client := meta.(*clients.Client).Invitations.MsClient

	invitedUserEmailAddress := d.Get("user_email_address").(string)
	inviteRedirectUrl := d.Get("redirect_url").(string)

	properties := msgraph.Invitation{
		InvitedUserEmailAddress: utils.String(invitedUserEmailAddress),
		InviteRedirectURL:       utils.String(inviteRedirectUrl),
	}

	if v, ok := d.GetOk("user_display_name"); ok {
		properties.InvitedUserDisplayName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("send_invitation_message"); ok {
		properties.SendInvitationMessage = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("user_message_info"); ok {
		// Since ValidateFunc and ValidateDiagFunc are not yet supported on lists or sets we must check for send_invitation_message value here
		if properties.SendInvitationMessage == nil || !*properties.SendInvitationMessage {
			return tf.ErrorDiagF(errors.New("Wrong value"), "When `user_message_info` is specified, `send_invitation_message` must be set to `true`")
		}

		properties.InvitedUserMessageInfo = expandInvitedUserMessageInfo(v.([]interface{}))
	}

	invitation, _, err := client.Create(ctx, properties)
	if err != nil {
		return tf.ErrorDiagF(err, "Could not create invitation")
	}

	if invitation.ID == nil || *invitation.ID == "" {
		return tf.ErrorDiagF(errors.New("Bad API response"), "Object ID returned for invitation is nil/empty")
	}
	d.SetId(*invitation.ID)

	if invitation.InvitedUser == nil || invitation.InvitedUser.ID == nil || *invitation.InvitedUser.ID == "" {
		return tf.ErrorDiagF(errors.New("Bad API response"), "Invited user object ID returned for invitation is nil/empty")
	}
	d.Set("user_id", invitation.InvitedUser.ID)

	if invitation.InviteRedeemURL == nil || *invitation.InviteRedeemURL == "" {
		return tf.ErrorDiagF(errors.New("Bad API response"), "Redeem URL returned for invitation is nil/empty")
	}
	d.Set("redeem_url", invitation.InviteRedeemURL)

	if err != nil {
		return tf.ErrorDiagF(err, "Waiting for User with object ID: %q", *invitation.InvitedUser.ID)
	}

	return invitationResourceRead(ctx, d, meta)
}

func invitationResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	userID := d.Get("user_id").(string)

	user, status, err := client.Get(ctx, userID)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] User with Object ID %q was not found - removing from state!", userID)
			d.Set("user_id", "")
			return nil
		}
		return tf.ErrorDiagF(err, "Retrieving user with object ID: %q", userID)
	}

	tf.Set(d, "user_id", user.ID)
	tf.Set(d, "user_email_address", user.Mail)

	return nil
}

func invitationResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	userID := d.Get("user_id").(string)

	_, status, err := client.Get(ctx, userID)
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(fmt.Errorf("User was not found"), "id", "Retrieving user with object ID %q", userID)
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving user with object ID %q", userID)
	}

	status, err = client.Delete(ctx, userID)
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Deleting user with object ID %q, got status %d", userID, status)
	}

	return nil
}

func expandInvitedUserMessageInfo(in []interface{}) *msgraph.InvitedUserMessageInfo {
	if len(in) == 0 {
		return nil
	}

	result := msgraph.InvitedUserMessageInfo{}
	config := in[0].(map[string]interface{})

	ccRecipients := config["cc_recipients"].([]interface{})
	messageBody := config["customised_message_body"].(string)
	messageLanguage := config["message_language"].(string)

	result.CCRecipients = expandRecipients(ccRecipients)
	result.CustomizedMessageBody = &messageBody
	result.MessageLanguage = &messageLanguage

	return &result
}

func expandRecipients(in []interface{}) *[]msgraph.Recipient {
	recipients := make([]msgraph.Recipient, 0, len(in))
	for _, recipientRaw := range in {
		recipient := recipientRaw.(string)

		newRecipient := msgraph.Recipient{
			EmailAddress: &msgraph.EmailAddress{
				Address: &recipient,
			},
		}

		recipients = append(recipients, newRecipient)
	}

	return &recipients
}
