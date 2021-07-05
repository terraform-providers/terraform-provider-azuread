package serviceprincipals_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
)

type ServicePrincipalDataSource struct{}

func TestAccServicePrincipalDataSource_byApplicationId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_service_principal", "test")
	r := ServicePrincipalDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.byApplicationId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("app_roles.#").HasValue("2"),
				check.That(data.ResourceName).Key("oauth2_permission_scopes.#").HasValue("2"),
			),
		},
	})
}

func TestAccServicePrincipalDataSource_byDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_service_principal", "test")
	r := ServicePrincipalDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.byDisplayName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("app_roles.#").HasValue("2"),
				check.That(data.ResourceName).Key("oauth2_permission_scopes.#").HasValue("2"),
			),
		},
	})
}

func TestAccServicePrincipalDataSource_byObjectId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_service_principal", "test")
	r := ServicePrincipalDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.byObjectId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("application_id").Exists(),
				check.That(data.ResourceName).Key("object_id").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("app_roles.#").HasValue("2"),
				check.That(data.ResourceName).Key("oauth2_permission_scopes.#").HasValue("2"),
			),
		},
	})
}

func (ServicePrincipalDataSource) byApplicationId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_service_principal" "test" {
  application_id = azuread_service_principal.test.application_id
}
`, ServicePrincipalResource{}.complete(data))
}

func (ServicePrincipalDataSource) byDisplayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_service_principal" "test" {
  display_name = azuread_service_principal.test.display_name
}
`, ServicePrincipalResource{}.complete(data))
}

func (ServicePrincipalDataSource) byObjectId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_service_principal" "test" {
  object_id = azuread_service_principal.test.object_id
}
`, ServicePrincipalResource{}.complete(data))
}
