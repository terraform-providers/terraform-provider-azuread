package applications_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
)

type ApplicationDataSource struct{}

func TestAccApplicationDataSource_byObjectId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_application", "test")
	r := ApplicationDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.objectId(data),
			Check:  r.testCheck(data),
		},
	})
}

func TestAccApplicationDataSource_byApplicationId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_application", "test")
	r := ApplicationDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.applicationId(data),
			Check:  r.testCheck(data),
		},
	})
}

func TestAccApplicationDataSource_byDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_application", "test")
	r := ApplicationDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.displayName(data),
			Check:  r.testCheck(data),
		},
	})
}

func (ApplicationDataSource) testCheck(data acceptance.TestData) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		check.That(data.ResourceName).Key("application_id").IsUuid(),
		check.That(data.ResourceName).Key("object_id").IsUuid(),
		check.That(data.ResourceName).Key("api.0.oauth2_permission_scopes.#").HasValue("2"),
		check.That(data.ResourceName).Key("app_roles.#").HasValue("2"),
		check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-APP-complete-%d", data.RandomInteger)),
		check.That(data.ResourceName).Key("group_membership_claims.#").HasValue("1"),
		check.That(data.ResourceName).Key("group_membership_claims.0").HasValue("All"),
		check.That(data.ResourceName).Key("identifier_uris.#").HasValue("1"),
		check.That(data.ResourceName).Key("identifier_uris.0").HasValue(fmt.Sprintf("api://hashicorptestapp-%d", data.RandomInteger)),
		check.That(data.ResourceName).Key("optional_claims.#").HasValue("1"),
		check.That(data.ResourceName).Key("optional_claims.0.access_token.#").HasValue("2"),
		check.That(data.ResourceName).Key("optional_claims.0.id_token.#").HasValue("1"),
		check.That(data.ResourceName).Key("required_resource_access.#").HasValue("2"),
		check.That(data.ResourceName).Key("sign_in_audience").HasValue("AzureADMultipleOrgs"),
		check.That(data.ResourceName).Key("web.0.homepage_url").HasValue(fmt.Sprintf("https://homepage-%d", data.RandomInteger)),
		check.That(data.ResourceName).Key("web.0.logout_url").HasValue("https://log.me.out"),
		check.That(data.ResourceName).Key("web.0.redirect_uris.#").HasValue("1"),
		check.That(data.ResourceName).Key("web.0.redirect_uris.0").HasValue("https://unittest.hashicorptest.com"),
	)
}

func (ApplicationDataSource) objectId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_application" "test" {
  object_id = azuread_application.test.object_id
}
`, ApplicationResource{}.complete(data))
}

func (ApplicationDataSource) applicationId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_application" "test" {
  application_id = azuread_application.test.application_id
}
`, ApplicationResource{}.complete(data))
}

func (ApplicationDataSource) displayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_application" "test" {
  display_name = azuread_application.test.display_name
}
`, ApplicationResource{}.complete(data))
}
