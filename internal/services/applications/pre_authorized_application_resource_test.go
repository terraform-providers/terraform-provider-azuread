package applications_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/applications/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type PreAuthorizedApplicationResource struct{}

func TestAccPreAuthorizedApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_pre_authorized_application", "test")
	r := PreAuthorizedApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("authorized_app_id").Exists(),
				check.That(data.ResourceName).Key("permission_ids.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPreAuthorizedApplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_pre_authorized_application", "test")
	r := PreAuthorizedApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport(data)),
	})
}

func (PreAuthorizedApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Applications.ApplicationsClient
	client.BaseClient.DisableRetries = true

	id, err := parse.PreAuthorizedApplicationID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing Pre-Authorized Application ID: %v", err)
	}

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Application with object ID %q does not exist", id.ObjectId)
		}
		return nil, fmt.Errorf("failed to retrieve Application with object ID %q: %+v", id.ObjectId, err)
	}

	if app.Api != nil && app.Api.PreAuthorizedApplications != nil {
		for _, a := range *app.Api.PreAuthorizedApplications {
			if a.AppId != nil && strings.EqualFold(*a.AppId, id.AppId) {
				return utils.Bool(true), nil
			}
		}
	}

	return nil, fmt.Errorf("Pre-Authorized Application %q was not found for Application %q", id.AppId, id.ObjectId)
}

func (PreAuthorizedApplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_application" "authorized" {
  display_name = "acctestApp-authorized-%[1]d"
}

resource "azuread_application" "authorizer" {
  display_name = "acctestApp-authorizer-%[1]d"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Administer the application"
      admin_consent_display_name = "Administer"
      enabled                    = true
      id                         = "%[2]s"
      type                       = "Admin"
      value                      = "administer"
    }

    oauth2_permission_scope {
      admin_consent_description  = "Access the application"
      admin_consent_display_name = "Access"
      enabled                    = true
      id                         = "%[3]s"
      type                       = "User"
      user_consent_description   = "Access the application"
      user_consent_display_name  = "Access"
      value                      = "user_impersonation"
    }
  }
}

resource "azuread_pre_authorized_application" "test" {
  application_object_id = azuread_application.authorizer.object_id
  authorized_app_id     = azuread_application.authorized.application_id
  permission_ids        = ["%[2]s", "%[3]s"]
}
`, data.RandomInteger, data.UUID(), data.UUID())
}

func (r PreAuthorizedApplicationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_pre_authorized_application" "import" {
  application_object_id = azuread_pre_authorized_application.test.application_object_id
  authorized_app_id     = azuread_pre_authorized_application.test.authorized_app_id
  permission_ids        = ["%[2]s"]
}
`, r.basic(data), data.UUID())
}
