---
subcategory: "Service Principals"
---

# Data Source: azuread_service_principal

Gets information about an existing service principal associated with an application within Azure Active Directory.

## Example Usage (by Application Display Name)

```terraform
data "azuread_service_principal" "example" {
  display_name = "my-awesome-application"
}
```

## Example Usage (by Application ID)

```terraform
data "azuread_service_principal" "example" {
  application_id = "00000000-0000-0000-0000-000000000000"
}
```

## Example Usage (by Object ID)

```terraform
data "azuread_service_principal" "example" {
  object_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional) The application ID (client ID) of the application associated with this service principal.
* `display_name` - (Optional) The display name of the application associated with this service principal.
* `object_id` - (Optional) The object ID of the service principal.

~> **NOTE:** At least one of `application_id`, `display_name` or `object_id` must be specified.

## Attributes Reference

The following attributes are exported:

* `app_roles` - A collection of `app_roles` blocks as documented below. For more information [official documentation](https://docs.microsoft.com/en-us/azure/architecture/multitenant-identity/app-roles).
* `object_id` - The object ID for the service principal.
* `oauth2_permission_scopes` - A collection of OAuth 2.0 delegated permissions exposed by the associated application. Each permission is covered by an `oauth2_permission_scopes` block as documented below.

---

`app_roles` block exports the following:

* `allowed_member_types` - Specifies whether this app role definition can be assigned to users and groups, or to other applications (that are accessing this application in daemon service scenarios). Possible values are: `User` and `Application`, or both.
* `description` - Permission help text that appears in the admin app assignment and consent experiences.
* `display_name` - Display name for the permission that appears in the admin consent and app assignment experiences.
* `id` - The unique identifier of the app role.
* `is_enabled` - Determines if the app role is enabled.
* `value` - Specifies the value of the roles claim that the application should expect in the authentication and access tokens.

---

`oauth2_permission_scopes` block exports the following:

* `admin_consent_description` - Delegated permission description that appears in all tenant-wide admin consent experiences, intended to be read by an administrator granting the permission on behalf of all users.
* `admin_consent_display_name` - Display name for the delegated permission, intended to be read by an administrator granting the permission on behalf of all users.
* `enabled` - Determines if the permission scope is enabled.
* `id` - The unique identifier of the delegated permission. Must be a valid UUID.
* `type` - Whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator should be required for consent to the permissions. Possible values are `User` or `Admin`.
* `user_consent_description` - Delegated permission description that appears in the end user consent experience, intended to be read by a user consenting on their own behalf.
* `user_consent_display_name` - Display name for the delegated permission that appears in the end user consent experience.
* `value` - The value that is used for the `scp` claim in OAuth 2.0 access tokens.
