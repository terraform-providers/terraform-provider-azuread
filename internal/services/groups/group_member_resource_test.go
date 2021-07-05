package groups_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/groups/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type GroupMemberResource struct{}

func TestAccGroupMember_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group_member", "test")
	r := GroupMemberResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.group(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("group_object_id").IsUuid(),
				check.That(data.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroupMember_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group_member", "test")
	r := GroupMemberResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.servicePrincipal(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("group_object_id").IsUuid(),
				check.That(data.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroupMember_user(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group_member", "testA")
	r := GroupMemberResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oneUser(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("group_object_id").IsUuid(),
				check.That(data.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroupMember_multipleUser(t *testing.T) {
	dataA := acceptance.BuildTestData(t, "azuread_group_member", "testA")
	dataB := acceptance.BuildTestData(t, "azuread_group_member", "testB")
	r := GroupMemberResource{}

	dataA.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.oneUser(dataA),
			Check: resource.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataA.ResourceName).Key("group_object_id").IsUuid(),
				check.That(dataA.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		dataA.ImportStep(),
		{
			Config: r.twoUsers(dataA),
			Check: resource.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataA.ResourceName).Key("group_object_id").IsUuid(),
				check.That(dataA.ResourceName).Key("member_object_id").IsUuid(),
				check.That(dataB.ResourceName).ExistsInAzure(r),
				check.That(dataB.ResourceName).Key("group_object_id").IsUuid(),
				check.That(dataB.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		// we rerun the config so the group resource updates with the number of members
		{
			Config: r.twoUsers(dataA),
			Check: resource.ComposeTestCheckFunc(
				check.That("azuread_group.test").Key("members.#").HasValue("2"),
			),
		},
		dataA.ImportStep(),
		{
			Config: r.oneUser(dataA),
			Check: resource.ComposeTestCheckFunc(
				check.That(dataA.ResourceName).ExistsInAzure(r),
				check.That(dataA.ResourceName).Key("group_object_id").IsUuid(),
				check.That(dataA.ResourceName).Key("member_object_id").IsUuid(),
			),
		},
		// we rerun the config so the group resource updates with the number of members
		{
			Config: r.oneUser(dataA),
			Check: resource.ComposeTestCheckFunc(
				check.That("azuread_group.test").Key("members.#").HasValue("1"),
			),
		},
	})
}

func TestAccGroupMember_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_group_member", "test")
	r := GroupMemberResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.group(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport(data)),
	})
}

func (r GroupMemberResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Groups.GroupsClient
	client.BaseClient.DisableRetries = true

	id, err := parse.GroupMemberID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing Group Member ID: %v", err)
	}

	members, _, err := client.ListMembers(ctx, id.GroupId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Group members (groupId: %q): %+v", id.GroupId, err)
	}

	if members != nil {
		for _, objectId := range *members {
			if strings.EqualFold(objectId, id.MemberId) {
				return utils.Bool(true), nil
			}
		}
	}

	return nil, fmt.Errorf("Member %q was not found in Group %q", id.MemberId, id.GroupId)
}

func (GroupMemberResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_group" "test" {
  display_name     = "acctestGroup-%[1]d"
  security_enabled = true
}
`, data.RandomInteger)
}

func (GroupMemberResource) templateThreeUsers(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

func (r GroupMemberResource) group(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "member" {
  display_name     = "acctestGroup-%[2]d-Member"
  security_enabled = true
}

resource "azuread_group_member" "test" {
  group_object_id  = azuread_group.test.object_id
  member_object_id = azuread_group.member.object_id
}
`, r.template(data), data.RandomInteger)
}

func (r GroupMemberResource) servicePrincipal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_application" "test" {
  display_name = "acctestServicePrincipal-%[2]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azuread_group_member" "test" {
  group_object_id  = azuread_group.test.object_id
  member_object_id = azuread_service_principal.test.object_id
}
`, r.template(data), data.RandomInteger)
}

func (r GroupMemberResource) oneUser(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azuread_group_member" "testA" {
  group_object_id  = azuread_group.test.object_id
  member_object_id = azuread_user.testA.object_id
}
`, r.template(data), r.templateThreeUsers(data))
}

func (r GroupMemberResource) twoUsers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "azuread_group_member" "testA" {
  group_object_id  = azuread_group.test.object_id
  member_object_id = azuread_user.testA.object_id
}

resource "azuread_group_member" "testB" {
  group_object_id  = azuread_group.test.object_id
  member_object_id = azuread_user.testB.object_id
}
`, r.template(data), r.templateThreeUsers(data))
}

func (r GroupMemberResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group_member" "import" {
  group_object_id  = azuread_group_member.test.group_object_id
  member_object_id = azuread_group_member.test.member_object_id
}
`, r.group(data))
}
