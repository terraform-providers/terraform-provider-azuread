---
subcategory: "Service Principals"
---

# Resource: azuread_service_principal

Manages a service principal associated with an application within Azure Active Directory.

## Example Usage

```terraform
resource "azuread_application" "example" {
  display_name = "example"
}

resource "azuread_service_principal" "example" {
  application_id               = azuread_application.example.application_id
  app_role_assignment_required = false

  tags = ["example", "tags", "here"]
}
```

## Argument Reference

The following arguments are supported:

* `app_role_assignment_required` - (Optional) Whether this service principal requires an app role assignment to a user or group before Azure AD will issue a user or access token to the application. Defaults to `false`.
* `application_id` - (Required) The application ID (client ID) of the application for which to create a service principal.
* `tags` - (Optional) A set of tags to apply to the service principal.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `app_roles` - A list of app roles published b the associated application, as documented below. For more information [official documentation](https://docs.microsoft.com/en-us/azure/architecture/multitenant-identity/app-roles).
* `display_name` - The display name of the application associated with this service principal.
* `oauth2_permission_scopes` - A list of OAuth 2.0 delegated permission scopes published by the associated application, as documented below.
* `object_id` - The object ID of the service principal.

---

`app_roles` is a list of objects with the following attributes:

* `allowed_member_types` - Specifies whether this app role definition can be assigned to users and groups, or to other applications (that are accessing this application in a standalone scenario). Possible values are: `User` and `Application`, or both.
* `description` - Description of the app role that appears when the role is being assigned and, if the role functions as an application permissions, during the consent experiences.
* `display_name` - Display name for the app role that appears during app role assignment and in consent experiences.
* `enabled` - Determines if the app role is enabled.
* `id` - The unique identifier of the `app_role`.
* `value` - The value that is used for the `roles` claim in ID tokens and OAuth 2.0 access tokens that are authenticating an assigned service or user principal.

---

`oauth2_permission_scopes` is a list of objects with the following attributes:

* `admin_consent_description` - Delegated permission description that appears in all tenant-wide admin consent experiences, intended to be read by an administrator granting the permission on behalf of all users.
* `admin_consent_display_name` - Display name for the delegated permission, intended to be read by an administrator granting the permission on behalf of all users.
* `enabled` - Specifies whether the permission scope is enabled.
* `id` - The unique identifier of the delegated permission.
* `type` - Whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator should be required for consent to the permissions. Possible values are `User` or `Admin`.
* `user_consent_description` - Delegated permission description that appears in the end user consent experience, intended to be read by a user consenting on their own behalf.
* `user_consent_display_name` - Display name for the delegated permission that appears in the end user consent experience.
* `value` - The value that is used for the `scp` claim in OAuth 2.0 access tokens.

## Import

Service principals can be imported using their object ID, e.g.

```shell
terraform import azuread_service_principal.test 00000000-0000-0000-0000-000000000000
```
