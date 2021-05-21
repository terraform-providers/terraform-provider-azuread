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
			),
		},
		data.ImportStep(),
	})
}

func (r InvitationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id *string

	userID := state.Attributes["invited_user_id"]

	user, status, err := clients.Users.MsClient.Get(ctx, userID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("User with object ID %q does not exist", userID)
		}
		return nil, fmt.Errorf("failed to retrieve User with object ID %q: %+v", userID, err)
	}
	id = user.ID

	return utils.Bool(id != nil && *id == state.ID), nil
}

func (InvitationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_invitation" "test" {
	invited_user_email_address = "test-user-%s@test.com"
	invite_redirect_url        = "https://portal.azure.com"
}
`, data.RandomString)
}
