package groups_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
)

type GroupDataSource struct{}

func TestAccGroupDataSource_byDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.displayName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
	})
}
func TestAccGroupDataSource_byDisplayNameWithSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.displayNameSecurity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccGroupDataSource_byCaseInsensitiveDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.caseInsensitiveDisplayName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccGroupDataSource_byObjectId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.objectId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccGroupDataSource_byObjectIdWithSecurity(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.objectIdSecurity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
	})
}

func TestAccGroupDataSource_members(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.members(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("members.#").HasValue("3"),
			),
		},
	})
}

func TestAccGroupDataSource_owners(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_group", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: GroupDataSource{}.owners(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("owners.#").HasValue("3"),
			),
		},
	})
}

func (GroupDataSource) displayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_group" "test" {
  display_name = azuread_group.test.display_name
}
`, GroupResource{}.basic(data))
}

func (GroupDataSource) displayNameSecurity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_group" "test" {
  display_name     = azuread_group.test.display_name
  mail_enabled     = false
  security_enabled = true
}
`, GroupResource{}.basic(data))
}

func (GroupDataSource) caseInsensitiveDisplayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
data "azuread_group" "test" {
  display_name = upper(azuread_group.test.display_name)
}
`, GroupResource{}.basic(data))
}

func (GroupDataSource) objectId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_group" "test" {
  object_id = azuread_group.test.object_id
}
`, GroupResource{}.basic(data))
}

func (GroupDataSource) objectIdSecurity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_group" "test" {
  object_id        = azuread_group.test.object_id
  mail_enabled     = false
  security_enabled = true
}
`, GroupResource{}.basic(data))
}

func (GroupDataSource) members(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_group" "test" {
  object_id = azuread_group.test.object_id
}
`, GroupResource{}.withThreeMembers(data))
}

func (GroupDataSource) owners(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azuread_group" "test" {
  object_id = azuread_group.test.object_id
}
`, GroupResource{}.withThreeOwners(data))
}
