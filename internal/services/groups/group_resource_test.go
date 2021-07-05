package groups_test

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

type GroupResource struct{}

func TestAccGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.unified(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.unified(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_owners(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withThreeOwners(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_members(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withThreeMembers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_membersAndOwners(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withOwnersAndMembers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_manyMembersAndOwners(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withManyOwnersAndMembers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_membersDiverse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDiverseMembers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_ownersDiverse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withDiverseOwners(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_membersUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withOneMember(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withThreeMembers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withServicePrincipalMember(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.noMembers(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_ownersUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withThreeOwners(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withOneOwner(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withServicePrincipalOwner(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_preventDuplicateNamesPass(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.preventDuplicateNamesPass(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctestGroup-%d", data.RandomInteger)),
			),
		},
		data.ImportStep("prevent_duplicate_names"),
	})
}

func TestAccGroup_preventDuplicateNamesFail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		data.RequiresImportErrorStep(r.preventDuplicateNamesFail(data)),
	})
}

func (r GroupResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Groups.GroupsClient
	client.BaseClient.DisableRetries = true

	group, status, err := client.Get(ctx, state.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Group with object ID %q does not exist", state.ID)
		}
		return nil, fmt.Errorf("failed to retrieve Group with object ID %q: %+v", state.ID, err)
	}
	return utils.Bool(group.ID != nil && *group.ID == state.ID), nil
}

func (GroupResource) templateDiverseDirectoryObjects(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_application" "test" {
  display_name = "acctestGroup-%[1]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_user" "test" {
  user_principal_name = "acctestGroup.%[1]d@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestGroup-%[1]d"
  password            = "%[2]s"
}

resource "azuread_group" "member" {
  display_name     = "acctestGroup-%[1]d-Member"
  security_enabled = true
}
`, data.RandomInteger, data.RandomPassword)
}

func (GroupResource) templateThreeUsers(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "testA" {
  user_principal_name = "acctestGroup.%[1]d.A@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestGroup-%[1]d-A"
  password            = "%[2]s"
}

resource "azuread_user" "testB" {
  user_principal_name = "acctestGroup.%[1]d.B@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestGroup-%[1]d-B"
  mail_nickname       = "acctestGroup-%[1]d-B"
  password            = "%[2]s"
}

resource "azuread_user" "testC" {
  user_principal_name = "acctestGroup.%[1]d.C@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestGroup-%[1]d-C"
  password            = "%[2]s"
}
`, data.RandomInteger, data.RandomPassword)
}

func (GroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
}
`, data.RandomInteger)
}

func (GroupResource) unified(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  types            = ["Unified"]
  mail_enabled     = true
  security_enabled = true
}
`, data.RandomInteger)
}

func (GroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  user_principal_name = "acctestGroup.%[1]d@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestGroup-%[1]d"
  password            = "%[2]s"
}

resource "azuread_group" "test" {
  description      = "Please delete me as this is a.test.AD group!"
  display_name     = "acctestGroup-complete-%[1]d"
  types            = ["Unified"]
  mail_enabled     = true
  security_enabled = true
  members          = [azuread_user.test.object_id]
  owners           = [azuread_user.test.object_id]
}
`, data.RandomInteger, data.RandomPassword)
}

func (GroupResource) noMembers(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
  members          = []
}
`, data.RandomInteger)
}

func (r GroupResource) withDiverseMembers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  members          = [azuread_user.test.object_id, azuread_group.member.object_id, azuread_service_principal.test.object_id]
}
`, r.templateDiverseDirectoryObjects(data), data.RandomInteger)
}

func (r GroupResource) withDiverseOwners(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  owners           = [azuread_user.test.object_id, azuread_service_principal.test.object_id]
}
`, r.templateDiverseDirectoryObjects(data), data.RandomInteger)
}

func (r GroupResource) withOneMember(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  members          = [azuread_user.testA.object_id]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}

func (r GroupResource) withOneOwner(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  owners           = [azuread_user.testA.object_id]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}

func (r GroupResource) withThreeMembers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  members          = [azuread_user.testA.object_id, azuread_user.testB.object_id, azuread_user.testC.object_id]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}

func (r GroupResource) withThreeOwners(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  owners           = [azuread_user.testA.object_id, azuread_user.testB.object_id, azuread_user.testC.object_id]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}

func (r GroupResource) withOwnersAndMembers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[2]d"
  security_enabled = true
  owners           = [azuread_user.testA.object_id]
  members          = [azuread_user.testB.object_id, azuread_user.testC.object_id]
}
`, r.templateThreeUsers(data), data.RandomInteger)
}

func (r GroupResource) withManyOwnersAndMembers(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_domains" "test" {
  only_initial = true
}

resource "azuread_user" "test" {
  count = 25

  user_principal_name = "acctestGroupParticipant${count.index}-%[1]d@${data.azuread_domains.test.domains.0.domain_name}"
  display_name        = "acctestGroupParticipant${count.index}-%[1]d"
  password            = "Qwer5678!@#"
}

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
  owners           = azuread_user.test.*.object_id
  members          = azuread_user.test.*.object_id
}
`, data.RandomInteger)
}

func (GroupResource) withServicePrincipalMember(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_application" "test" {
  display_name = "acctestGroup-%[1]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
  members          = [azuread_service_principal.test.object_id]
}
`, data.RandomInteger)
}

func (GroupResource) withServicePrincipalOwner(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_application" "test" {
  display_name = "acctestGroup-%[1]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
  owners           = [azuread_service_principal.test.object_id]
}
`, data.RandomInteger)
}

func (GroupResource) preventDuplicateNamesPass(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_group" "test" {
  display_name            = "acctestGroup-%[1]d"
  security_enabled        = true
  prevent_duplicate_names = true
}
`, data.RandomInteger)
}

func (r GroupResource) preventDuplicateNamesFail(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "duplicate" {
  display_name            = azuread_group.test.display_name
  security_enabled        = true
  prevent_duplicate_names = true
}
`, r.basic(data))
}
