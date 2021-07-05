---
subcategory: "Groups"
---

# Resource: azuread_group_member

Manages a single group membership within Azure Active Directory.

-> **Warning** Do not use this resource at the same time as the `members` property of the `azuread_group` resource.

## Example Usage

```terraform

data "azuread_user" "example" {
  user_principal_name = "jdoe@hashicorp.com"
}

resource "azuread_group" "example" {
  display_name     = "my_group"
  security_enabled = true
}

resource "azuread_group_member" "example" {
  group_object_id  = azuread_group.example.id
  member_object_id = data.azuread_user.example.id
}
```

## Argument Reference

The following arguments are supported:

* `group_object_id` - (Required) The object ID of the group you want to add the member to. Changing this forces a new resource to be created.
* `member_object_id` - (Required) The object ID of the principal you want to add as a member to the group. Supported object types are Users, Groups or Service Principals. Changing this forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*No additional attributes are exported*

## Import

Group members can be imported using the object ID of the group and the object ID of the member, e.g.

```shell
terraform import azuread_group_member.test 00000000-0000-0000-0000-000000000000/member/11111111-1111-1111-1111-111111111111
```

-> **NOTE:** This ID format is unique to Terraform and is composed of the Azure AD Group Object ID and the target Member Object ID in the format `{GroupObjectID}/member/{MemberObjectID}`.
