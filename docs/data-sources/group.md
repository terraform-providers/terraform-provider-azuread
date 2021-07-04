---
subcategory: "Groups"
---

# Data Source: azuread_group

Gets information about an Azure Active Directory group.

## Example Usage (by Group Display Name)

```terraform
data "azuread_group" "example" {
  display_name     = "MyGroupName"
  security_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name for the group.
* `mail_enabled` - (Optional) Whether the group is mail-enabled.
* `object_id` - (Optional) Specifies the object ID of the group.
* `security_enabled` - (Optional) Whether the group is a security group.

~> **NOTE:** One of `display_name` or `object_id` must be specified.

## Attributes Reference

The following attributes are exported:

* `assignable_to_role` - Indicates whether this group can be assigned to an Azure Active Directory role.
* `description` - The optional description of the group.
* `display_name` - The display name for the group.
* `object_id` - The object ID of the group.
* `mail_enabled` - Whether the group is mail-enabled.
* `members` - The object IDs of the group members.
* `owners` - The object IDs of the group owners.
* `security_enabled` - Whether the group is a security group.
* `types` - A list of group types configured for the group. The only supported type is `Unified`, which specifies a Microsoft 365 group.
