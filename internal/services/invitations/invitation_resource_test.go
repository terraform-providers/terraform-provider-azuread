package invitations_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type InvitationResource struct{}

func TestAccInvitation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_invitation", "test")
	r := InvitationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("user_id").Exists(),
				check.That(data.ResourceName).Key("redeem_url").Exists(),
				check.That(data.ResourceName).Key("user_email_address").HasValue(fmt.Sprintf("test-user-%s@test.com", data.RandomString)),
				check.That(data.ResourceName).Key("redirect_url").HasValue("https://portal.azure.com"),
			),
		},
	})
}

func TestAccInvitation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_invitation", "test")
	r := InvitationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("user_id").Exists(),
				check.That(data.ResourceName).Key("redeem_url").Exists(),
				check.That(data.ResourceName).Key("user_email_address").HasValue(fmt.Sprintf("test-user-%s@test.com", data.RandomString)),
				check.That(data.ResourceName).Key("redirect_url").HasValue("https://portal.azure.com"),
				check.That(data.ResourceName).Key("user_display_name").HasValue("Test user"),
				check.That(data.ResourceName).Key("send_invitation_message").HasValue("true"),
			),
		},
	})
}

func (r InvitationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id *string

	userID := state.Attributes["user_id"]

	user, status, err := clients.Users.MsClient.Get(ctx, userID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("User with object ID %q does not exist", userID)
		}
		return nil, fmt.Errorf("failed to retrieve User with object ID %q: %+v", userID, err)
	}
	id = user.ID

	return utils.Bool(id != nil && *id == userID), nil
}

func (InvitationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_invitation" "test" {
	user_email_address = "test-user-%s@test.com"
	redirect_url       = "https://portal.azure.com"
}
`, data.RandomString)
}

func (InvitationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_invitation" "test" {
	user_email_address = "test-user-%s@test.com"
	redirect_url       = "https://portal.azure.com"

	user_display_name = "Test user"

	send_invitation_message = true

	user_message_info {
		cc_recipients           = ["test-user-%s@test.com"]
		customized_message_body = "Hello there! You are invited to join my Azure tenant."
		message_language        = "en-US"
	}
}
`, data.RandomString, data.RandomString)
}
