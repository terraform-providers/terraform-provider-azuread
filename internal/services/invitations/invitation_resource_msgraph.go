package invitations

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

func invitationResourceCreateMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Invitations.MsClient

	invitedUserEmailAddress := d.Get("invited_user_email_address").(string)
	inviteRedirectUrl := d.Get("invite_redirect_url").(string)

	properties := msgraph.Invitation{
		InvitedUserEmailAddress: utils.String(invitedUserEmailAddress),
		InviteRedirectURL:       utils.String(inviteRedirectUrl),
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
	d.Set("invited_user_id", *invitation.InvitedUser.ID)

	return invitationResourceReadMsGraph(ctx, d, meta)
}

func invitationResourceReadMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func invitationResourceDeleteMsGraph(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.MsClient

	userID := d.Get("invited_user_id").(string)

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
