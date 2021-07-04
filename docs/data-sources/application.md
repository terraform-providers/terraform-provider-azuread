---
subcategory: "Applications"
---

# Data Source: azuread_application

Use this data source to access information about an existing Application within Azure Active Directory.

## Example Usage

```terraform
data "azuread_application" "example" {
  display_name = "My First AzureAD Application"
}

output "azure_ad_object_id" {
  value = data.azuread_application.example.id
}
```

## Argument Reference

* `application_id` - (Optional) Specifies the Application ID (also called Client ID).
* `display_name` - (Optional) Specifies the display name of the application.
* `object_id` - (Optional) Specifies the Object ID of the application.

~> **NOTE:** One of `object_id`, `application_id` or `display_name` must be specified.

## Attributes Reference

The following attributes are exported:

* `api` - An `api` block as documented below.
* `app_role_ids` - A mapping of app role values to app role IDs, intended to be useful when referencing app roles in other resources in your configuration.
* `app_roles` - A collection of `app_role` blocks as documented below. For more information see [official documentation on Application Roles](https://docs.microsoft.com/en-us/azure/architecture/multitenant-identity/app-roles).
* `application_id` - The Application ID (also called Client ID).
* `device_only_auth_enabled` - Specifies whether this application supports device authentication without a user.
* `disabled_by_microsoft` - Whether Microsoft has disabled the registered application. If the application is disabled, this will be a string indicating the status/reason, e.g. `DisabledDueToViolationOfServicesAgreement`
* `display_name` - The display name for the application.
* `fallback_public_client_enabled` - The fallback application type as public client, such as an installed application running on a mobile device.
* `group_membership_claims` - The `groups` claim issued in a user or OAuth 2.0 access token that the app expects.
* `identifier_uris` - A list of user-defined URI(s) that uniquely identify a Web application within it's Azure AD tenant, or within a verified custom domain if the application is multi-tenant.
* `logo_url` - CDN URL to the application's logo.
* `marketing_url` - URL of the application's marketing page.
* `oauth2_permission_scope_ids` - A mapping of OAuth2.0 permission scope values to scope IDs, intended to be useful when referencing permission scopes in other resources in your configuration.
* `oauth2_post_response_required` - Specifies whether, as part of OAuth 2.0 token requests, Azure AD allows POST requests, as opposed to GET requests. When `false`, only GET requests are allowed.
* `object_id` - The application's object ID.
* `optional_claims` - An `optional_claims` block as documented below.
* `owners` - A list of object IDs of principals that are assigned ownership of the application.
* `privacy_statement_url` - URL of the application's privacy statement.
* `public_client` - A `public_client` block as documented below.
* `publisher_domain` - The verified publisher domain for the application.
* `required_resource_access` - A collection of `required_resource_access` blocks as documented below.
* `sign_in_audience` - The Microsoft account types that are supported for the current application. One of `AzureADMyOrg`, `AzureADMultipleOrgs`, `AzureADandPersonalMicrosoftAccount` or `PersonalMicrosoftAccount`.
* `single_page_application` - A `single_page_application` block as documented below.
* `support_url` - URL of the application's support page.
* `terms_of_service_url` - URL of the application's terms of service statement.
* `web` - A `web` block as documented below.

---

`api` block exports the following:

* `known_client_applications` - A set of application IDs (client IDs), used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app.
* `mapped_claims_enabled` - Allows an application to use claims mapping without specifying a custom signing key.
* `oauth2_permission_scope` - One or more `oauth2_permission_scope` blocks as documented below, to describe delegated permissions exposed by the web API represented by this application.
* `requested_access_token_version` - The access token version expected by this resource. Possible values are `1` or `2`.

---

`oauth2_permission_scope` block exports the following:

* `admin_consent_description` - Delegated permission description that appears in all tenant-wide admin consent experiences, intended to be read by an administrator granting the permission on behalf of all users.
* `admin_consent_display_name` - Display name for the delegated permission, intended to be read by an administrator granting the permission on behalf of all users.
* `enabled` - Determines if the permission scope is enabled.
* `id` - The unique identifier of the delegated permission. Must be a valid UUID.
* `type` - Whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator should be required for consent to the permissions. Possible values are `User` or `Admin`.
* `user_consent_description` - Delegated permission description that appears in the end user consent experience, intended to be read by a user consenting on their own behalf.
* `user_consent_display_name` - Display name for the delegated permission that appears in the end user consent experience.
* `value` - The value that is used for the `scp` claim in OAuth 2.0 access tokens.

---

`app_role` block exports the following:

* `allowed_member_types` - Specifies whether this app role definition can be assigned to users and groups, or to other applications (that are accessing this application in a standalone scenario). Possible values are `User` or `Application`, or both.
* `description` - Description of the app role that appears when the role is being assigned and, if the role functions as an application permissions, during the consent experiences.
* `display_name` - Display name for the app role that appears during app role assignment and in consent experiences.
* `enabled` - Determines if the app role is enabled.
* `id` - The unique identifier of the app role.
* `value` - The value that is used for the `roles` claim in ID tokens and OAuth 2.0 access tokens that are authenticating an assigned service or user principal.

---

`optional_claims` block exports the following:

* `access_token` - One or more `access_token` blocks as documented below.
* `id_token` - One or more `id_token` blocks as documented below.
* `saml2_token` - One or more `saml2_token` blocks as documented below.

---

`access_token`, `id_token` and `saml2_token` blocks export the following:

* `additional_properties` - List of Additional Properties of the claim. If a property exists in this list, it modifies the behaviour of the optional claim.
* `essential` - Whether the claim specified by the client is necessary to ensure a smooth authorization experience.
* `name` - The name of the optional claim.
* `source` - The source of the claim. If `source` is absent, the claim is a predefined optional claim. If `source` is `user`, the value of `name` is the extension property from the user object.

---

`public_client` block exports the following:

* `redirect_uris` - A list of URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.

---

`required_resource_access` block exports the following:

* `resource_access` - A collection of `resource_access` blocks as documented below, describing OAuth2.0 permission scopes and app roles that the application requires from the specified resource.
* `resource_app_id` - The unique identifier for the resource that the application requires access to. This is the Application ID of the target application.

---

`resource_access` block exports the following:

* `id` - The unique identifier for an app role or OAuth2 permission scope published by the resource application.
* `type` - Specifies whether the `id` property references an app role or an OAuth2 permission scope. Possible values are `Role` or `Scope`.

---

`single_page_application` block exports the following:

* `redirect_uris` - A list of URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.

---

`web` block exports the following:

* `homepage_url` - Home page or landing page of the application.
* `implicit_grant` - An `implicit_grant` block as documented above.
* `logout_url` - The URL that will be used by Microsoft's authorization service to sign out a user using front-channel, back-channel or SAML logout protocols.
* `redirect_uris` - A list of URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.

---

`implicit_grant` block exports the following:

* `access_token_issuance_enabled` - Whether this web application can request an access token using OAuth 2.0 implicit flow.
* `id_token_issuance_enabled` - Whether this web application can request an ID token using OAuth 2.0 implicit flow.
