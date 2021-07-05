package applications_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type ApplicationResource struct{}

func TestAccApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-APP-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_appRoles(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_role.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appRoleNoValue(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_role.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appRole(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_role.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.appRolesUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_role.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("app_role.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_duplicateAppRolesOauth2PermissionsValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.duplicateAppRolesOauth2PermissionsValues(data),
			ExpectError: regexp.MustCompile("validation failed: duplicate value found:"),
		},
	})
}

func TestAccApplication_groupMembershipClaimsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withGroupMembershipClaims(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_oauth2PermissionScopes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}
	scopeIDs := []string{
		data.UUID(),
		data.UUID(),
		data.UUID(),
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api.0.oauth2_permission_scope.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.oauth2PermissionScopes(data, scopeIDs),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api.0.oauth2_permission_scope.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.oauth2PermissionScopesUpdate(data, scopeIDs),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api.0.oauth2_permission_scope.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.oauth2PermissionScopes(data, scopeIDs),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("api.0.oauth2_permission_scope.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_owners(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("owners.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.singleOwner(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("owners.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.threeOwners(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("owners.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("owners.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApplication_preventDuplicateNamesPass(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.preventDuplicateNamesPass(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("prevent_duplicate_names"),
	})
}

func TestAccApplication_preventDuplicateNamesFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_application", "test")
	r := ApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.RequiresImportErrorStep(r.preventDuplicateNamesFail(data)),
	})
}

func (r ApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Applications.ApplicationsClient
	client.BaseClient.DisableRetries = true
	app, status, err := client.Get(ctx, state.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Application with object ID %q does not exist", state.ID)
		}
		return nil, fmt.Errorf("failed to retrieve Application with object ID %q: %+v", state.ID, err)
	}
	return utils.Bool(app.ID != nil && *app.ID == state.ID), nil
}

func (ApplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[1]d"
}
`, data.RandomInteger)
}

func (ApplicationResource) withGroupMembershipClaims(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name            = "acctest-APP-%[1]d"
  group_membership_claims = ["DirectoryRole", "SecurityGroup", "ApplicationGroup"]
}
`, data.RandomInteger)
}

func (ApplicationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name = "acctestUser.%[1]d@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d"
  password            = "%[2]s"
}

resource "azuread_application" "test" {
  display_name            = "acctest-APP-complete-%[1]d"
  identifier_uris         = ["api://hashicorptestapp-%[1]d"]
  group_membership_claims = ["All"]
  sign_in_audience        = "AzureADMultipleOrgs"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Administer the application"
      admin_consent_display_name = "Administer"
      enabled                    = true
      id                         = "%[3]s"
      type                       = "Admin"
      value                      = "administer"
    }

    oauth2_permission_scope {
      admin_consent_description  = "Allow the application to access acctest-APP-%[1]d on behalf of the signed-in user."
      admin_consent_display_name = "Access acctest-APP-%[1]d"
      enabled                    = true
      id                         = "%[4]s"
      type                       = "User"
      user_consent_description   = "Allow the application to access acctest-APP-%[1]d on your behalf."
      user_consent_display_name  = "Access acctest-APP-%[1]d"
      value                      = "user_impersonation"
    }
  }

  app_role {
    allowed_member_types = ["User"]
    description          = "Admins can manage roles and perform all task actions"
    display_name         = "Admin"
    enabled              = true
    id                   = "%[5]s"
    value                = ""
  }

  app_role {
    allowed_member_types = ["User"]
    description          = "ReadOnly roles have limited query access"
    display_name         = "ReadOnly"
    enabled              = true
    id                   = "%[6]s"
    value                = "User"
  }

  optional_claims {
    access_token {
      name = "myclaim"
    }

    access_token {
      name = "otherclaim"
    }

    id_token {
      name                  = "userclaim"
      source                = "user"
      essential             = true
      additional_properties = ["emit_as_roles"]
    }

    saml2_token {
      name = "samlexample"
    }
  }

  required_resource_access {
    resource_app_id = "00000003-0000-0000-c000-000000000000"

    resource_access {
      id   = "7ab1d382-f21e-4acd-a863-ba3e13f7da61"
      type = "Role"
    }

    resource_access {
      id   = "e1fe6dd8-ba31-4d61-89e7-88639da4683d"
      type = "Scope"
    }

    resource_access {
      id   = "06da0dbc-49e2-44d2-8312-53f166ab848a"
      type = "Scope"
    }
  }

  required_resource_access {
    resource_app_id = "00000002-0000-0000-c000-000000000000"

    resource_access {
      id   = "311a71cc-e848-46a1-bdf8-97ff7156d8e6"
      type = "Scope"
    }
  }

  web {
    homepage_url  = "https://homepage-%[1]d"
    logout_url    = "https://log.me.out"
    redirect_uris = ["https://unittest.hashicorptest.com"]

    implicit_grant {
      access_token_issuance_enabled = true
    }
  }

  owners = [azuread_user.test.object_id]
}
`, data.RandomInteger, data.RandomPassword, data.UUID(), data.UUID(), data.UUID(), data.UUID())
}

func (ApplicationResource) appRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[1]d"

  app_role {
    allowed_member_types = ["User", "Application"]
    description          = "Admins can manage roles and perform all task actions"
    display_name         = "Admin"
    enabled              = true
    id                   = "%[2]s"
    value                = "Admin"
  }
}
`, data.RandomInteger, data.UUID())
}

