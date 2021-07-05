---
subcategory: "Groups"
---

# Resource: azuread_group

Manages a group within Azure Active Directory.

## Example Usage

*Basic example*

```terraform
resource "azuread_group" "example" {
  display_name     = "example"
  security_enabled = true
}
```

*Microsoft 365 group*

```terraform
resource "azuread_group" "example" {
  display_name     = "example"
  mail_enabled     = true
  mail_nickname    = "ExampleGroup"
  security_enabled = true
  types            = ["Unified"]
}
```

*Group with members*

```terraform
resource "azuread_user" "example" {
  display_name        = "J Doe"
  password            = "notSecure123"
  user_principal_name = "jdoe@hashicorp.com"
}

resource "azuread_group" "example" {
  display_name     = "MyGroup"
  security_enabled = true

  members = [
    azuread_user.example.object_id,
    /* more users */
  ]
}
```

## Argument Reference

The following arguments are supported:

* `assignable_to_role` - (Optional) Indicates whether this group can be assigned to an Azure Active Directory role. Can only be `true` for security-enabled groups. Changing this forces a new resource to be created.
* `behaviors` - (Optional) A set of behaviors for a Microsoft 365 group. Possible values are `AllowOnlyMembersToPost`, `HideGroupInOutlook`, `SubscribeNewGroupMembers` and `WelcomeEmailDisabled`. See [official documentation](https://docs.microsoft.com/en-us/graph/group-set-options) for more details. Changing this forces a new resource to be created.
* `description` - (Optional) The description for the group.
* `display_name` - (Required) The display name for the group.
* `mail_enabled` - (Optional) Whether the group is a mail enabled, with a shared group mailbox. At least one of `mail_enabled` or `security_enabled` must be specified. A group can be mail enabled _and_ security enabled.
* `mail_nickname` - (Optional) The mail alias for the group, unique in the organisation. Required for mail-enabled groups. Changing this forces a new resource to be created.
* `members` - (Optional) A set of members who should be present in this group. Supported object types are Users, Groups or Service Principals.
* `owners` - (Optional) A set of owners who own this group. Supported object types are Users or Service Principals.
* `prevent_duplicate_names` - (Optional) If `true`, will return an error if an existing group is found with the same name. Defaults to `false`.
* `provisioning_options` - (Optional) A set of provisioning options for a Microsoft 365 group. The only supported value is `Team`. See [official documentation](https://docs.microsoft.com/en-us/graph/group-set-options) for details. Changing this forces a new resource to be created.
* `security_enabled` - (Optional) Whether the group is a security group for controlling access to in-app resources. At least one of `security_enabled` or `mail_enabled` must be specified. A group can be security enabled _and_ mail enabled.
* `theme` - (Optional) The colour theme for a Microsoft 365 group. Possible values are `Blue`, `Green`, `Orange`, `Pink`, `Purple`, `Red` or `Teal`. By default, no theme is set.
* `types` - (Optional) A set of group types to configure for the group. The only supported type is `Unified`, which specifies a Microsoft 365 group. Required when `mail_enabled` is true. Changing this forces a new resource to be created.
* `visibility` - (Optional) The group join policy and group content visibility. Possible values are `Private`, `Public`, or `Hiddenmembership`. Only Microsoft 365 groups can have `Hiddenmembership` visibility and this value must be set when the group is created. By default, security groups will receive `Private` visibility and Microsoft 365 groups will receive `Public` visibility.

-> **Group Name Uniqueness** Group names are not unique within Azure Active Directory. Use the `prevent_duplicate_names` argument to check for existing groups if you want to avoid name collisions.

!> **Warning** Do not use the `azuread_group_member` resource at the same time as the `members` argument.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mail` - The SMTP address for the group.
* `object_id` - The object ID of the group.
* `onpremises_domain_name` - The on-premises FQDN, also called dnsDomainName, synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_netbios_name` - The on-premises NetBIOS name, synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_sam_account_name` - The on-premises SAM account name, synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_security_identifier` - The on-premises security identifier (SID), synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_sync_enabled` - Whether this group is synchronised from an on-premises directory (`true`), no longer synchronised (`false`), or has never been synchronised (`null`).
* `preferred_language` - The preferred language for a Microsoft 365 group, in ISO 639-1 notation.
* `proxy_addresses` - List of email addresses for the group that direct to the same group mailbox.

## Import

Groups can be imported using their object ID, e.g.

```shell
terraform import azuread_group.my_group 00000000-0000-0000-0000-000000000000
```
