---
page_title: "AzureAD v2.0 and Microsoft Graph"
subcategory: "Upgrade Guides"
---

# AzureAD v2.0 and Microsoft Graph

From version 2.0 the AzureAD provider exclusively uses [Microsoft Graph](https://docs.microsoft.com/en-us/graph/overview) to connect to Azure Active Directory and has ceased to support using the [Azure Active Directory Graph API](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-graph-api).

Due to differences between the two APIs, some schema deprecations have already been introduced prior to v2.0 and several fields have been renamed, removed or otherwise changed in v2.0. Please consult the following guide to determine any configuration changes you will need to make in order to upgrade to version 2.0.

We take semantic versioning very seriously which is why these changes have been introduced in a new major version of the provider. You can read up the background behind the API migration and proposed development path in [this GitHub issue](https://github.com/hashicorp/terraform-provider-azuread/issues/323), and follow along with the work being done for v2.0 in the [GitHub milestone](https://github.com/hashicorp/terraform-provider-azuread/milestone/16)

In version 1.5.0 or later of the AzureAD provider, limited beta support for Microsoft Graph can be enabled. See [Beta support for Microsoft Graph](#beta-support-for-microsoft-graph-in-v150) for more details.

## Pinning your provider version

We recommend pinning the version of each provider you use in Terraform. You can do this using the `version` attribute in the `required_providers` block of your `terraform` configuration block.

To pin to a specific version of the AzureAD provider:

```hcl
terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "= 1.5.1"
    }
  }
}
```

To pin to any 1.5 release:

```hcl
terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 1.5.0"
    }
  }
}
```

Older versions of Terraform (0.12.x) can pin the provider version using the `version` attribute within the `provider` block:
```hcl
provider "azuread" {
  version = "~> 1.5.0"
}
```

This will enable you to upgrade to version 2.0 at your convenience, by simply advancing the desired target version in your configuration. See the [Lock and Upgrade Provider Versions](https://learn.hashicorp.com/tutorials/terraform/provider-versioning) guide on HashiCorp Learn for more details.

## Authentication and National Clouds

Existing authentication methods will continue to work unchanged, whether you authenticate with a service principal ([client certificate](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_client_certificate) or [client secret](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_client_secret)), [managed identity](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/managed_service_identity), or using [Azure CLI](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/azure_cli).

However, you may need to assign new API permissions depending on your configuration and authentication scenario.

For users connecting to national clouds (e.g. germany, china and usgovernment), these are all supported using the existing provider configuration property `environment`, or the environment variable `ARM_ENVIRONMENT`. The "usgovernment" environment has been split into two environments "usgovernmentl4" and "usgovernmentl5" - see [this post](https://developer.microsoft.com/en-us/office/blogs/new-microsoft-graph-endpoints-in-us-government-cloud/) for more information. Specifying the "usgovernment" environment will use the "usgovernmentl4" cloud.

## New API permissions

Microsoft Graph is a different web service to Azure Active Directory Graph, and as such if you are authenticating using a service principal, you may need to assign new permissions to [your authenticated principal](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration).

-> If you have assigned API permissions specific to the Azure Directory Graph API, you can safely unassign these permissions after upgrading to version 2.0.

### Assigning directory roles

If you are using directory roles to assign effective permissions to your authenticated principal, you may not necessarily need to assign new API permissions.

Whilst assigning directory roles is the recommended approach for user principals, if you are authenticating using a service principal, we recommend assigning permissions using app roles as detailed below. Note that this differs from our advice for earlier (v0.x and v1.x) releases, where we recommend the use of directory roles for service principals due to bugs previously encountered with Azure Active Directory Graph that do not occur with Microsoft Graph.

### Assigning new API permissions for a service principal

To assign permissions to your Application for use with Service Principal authentication, navigate to the [**Azure Active Directory** overview](https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview) within the [Azure Portal](https://portal.azure.com/) and select the [**App Registrations** blade](https://portal.azure.com#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps). Locate your registered application and click on its display name to manage it.

Go to the API Permissions pane for the Application and click the "Add a permission" button. In the pane that opens, select **Microsoft Graph**.

Choose "Application Permissions" for the permission type, and check the permissions you would like to assign. The permissions you need will depend on which directory objects you wish to manage with Terraform. The following table show the required permissions for some common resources:

Resource(s) | Role Name(s)
-------- | ---------------
`data.azuread_application`<br>`data.azuread_service_principal` | Application.Read.All
`data.azuread_domains` | Domain.Read.All
`data.azuread_group`<br>`data.azuread_groups` | Group.Read.All
`data.azuread_user`<br>`data.azuread_users` | User.Read.All
`azuread_application`<br>`azuread_application_certificate`<br>`azuread_application_password`<br>`azuread_service_principal`<br>`azuread_service_principal_certificate`<br>`azuread_service_principal_password` | Application.ReadWrite.All
`azuread_group`<br>`azuread_group_member` | Group.ReadWrite.All
`azuread_user` | User.ReadWrite.All

-> **Permissions for other resources** If the resource you are using is not shown in the above table, consult the documentation page for the resource for a guide to the required permissions.

Depending on the configuration of your AAD tenant, you may also need to grant the Directory.Read.All and/or Directory.ReadWrite.All roles.

After assigning permissions, you will need to grant consent for the service principal to utilise them. The easiest way to do this is by clicking the Grant Admin Consent button in the same API Permissions pane, which will create the necessary app role assignments for the Service Principal.

## New required fields

### Resource: `azuread_group`

The `mail_enabled` and `security_enabled` fields are no longer read-only, and at least one of these fields must be set to `true` in order to create a new group.

## Removal of deprecated fields

The following attributes/properties have been deprecated in the AzureAD provider, and has been removed in version 2.0.

~> **Compatibility Note** You will need to update your Terraform configuration in the latest v1.x release to use the new fields, prior to upgrading to 2.0.

### Provider

The deprecated field `metadata_host` is no longer used and has been removed.

### Data source: `azuread_application`

The deprecated field `name` has been replaced by the `display_name` field and has been removed.

The deprecated field `is_enabled` in the `app_roles` block has been replaced by the `enabled` field and has been removed.

The deprecated field `available_to_other_tenants` has been replaced by the `sign_in_audience` field and has been removed.

The deprecated field `homepage` has been replaced by the `homepage_url` field in the `web` block and has been removed.

The deprecated field `logout_url` has been replaced by the `logout_url` field in the `web` block and has been removed.

The deprecated field `oauth2_allow_implicit_flow` has been replaced by the `access_token_issuance_enabled` field in the `implicit_grant` block and has been removed.

The deprecated `oauth2_permissions` block has been replaced by the `oauth2_permission_scopes` block within the `api` block and has been removed.

The deprecated field `reply_urls` has been replaced by the `redirect_uris` field in the `web` block and has been removed.

The legacy `type` field is deprecated and has been removed.

### Data source: `azuread_group`

The deprecated field `name` has been replaced by the `display_name` field and has been removed.

### Data source: `azuread_groups`

The deprecated field `names` has been replaced by the `display_names` field and has been removed.

### Data source: `azuread_user`

The deprecated field `immutable_id` has been replaced by the `onpremises_immutable_id` field and has been removed.

The deprecated field `physical_delivery_office_name` has been replaced by the `office_location` field and has been removed.

The deprecated field `mobile` has been replaced by the `mobile_phone` field and has been removed.

### Data source: `azuread_users`

The deprecated field `immutable_id` in the `users` block has been replaced by the `onpremises_immutable_id` field and has been removed.

### Resource: `azuread_application`

The deprecated field `name` has been replaced by the `display_name` field and has been removed.

The deprecated field `is_enabled` in the `app_role` block has been replaced by the `enabled` field and has been removed.

The deprecated field `available_to_other_tenants` has been replaced by the `sign_in_audience` field and has been removed.

The deprecated field `homepage` has been replaced by the `homepage_url` field in the `web` block and has been removed.

The deprecated field `logout_url` has been replaced by the `logout_url` field in the `web` block and has been removed.

The deprecated field `oauth2_allow_implicit_flow` has been replaced by the `access_token_issuance_enabled` field in the `implicit_grant` block and has been removed.

The deprecated field `is_enabled` in the `app_role` block has been replaced by the `enabled` field and has been removed.

The deprecated `oauth2_permissions` block has been replaced by the `oauth2_permission_scope` block within the `api` block and has been removed.

-> In the new `oauth2_permission_scope` block, the `is_enabled` field has been renamed to `enabled` and the `id` field is now **required**. See the [New required UUID fields](#new-required-uuid-fields) section below for more information.

The deprecated field `public_client` has been replaced by the `fallback_public_client_enabled` field and has been removed.

The deprecated field `reply_urls` has been replaced by the `redirect_uris` field in the `web` block and has been removed.

The legacy `type` field is deprecated and has been removed.

### Resource: `azuread_application_app_role`

This resource has been removed in version 2.0 of the provider. Limitations of the `azuread_application` resource, that previously necessitated use of this resource, have been resolved in version 2.0. See [this pull request](https://github.com/hashicorp/terraform-provider-azuread/pull/465) for more information about this change.

### Resources: `azuread_application_oauth2_permission` and `azuread_application_oauth2_permission_scope`

In version 1.5.0 of the provider, the `azuread_application_oauth2_permission` resource was deprecated and replaced by the `azuread_application_oauth2_permission_scope` resource.

However, in version 2.0 of the provider, both of these resources have been removed. Limitations of the `azuread_application` resource, that previously necessitated use of these resources, have been resolved in version 2.0. See [this pull request](https://github.com/hashicorp/terraform-provider-azuread/pull/465) for more information about this change.

### Resource: `azuread_application_password`

The deprecated field `description` has been replaced by the `display_name` field and has been removed.

-> The following also applies when the Microsoft Graph beta is enabled in version 1.5 or later

The `key_id` field has become read-only as Azure Active Directory no longer allows user-specified key IDs for passwords. This also means that the `azuread_application_password` resource no longer supports importing in version 2.0 of the provider.

The `value` field has become read-only as Azure Active Directory no longer accepts user-supplied password values. Passwords are instead auto-generated by Azure and exported with the `value` attribute by the resource.

### Resource: `azuread_group`

The deprecated field `name` has been replaced by the `display_name` field and has been removed.

### Resource: `azuread_service_principal_password`

The deprecated field `description` has been replaced by the `display_name` field and has been removed.

-> The following also applies when the Microsoft Graph beta is enabled in version 1.5 or later

The `display_name`, `start_date` and `end_date` fields are no longer respected by the API and have been made read-only. Accordingly the `end_date_relative` field has been removed.

The `key_id` field has become read-only as Azure Active Directory no longer allows user-specified key IDs for passwords. This also means that the `azuread_service_principal_password` resource no longer supports importing in version 2.0 of the provider.

The `value` field has become read-only as Azure Active Directory no longer accepts user-supplied password values. Passwords are instead auto-generated by Azure and exported with the `value` attribute by the resource.

### Resource: `azuread_user`

The deprecated field `immutable_id` has been replaced by the `onpremises_immutable_id` field and has been removed.

The deprecated field `physical_delivery_office_name` has been replaced by the `office_location` field and has been removed.

The deprecated field `mobile` has been replaced by the `mobile_phone` field and has been removed.

## New required UUID fields

Several fields that were previously optional or read-only are now required in version 2.0. Currently, these fields are optionally (or always) autogenerated by the provider, as not all users need to be able to set these fields to specific values.

In v2.0, you can achieve the same behaviour by using the [Random provider](https://registry.terraform.io/providers/hashicorp/random/latest/docs), for example:

```hcl
resource "random_uuid" "example_role_id" {}

resource "azuread_application" "example" {
  display_name = "example"

  app_role {
    id = random_uuid.example_role_id.result
    # additional properties...
  }
}
```

Requiring these properties will enable more predictable and less disruptive management of app roles and OAuth 2.0 permission scopes such that existing roles/scopes will not be disabled or overwritten to accommodate changes in an application's configuration.

### Resource: `azuread_application`

The `id` field in the `app_role` block was previously currently Computed (read-only) but is now Required.

The `id` field in the deprecated `oauth2_permissions` block was previously Computed (read-only) but its replacement field `id` in the `oauth2_permission_scope` block is Required.

## Computed fields

In previous version of the provider, many fields were introduced as Optional + Computed fields. This meant that omitting such fields would cause Terraform to ignore them and not attempt to manage them. However, this approach has many side effects including the inability to unset or clear these fields, and sometimes being forced to accept an undesired default value.

To resolve these issues, many of these fields are no longer Computed in version 2.0 of the provider. This means that Terraform will manage these fields and if you do not specify their values in your configuration, they will be unset or set to their zero values. It's also appropriate in some cases for a field to be Computed, particularly where it helps prevent disruption to services or users.

Accordingly, in version 2.0 of the provider the following fields have changed.

### Resource: `azuread_application`

The `app_role` block is no longer Computed, omitting this block will cause Terraform to remove any app roles published by an application.

The `value` field in the `app_role` block is no longer Computed, omitting this field will cause Terraform to clear this field for an app role.

The `fallback_public_client_enabled` field is no longer Computed, omitting this field will cause Terraform to default this value to `false`.

The `identifier_uris` field is no longer Computed, omitting this field will cause Terraform to remove any identifier URIs configured for an application.

The `oauth2_permission_scope` block is no longer Computed, omitting this block will cause Terraform to remove any OAuth2 permission scopes published by an application.

The `owners` field is no longer Computed, omitting this field will cause Terraform remove any owners for an application. It's recommended to specify the object ID of the authenticated principal running Terraform, to ensure sufficient permissions that the application can be subsequently updated.

The `sign_in_audience` field is no longer Computed, omitting this field will cause Terraform to default this value to `AzureADMyOrg`.

The `web` block is no longer Computed, omitting this field will cause Terraform to remove all field values contained in this block, including the `homepage_url`, `logout_url`, `redirect_uris` and `access_token_issuance_enabled` fields.

## Resource: `azuread_user`

The `password` field is now Optional + Computed, which makes it possible to import existing users without forcibly resetting their password.

The `city` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `company_name` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `country` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `department` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `given_name` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `job_title` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `mobile_phone` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `office_location` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `postal_code` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `state` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `street_address` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `surname` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

The `usage_location` field is no longer Computed, omitting this field will cause Terraform to remove this value for a user.

### Avoiding diffs for already-published applications

In order to determine the assigned value of a given UUID-type property where you had previously not defined one, you can use the `terraform state show` command to inspect your existing application(s), and then add the current value to your configuration.

```shell
$ terraform state show azuread_application.example

# azuread_application.example:
resource "azuread_application" "example" {
    app_role           = [
        {
            allowed_member_types = [
                "User",
            ]
            description          = "Just an example"
            display_name         = "Example Role"
            id                   = "3dbd749f-c2ba-4796-0d33-273878f8a31b"
            is_enabled           = true
            value                = ""
        },
    ]
    application_id     = "cd986e6e-5a90-42f0-8c01-82f0da5b286c"
    display_name       = "example-app"
    id                 = "61a9d04a-6694-497a-afe8-6303aa3435df"
    oauth2_permissions = [
        {
            admin_consent_description  = "Allow the application to access example-app on behalf of the signed-in user."
            admin_consent_display_name = "Access example-app"
            id                         = "c32d857d-02d9-4887-af7a-cb1f1fd61b9a"
            is_enabled                 = true
            type                       = "User"
            user_consent_description   = "Allow the application to access example-app on your behalf."
            user_consent_display_name  = "Access example-app"
            value                      = "user_impersonation"
        },
    ]
    object_id          = "61a9d04a-6694-497a-afe8-6303aa3435df"
}
```

## Behaviour change: default `user_impersonation` scope for applications

With AzureAD v1.x using Azure Active Directory Graph (or using the beta support for Microsoft Graph), newly created applications are assigned a default `user_impersonation` scope which enables users to sign in to your application.

With AzureAD v2.0 and later using Microsoft Graph, this default scope is **not** automatically granted to your application, and you will need to specify it in your configuration. The following example can be used to replicate the earlier behaviour.

```hcl
resource "random_uuid" "example_app_user_impersonation_scope" {}

resource "azuread_application" "example" {
  display_name = "example-app"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Allow the application to access example-app on behalf of the signed-in user."
      admin_consent_display_name = "Access example-app"
      id                         = random_uuid.example_app_user_impersonation_scope.result
      is_enabled                 = true
      type                       = "User"
      user_consent_description   = "Allow the application to access example-app on your behalf."
      user_consent_display_name  = "Access example-app"
      value                      = "user_impersonation"
    }
  }
}
```

## Beta support for Microsoft Graph in v1.5.0

In version 1.5.0 or later of the AzureAD provider, beta support for Microsoft Graph can be enabled in the provider block.

~> Please note that whilst we do not recommend using this feature in production during the beta phase, you can use this feature flag to determine any configuration changes that you may need to make in advance of the v2.0 release.

With this feature enabled, most requests will be sent to Microsoft Graph, however for compatibility some requests may still be sent to Azure Active Directory Graph. You should therefore ensure your authenticated principal is assigned the required permissions for both APIs.

To enable this support, set the following in your provider block:

```hcl
provider "azuread" {
  use_microsoft_graph = true
}
```

Alternatively, you can enable the beta by setting the following environment variable to any non-empty value:

```shell
# sh
export AAD_USE_MICROSOFT_GRAPH=1

# PowerShell
$env:AAD_USE_MICROSOFT_GRAPH = 1
```

We appreciate any feedback you might have whilst using beta support for Microsoft Graph. Bug reports and feature requests can be logged on our [GitHub issue tracker](https://github.com/hashicorp/terraform-provider-azuread/issues), and you can also check out the [project README](https://github.com/hashicorp/terraform-provider-azuread/blob/main/README.md) for information on how to get in touch with the maintainers. If you encounter an issue that has already been reported, please upvote it and comment to add any additional context you might have.