func (ApplicationResource) appRoleNoValue(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[1]d"

  app_role {
    allowed_member_types = ["User"]
    description          = "Admins can manage roles and perform all task actions"
    display_name         = "Admin"
    id                   = "%[2]s"
    enabled              = true
  }
}
`, data.RandomInteger, data.UUID())
}

func (ApplicationResource) appRolesUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctestApp-%[1]d"

  app_role {
    allowed_member_types = ["User"]
    description          = "Admins can manage roles and perform all task actions"
    display_name         = "Admin"
    enabled              = true
    id                   = "%[2]s"
    value                = ""
  }

  app_role {
    allowed_member_types = ["User"]
    description          = "ReadOnly roles have limited query access"
    display_name         = "ReadOnly"
    enabled              = true
    id                   = "%[3]s"
    value                = "User"
  }
}
`, data.RandomInteger, data.UUID(), data.UUID())
}

func (ApplicationResource) oauth2PermissionScopes(data acceptance.TestData, scopeIDs []string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[1]d"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Allow the application to access acctest-APP-%[1]d on behalf of the signed-in user."
      admin_consent_display_name = "Access acctest-APP-%[1]d"
      enabled                    = true
      id                         = "%[2]s"
      type                       = "User"
      user_consent_description   = "Allow the application to access acctest-APP-%[1]d on your behalf."
      user_consent_display_name  = "Access acctest-APP-%[1]d"
      value                      = "user_impersonation"
    }

    oauth2_permission_scope {
      admin_consent_description  = "Administer the application"
      admin_consent_display_name = "Administer"
      enabled                    = true
      id                         = "%[3]s"
      type                       = "Admin"
      value                      = "administer"
    }
  }
}
`, data.RandomInteger, scopeIDs[0], scopeIDs[1])
}

func (ApplicationResource) oauth2PermissionScopesUpdate(data acceptance.TestData, scopeIDs []string) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[1]d"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Allow the application to access on behalf blah... this changed"
      admin_consent_display_name = "Access acctest-APP-%[1]d"
      enabled                    = true
      id                         = "%[2]s"
      type                       = "User"
      user_consent_description   = "Allow the application to access acctest-APP-%[1]d on your behalf."
      user_consent_display_name  = "Access acctest-APP-%[1]d"
      value                      = "user_impersonation"
    }

    oauth2_permission_scope {
      admin_consent_description  = "Administer the application"
      admin_consent_display_name = "Administer"
      enabled                    = true
      id                         = "%[3]s"
      type                       = "Admin"
      value                      = "administer"
    }

    oauth2_permission_scope {
      admin_consent_description  = "Audit the application"
      admin_consent_display_name = "Audit"
      enabled                    = true
      id                         = "%[4]s"
      type                       = "Admin"
      value                      = "audit"
    }
  }
}
`, data.RandomInteger, scopeIDs[0], scopeIDs[1], scopeIDs[2])
}

func (ApplicationResource) preventDuplicateNamesPass(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name            = "acctest-APP-%[1]d"
  prevent_duplicate_names = true
}
`, data.RandomInteger)
}

func (r ApplicationResource) preventDuplicateNamesFail(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "duplicate" {
  display_name            = azuread_application.test.display_name
  prevent_duplicate_names = true
}
`, r.basic(data))
}

func (ApplicationResource) duplicateAppRolesOauth2PermissionsValues(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[1]d"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Administer the application"
      admin_consent_display_name = "Administer"
      enabled                    = true
      id                         = "%[2]s"
      type                       = "Admin"
      value                      = "administer"
    }
  }

  app_role {
    allowed_member_types = ["User"]
    description          = "Admins can manage roles and perform all task actions"
    display_name         = "Admin"
    enabled              = true
    id                   = "%[3]s"
    value                = "administer"
  }
}
`, data.RandomInteger, data.UUID(), data.UUID())
}

func (ApplicationResource) templateThreeUsers(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azuread" {}

data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "testA" {
  user_principal_name = "acctestUser.%[1]d.A@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-A"
  password            = "%[2]s"
}

resource "azuread_user" "testB" {
  user_principal_name = "acctestUser.%[1]d.B@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-B"
  mail_nickname       = "acctestUser-%[1]d-B"
  password            = "%[2]s"
}

resource "azuread_user" "testC" {
  user_principal_name = "acctestUser.%[1]d.C@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d-C"
  password            = "%[2]s"
}
`, data.RandomInteger, data.RandomPassword)
}

func (r ApplicationResource) singleOwner(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[2]d"
  owners = [
    azuread_user.testA.object_id,
  ]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}

func (r ApplicationResource) threeOwners(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "test" {
  display_name = "acctest-APP-%[2]d"
  owners = [
    azuread_user.testA.object_id,
    azuread_user.testB.object_id,
    azuread_user.testC.object_id,
  ]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}
