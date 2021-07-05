---
subcategory: "Users"
---

# Data Source: azuread_users

Gets object IDs or user principal names for multiple Azure Active Directory users.

## Example Usage

```terraform
data "azuread_users" "users" {
  user_principal_names = ["kat@hashicorp.com", "byte@hashicorp.com"]
}
```

## Argument Reference

The following arguments are supported:

* `mail_nicknames` - (Optional) The email aliases of the users.
* `ignore_missing` - (Optional) Ignore missing users and return users that were found. The data source will still fail if no users are found. Defaults to false.
* `object_ids` - (Optional) The object IDs of the users.
* `user_principal_names` - (Optional) The user principal names (UPNs) of the users.

~> **NOTE:** One of `user_principal_names`, `object_ids` or `mail_nicknames` must be specified. These _may_ be specified as an empty list, in which case no results will be returned.

## Attributes Reference

The following attributes are exported:

* `mail_nicknames` - The email aliases of the users.
* `object_ids` - The object IDs of the users.
* `user_principal_names` - The user principal names (UPNs) of the users.
* `users` - A list of users. Each `user` object provides the attributes documented below.

___

`user` object exports the following:

* `account_enabled` - Whether or not the account is enabled.
* `display_name` - The display name of the user.
* `mail_nickname` - The email alias of the user.
* `mail` - The primary email address of the user.
* `object_id` - The object ID of the user.
* `onpremises_immutable_id` - The value used to associate an on-premises Active Directory user account with their Azure AD user object.
* `onpremises_sam_account_name` - The on-premise SAM account name of the user.
* `onpremises_user_principal_name` - The on-premise user principal name of the user.
* `usage_location` - The usage location of the user.
* `user_principal_name` - The user principal name (UPN) of the user.
