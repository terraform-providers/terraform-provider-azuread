## 2.0.0 (Unreleased)

NOTES:

* **Major Version:** This is a major version upgrade which contains breaking changes. Please read the [Upgrade Guide](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/microsoft-graph) before upgrading, which details all the known breaking changes that practitioners should be aware of.
* **Microsoft Graph:** The upstream API for Azure Active Directory is now Microsoft Graph, and the deprecated Azure Active Directory Graph API is no longer supported.

FEATURES:

* **New Resource** `application_pre_authorized` [GH-472]

IMPROVEMENTS:

* `data.azuread_application` - the `implicit_grant` block now supports the `id_token_issuance_enabled` attribute [GH-461]
* `data.azuread_application` - the `optional_claims` block now supports the `saml2_token` attribute [GH-461]
* `data.azuread_domains` - export the `admin_managed`, `root` and `supported_services` attributes for each domain [GH-461]
* `data.azuread_domains` - support the `admin_managed`, `only_root` and `supports_services` properties [GH-461]
* `azuread_application` - the `implicit_grant` block now supports the `id_token_issuance_enabled` property [GH-461]
* `azuread_application` - the `optional_claims` block now supports the `saml2_token` block [GH-461]
* `azuread_application` - the `sign_in_audience` property now supports the `AzureADandPersonalMicrosoftAccount` and `PersonalMicrosoftAccount` values [GH-461]
* `azuread_group` - support for creating mail-enabled groups [GH-461]
* `azuread_group` - support for creating Microsoft 365 groups [GH-461]
* `azuread_group` - support for updating groups without recreating them [GH-461]

BUG FIXES:

* `azuread_application` - resolved an issue where `identifier_uris` could be reordered and cause a persistent diff [GH-461]
* `azuread_application` - the `identifier_uris` property can now be set for all applications regardless of target platform [GH-461]
* `azuread_application` - fixed a bug where app roles could be duplicated or left in a disabled state [GH-461]
* `azuread_application` - fixed a bug where app roles could not be removed from an application [GH-461]
* `azuread_application` - fixed a bug where the `enabled` property of app roles could be ignored [GH-461]
* `azuread_application` - fixed a bug where the `id` property of app roles could be undesirably changed [GH-461]
* `azuread_application` - resolved an issue where the default scope could not be removed from an application [GH-461]
* `azuread_application` - resolved an issue where multiple `group_membership_claims` could not be specified [GH-461]
* `azuread_application_password` - the `display_name` / `description` properties are no longer stored using the `customKeyIdentifier` API field, lifting the 32 byte limit [GH-461]
* `azuread_user` - resolved an issue where importing users would inadvertently reset their password [GH-461]

BREAKING CHANGES:

* `data.azuread_domains` - the `is_` prefix has been dropped from all exported attributes [GH-461]
* `data.azuread_application` - the deprecated property `name` has been removed [GH-461]
* `data.azuread_application` - the deprecated attribute `available_to_other_tenants` has been removed [GH-461]
* `data.azuread_application` - the `group_membership_claims` attribute has changed from a string to a list of strings [GH-461]
* `data.azuread_application` - the deprecated attribute `homepage` has been removed [GH-461]
* `data.azuread_application` - the deprecated attribute `logout_url` has been removed [GH-461]
* `data.azuread_application` - the deprecated attribute `oauth2_allow_implicit_flow` has been removed [GH-461]
* `data.azuread_application` - the deprecated attribute `oauth2_permissions` has been removed [GH-461]
* `data.azuread_application` - the `public_client` attribute is now a block containing public client settings [GH-461]
* `data.azuread_application` - the deprecated attribute `reply_urls` has been removed [GH-461]
* `data.azuread_application` - the deprecated attribute `type` has been removed [GH-461]
* `data.azuread_group` - the deprecated property `name` has been removed [GH-461]
* `data.azuread_groups` - the deprecated property `names` has been removed [GH-461]
* `data.azuread_service_principal` - the deprecated attribute `oauth2_permissions` has been removed [GH-461]
* `data.azuread_user` - the deprecated attribute `immutable_id` has been removed [GH-461]
* `data.azuread_user` - the deprecated attribute `physical_delivery_office_name` has been removed [GH-461]
* `data.azuread_user` - the deprecated attribute `mobile` has been removed [GH-461]
* `data.azuread_users` - the deprecated attribute `immutable_id` in the `users` block has been removed [GH-461]
* `azuread_application` - the deprecated property `name` has been removed [GH-461]
* `azuread_application` - the `api` block is no longer Computed, omitting this block will cause it to be removed from your configuration [GH-461]
* `azuread_application` - the `app_role` block is no longer Computed, omitting this block will cause it to be removed from your configuration [GH-461]
* `azuread_application` - the `id` property in the `app_role` block is now Required [GH-461]
* `azuread_application` - the deprecated property `available_to_other_tenants` has been removed [GH-461]
* `azuread_application` - the `fallback_public_client_enabled` property is no longer Computed, omitting this property will cause the default value to be applied [GH-461]
* `azuread_application` - the `group_membership_claims` property has changed from a string to a set of strings [GH-461]
* `azuread_application` - the deprecated property `homepage` has been removed [GH-461]
* `azuread_application` - the `identifier_uris` property is no longer Computed, omitting this property will cause it to be removed from your configuration [GH-461]
* `azuread_application` - the deprecated property `logout_url` has been removed [GH-461]
* `azuread_application` - the deprecated property `oauth2_allow_implicit_flow` has been removed [GH-461]
* `azuread_application` - the `oauth2_permission_scope` block is no longer Computed, omitting this block will cause it to be removed from your configuration [GH-461]
* `azuread_application` - the deprecated block `oauth2_permissions` has been removed [GH-461]
* `azuread_application` - the `owners` property is no longer Computed, omitting this property will cause it to be removed from your configuration [GH-461]
* `azuread_application` - the `public_client` property is now a block containing public client settings [GH-461]
* `azuread_application` - the deprecated property `reply_urls` has been removed [GH-461]
* `azuread_application` - the `sign_in_audience` property is no longer Computed, omitting this property will cause the default value to be applied [GH-461]
* `azuread_application` - the deprecated property `type` has been removed [GH-461]
* `azuread_application` - the `web` block is no longer Computed, omitting this block will cause it to be removed from your configuration [GH-461]
* `azuread_application_password` - the `key_id` and `value` properties are now Computed, due to API changes it is no longer possible to specify these values [GH-461]
* `azuread_group` - the deprecated property `name` has been removed [GH-461]
* `azuread_group` - at least one of the `mail_enabled` or `security_enabled` properties are now Required [GH-461]
* `azuread_service_principal` - the deprecated attribute `oauth2_permissions` has been removed [GH-461]
* `azuread_service_principal_password` - the `key_id` and `value` properties are now Computed, due to API changes it is no longer possible to specify these values [GH-461]
* `azuread_service_principal_password` - the `start_date` and `end_date` properties are now Computed, due to an API issue it is no longer possible to specify these values [GH-461]
* `azuread_user` - the deprecated property `immutable_id` has been removed [GH-461]
* `azuread_user` - the deprecated property `physical_delivery_office_name` has been removed [GH-461]
* `azuread_user` - the deprecated property `mobile` has been removed [GH-461]

## 1.6.0 (June 24, 2021)

DEPRECATIONS:

* `azuread_application_app_role` - this resource is deprecated and will be removed in version 2.0 ([#465](https://github.com/terraform-providers/terraform-provider-azuread/issues/465))
* `azuread_application_oauth2_permission` - this resource is deprecated and will be removed in version 2.0 ([#465](https://github.com/terraform-providers/terraform-provider-azuread/issues/465))
* `azuread_application_oauth2_permission_scope` - this resource is deprecated and will be removed in version 2.0 ([#465](https://github.com/terraform-providers/terraform-provider-azuread/issues/465))

## 1.5.1 (June 10, 2021)

BUG FIXES:

* **Provider:** Suppress a spurious deprecation notice for the `metadata_host` provider field ([#439](https://github.com/terraform-providers/terraform-provider-azuread/issues/439))
* `azuread_application_password` - fix a bug that prevented specifying the `display_name`, `start_date`, `end_date` or `end_date_relative` properties when using Microsoft Graph ([#444](https://github.com/terraform-providers/terraform-provider-azuread/issues/444))
* `azuread_group` - fix a bug that prevented creating a group with more than 20 owners or members ([#454](https://github.com/terraform-providers/terraform-provider-azuread/issues/454))
* `azuread_service_principal_password` - fix a bug that prevented specifying the `display_name`, `start_date`, `end_date` or `end_date_relative` properties when using Microsoft Graph ([#444](https://github.com/terraform-providers/terraform-provider-azuread/issues/444))

## 1.5.0 (May 20, 2021)

NOTES:

* **Support for Microsoft Graph:** This release introduces beta support for [Microsoft Graph](https://docs.microsoft.com/en-us/graph/overview) in a way that is forward (and backward) compatible with the current [Azure Active Directory Graph](https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-graph-api) API implementation. We do not recommend enabling this beta _in production_ at this time, but encourage you to try it out in test environments where minimal impact can occur if something doesn't work as expected. See the [Migration Guide](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/microsoft-graph#beta-support-for-microsoft-graph-in-v150) for more details.

* **Deprecations:** This release contains a number of additional deprecations to aid in future upgrades to version 2.0 of this provider. These will be flagged when running Terraform, and are documented in detail in the [Migration Guide](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/microsoft-graph). Existing configurations will continue to work unchanged for any v1.x release, regardless of which API is used.

IMPROVEMENTS:

* `data.azuread_user` - export the `user_type` attribute ([#406](https://github.com/terraform-providers/terraform-provider-azuread/issues/406))
* `azuread_user` - export the `user_type` attribute ([#401](https://github.com/terraform-providers/terraform-provider-azuread/issues/401)] / [[#413](https://github.com/terraform-providers/terraform-provider-azuread/issues/413))

BUG FIXES:

* `azuread_application` - validation for the `identifier_uris` property now supports URNs ([#426](https://github.com/terraform-providers/terraform-provider-azuread/issues/426))

## 1.4.0 (February 18, 2021)

IMPROVEMENTS:

* dependencies: updating to build using Go 1.16 which adds support for `darwin/arm64` (Apple Silicon) ([#403](https://github.com/hashicorp/terraform-provider-azuread/issues/403))
* Data Source: `azuread_group` - support for the `mail_enabled` and `security_enabled` properties ([#393](https://github.com/hashicorp/terraform-provider-azuread/issues/393))
* `azuread_group` - support for the `mail_enabled` and `security_enabled` attributes ([#393](https://github.com/hashicorp/terraform-provider-azuread/issues/393))

## 1.3.0 (January 28, 2021)

IMPROVEMENTS:

* `azuread_application_certificate` - support for base64 and hex encoded certificate values ([#386](https://github.com/hashicorp/terraform-provider-azuread/issues/386))
* `azuread_service_principal_certificate` - support for base64 and hex encoded certificate values ([#386](https://github.com/hashicorp/terraform-provider-azuread/issues/386))

## 1.2.2 (January 16, 2021)

BUGFIXES:

* `azuread_application` - set the display name correctly when creating/updating applications using the `display_name` property

## 1.2.1 (January 14, 2021)

BUGFIXES:

* `data.azuread_application` - correctly set the `display_name` attribute in state.
* `azuread_application` - correctly set the `display_name` attribute in state.

## 1.2.0 (January 14, 2021)

NOTES:

* **Terraform Plugin SDK Upgrade:** This version upgrades the Terraform Plugin SDK to v2.3.0. This does not provide any additional provider features or resources but is useful for developers and part of our development roadmap.
* **Refactor into multiple packages:** As part of our preparation for Microsoft Graph support, this release refactors resources and data sources into separate Go packages.

IMPROVEMENTS:

* `azuread_application` - support new values `include_externally_authenticated_upn`, `include_externally_authenticated_upn_without_hash`, and `use_guid` for the `additional_properties` property of the `optional_claims` block.

DEPRECATIONS:

* `data.azuread_application` - the `name` property has been renamed to `display_name` and will be removed in version 2.0.
* `data.azuread_group` - the `name` property has been renamed to `display_name` and will be removed in version 2.0.
* `data.azuread_groups` - the `names` property has been renamed to `display_names` and will be removed in version 2.0.
* `azuread_application` - the `name` property has been renamed to `display_name` and will be removed in version 2.0.
* `azuread_application` - the `type` property is now deprecated and will be removed in version 2.0, as there is no longer any distinction between native and webapp/api applications.
* `azuread_group` - the `name` property has been renamed to `display_name` and will be removed in version 2.0.

## 1.1.1 (November 26, 2020)

BUG FIXES:

* `azuread_application` - resolves an issue where setting `prevent_duplicate_names = true` causes an error for new applications ([#367](https://github.com/hashicorp/terraform-provider-azuread/issues/367))
* `azuread_application` - fixes a bug where the default owner for a new application is removed ([#366](https://github.com/hashicorp/terraform-provider-azuread/issues/366))

## 1.1.0 (November 25, 2020)

FEATURES:

* Added a flag to allow users to customize the Partner ID or opt-out of the default Terraform Partner ID ([#350](https://github.com/hashicorp/terraform-provider-azuread/issues/350))
* This release includes updated support for working directly with tenants using Azure CLI authentication. We recommend the use of `az login --allow-no-subscription` to populate tenant-level accounts (which have no subscriptions).

IMPROVEMENTS:

* `data.azuread_user` - support the `given_name`, `surname`, `job_title`, `department`, `company_name`, `physical_delivery_office_name`, `street_address`, `city`, `state`, `country`, `postal_code` and `mobile` attribute ([#351](https://github.com/hashicorp/terraform-provider-azuread/issues/351))
* `azuread_user` - support the `given_name`, `surname`, `job_title`, `department`, `company_name`, `physical_delivery_office_name`, `street_address`, `city`, `state`, `country`, `postal_code` and `mobile` properties ([#351](https://github.com/hashicorp/terraform-provider-azuread/issues/351))

BUG FIXES:

* **Provider:** Fixed an issue where CLI authentication produced a `parsing json result` error during provider initialization ([#358](https://github.com/hashicorp/terraform-provider-azuread/issues/358))
* `azuread_application` - enable removal of owners on existing applications, and creation of applications with no owners ([#355](https://github.com/hashicorp/terraform-provider-azuread/issues/355))
* `azuread_application` - fixed a bug where specifying the `prevent_duplicate_names` property would report a false positive on update. ([#338](https://github.com/hashicorp/terraform-provider-azuread/issues/338))

## 1.0.0 (September 03, 2020)

NOTES:

* **Major Version:** This is a major version upgrade which contains some breaking changes as detailed below.
* **Terraform 0.10/0.11:** This version of the provider requires Terraform 0.12.x or later and will not work with earlier versions.

FEATURES:

* New resource: `azuread_application_app_role` ([#150](https://github.com/hashicorp/terraform-provider-azuread/issues/150)] [[#306](https://github.com/hashicorp/terraform-provider-azuread/issues/306))
* New resource: `azuread_application_oauth2_permission` ([#267](https://github.com/hashicorp/terraform-provider-azuread/issues/267))

BREAKING CHANGES:

* `azuread_application` - a default value for the `homepage` property is no longer derived when unspecified ([#268](https://github.com/hashicorp/terraform-provider-azuread/issues/268))
* `azuread_application_password` - the deprecated `application_id` property has been removed
* `data.azuread_group` - the `name` property is now case-insensitive ([#246](https://github.com/hashicorp/terraform-provider-azuread/issues/246))
* `data.azuread_groups` and `data.azuread_users` will not error if no results found

## 0.11.0 (July 09, 2020)

IMPROVEMENTS:

* Provider: no longer require configuring `subscription_id` (configuration value) / `ARM_SUBSCRIPTION_ID` (environment variable). ([#271](https://github.com/hashicorp/terraform-provider-azuread/issues/271))
* `data.azuread_client_config` - deprecate the `subscription_id` property. For compatibility, still populates `subscription_id` if the provider is configured with a subscription ID ([#271](https://github.com/hashicorp/terraform-provider-azuread/issues/271))
* `data.azuread_application` - support for the `application_id` property ([#274](https://github.com/hashicorp/terraform-provider-azuread/issues/274))
* `data.azuread_users` - support the `ignore_missing` property ([#256](https://github.com/hashicorp/terraform-provider-azuread/issues/256))
* `data.azuread_users` - export the `users` attribute containing a list of users with additional properties ([#256](https://github.com/hashicorp/terraform-provider-azuread/issues/256))
* `azuread_application` - support the `prevent_duplicate_names` property ([#279](https://github.com/hashicorp/terraform-provider-azuread/issues/279))
* `azuread_application` - validate `app_roles` and `oauth2_permissions` to check for duplicate `value`s ([#287](https://github.com/hashicorp/terraform-provider-azuread/issues/287))
* `azuread_group` - support the `prevent_duplicate_names` property ([#279](https://github.com/hashicorp/terraform-provider-azuread/issues/279))

BUG FIXES:

* `azuread_group` - remediate AAD replication delays when adding/removing group members ([#283](https://github.com/hashicorp/terraform-provider-azuread/issues/283))
* `azuread_group` - remediate AAD replication delays after group creation, before setting owners/members ([#290](https://github.com/hashicorp/terraform-provider-azuread/issues/290))

## 0.10.0 (June 05, 2020)

BREAKING CHANGES:

* `azuread_application` - the `oauth2_permissions` attribute has changed from a list to a set. If you are referencing this attribute with explicit list indexes, you will need to update your configuration to use a `for` expression. For example:

    ```hcl
    id = azuread_application.example.oauth2_permissions[0].id
    ```

    becomes

    ```hcl
    id = [for permission in azuread_application.example.oauth2_permissions : permission.id][0]
    ```

FEATURES:

* **New Resource:** `azuread_application_certificate` ([#262](https://github.com/hashicorp/terraform-provider-azuread/issues/262))
* **New Resource:** `azuread_service_principal_certificate` ([#262](https://github.com/hashicorp/terraform-provider-azuread/issues/262))

IMPROVEMENTS:

* `azuread_application` - support for the `optional_claims` property, for access tokens and ID tokens ([#260](https://github.com/hashicorp/terraform-provider-azuread/issues/260))
* `azuread_application` - support for the `oauth2_permissions` property ([#252](https://github.com/hashicorp/terraform-provider-azuread/issues/252))
* `azuread_application_password` - support the `description` property ([#253](https://github.com/hashicorp/terraform-provider-azuread/issues/253))
* `azuread_service_principal_password` - support the `description` property ([#253](https://github.com/hashicorp/terraform-provider-azuread/issues/253))
* `data.azuread_users` - support empty lists for `user_principal_names`/`object_ids`/`mail_nicknames` properties ([#258](https://github.com/hashicorp/terraform-provider-azuread/issues/258))
* `data.azuread_groups` - support empty lists for `names`/`object_ids` properties ([#257](https://github.com/hashicorp/terraform-provider-azuread/issues/257))

BUG FIXES:

* `azuread_application_password` and `azuread_service_principal_password` - Plan-time validation for `end_date` / `end_date_relative` ([#261](https://github.com/hashicorp/terraform-provider-azuread/issues/261))
* `azuread_application_password` and `azuread_service_principal_password` - Change the resource ID format to mitigate potential UUID collision ([#264](https://github.com/hashicorp/terraform-provider-azuread/issues/264))

## 0.9.0 (May 15, 2020)

DEPENDENCIES:

* upgrade `azure-sdk-for-go` to `v42.1.0` ([#247](https://github.com/hashicorp/terraform-provider-azuread/issues/247))

IMPROVEMENTS:

* `azuread_application` - the `group_membership_claims` property now supports `ApplicationGroup` ([#238](https://github.com/hashicorp/terraform-provider-azuread/issues/238))
* `azuread_service_principal` - changing the `tags` property no longer forces a new resource ([#245](https://github.com/hashicorp/terraform-provider-azuread/issues/245))

BUG FIXES:

* `data.azuread_user` - use `equals` instead of `startsWith` when looking uo users by `mailNickname` ([#251](https://github.com/hashicorp/terraform-provider-azuread/issues/251))
* `data.azuread_users` - use `equals` instead of `startsWith` when looking uo users by `mailNickname` ([#251](https://github.com/hashicorp/terraform-provider-azuread/issues/251))

## 0.8.0 (March 16, 2020)

FEATURES:

* **New Data Source:** `azuread_client_config` ([#229](https://github.com/hashicorp/terraform-provider-azuread/issues/229))

IMPROVEMENTS:

* dependencies: upgrade `azure-sdk-for-go` to `v40.3.0` ([#225](https://github.com/hashicorp/terraform-provider-azuread/issues/225))
* dependencies: upgrade `go-autorest/autorest` to `v0.10.0` ([#225](https://github.com/hashicorp/terraform-provider-azuread/issues/225))
* dependencies: upgrade `terraform-plugin-sdk` to `v1.6.0` ([#225](https://github.com/hashicorp/terraform-provider-azuread/issues/225))
* `azuread_application` - support for the `logout_url` property ([#226](https://github.com/hashicorp/terraform-provider-azuread/issues/226))
* `azuread_group` - support for the `description` property ([#216](https://github.com/hashicorp/terraform-provider-azuread/issues/216))
* `azuread_user` - support for the `onpremises_sam_account_name` and `onpremises_user_principal_name` properties ([#222](https://github.com/hashicorp/terraform-provider-azuread/issues/222))
* `azuread_user` - support for the `immutable_id` property ([#207](https://github.com/hashicorp/terraform-provider-azuread/issues/207))

BUG FIXES:

* `azuread_application` - ensure all owners are added before removed ([#226](https://github.com/hashicorp/terraform-provider-azuread/issues/226))
* `azuread_application_password` - validate the `length` property is less then `863` ([#228](https://github.com/hashicorp/terraform-provider-azuread/issues/228))
* `azuread_group` - the `owners` property is now additive during creation allowing an existing owner to be provided ([#211](https://github.com/hashicorp/terraform-provider-azuread/issues/211))
* `azuread_group_member` - mark as missing when member cannot be found instead of erroring ([#227](https://github.com/hashicorp/terraform-provider-azuread/issues/227))
* `azuread_service_principal_password` - validate the `length` property is less then `863` ([#228](https://github.com/hashicorp/terraform-provider-azuread/issues/228))

## 0.7.0 (November 15, 2019)

IMPROVEMENTS:

* provider: migrate to standalone plugin SDK v1.1.0 ([#154](https://github.com/hashicorp/terraform-provider-azuread/issues/154))
* provider: using the current (rather than the vendored) version of Terraform Core in user agents ([#154](https://github.com/hashicorp/terraform-provider-azuread/issues/154))
* `azuread_application` - adds ability to build homepage with HTTP in addition to HTTPS ([#155](https://github.com/hashicorp/terraform-provider-azuread/issues/155))
* `azuread_application` - allow the `app_role` block `value` property to be nil ([#157](https://github.com/hashicorp/terraform-provider-azuread/issues/157))
* `azuread_user` - support for the `usage_location` property ([#141](https://github.com/hashicorp/terraform-provider-azuread/issues/141))
* `data.azuread_user` - support looking up a user with `mail_nickname` ([#161](https://github.com/hashicorp/terraform-provider-azuread/issues/161))
* `data.azuread_users` - support looking up users with `mail_nicknames` ([#161](https://github.com/hashicorp/terraform-provider-azuread/issues/161))

## 0.6.0 (August 21, 2019)

IMPROVEMENTS:

* dependencies: upgrading `github.com/Azure/azure-sdk-for-go` to `v32.5.0` ([#140](https://github.com/hashicorp/terraform-provider-azuread/issues/140))
* dependencies: upgrading `github.com/Azure/go-autorest` to `v13.0.0` ([#140](https://github.com/hashicorp/terraform-provider-azuread/issues/140))
* dependencies: upgrading `github.com/hashicorp/go-azure-helpers` to `v0.7.0` ([#140](https://github.com/hashicorp/terraform-provider-azuread/issues/140))
* dependencies: upgrading `github.com/hashicorp/terraform` to `0.12.6` ([#133](https://github.com/hashicorp/terraform-provider-azuread/issues/133))
* `azuread_service_principal` - support for the `app_role_assignment_required` property ([#127](https://github.com/hashicorp/terraform-provider-azuread/issues/127))


## 0.5.1 (July 24, 2019)

BUG FIXES:

* `azuread_application_password` - fix incorrect conflicts with ([#129](https://github.com/hashicorp/terraform-provider-azuread/issues/129))

## 0.5.0 (July 24, 2019)

FEATURES:

* **New Data Source:** `azuread_users` ([#109](https://github.com/hashicorp/terraform-provider-azuread/issues/109))
* **New Resource:** `azuread_group_member` ([#100](https://github.com/hashicorp/terraform-provider-azuread/issues/100))

IMPROVEMENTS:

* `azuread_application` - support for the `app_roles` property ([#98](https://github.com/hashicorp/terraform-provider-azuread/issues/98))
* `azuread_application` - the `identifier_uris` property now allows `api`,`urn`, and `ms-appx` URI schemas ([#115](https://github.com/hashicorp/terraform-provider-azuread/issues/115))
* `azuread_application_password` - deprecation of `application_id` in favour of `application_object_id` ([#107](https://github.com/hashicorp/terraform-provider-azuread/issues/107))
* `azuread_group` - support for the `members` property ([#100](https://github.com/hashicorp/terraform-provider-azuread/issues/100))
* `azuread_group` - support for the `owners` property ([#62](https://github.com/hashicorp/terraform-provider-azuread/issues/62))
* `azuread_service_principal` - export the `oauth2_permissions` property ([#103](https://github.com/hashicorp/terraform-provider-azuread/issues/103))
* `data.azuread_application` - support for the `app_roles` property ([#110](https://github.com/hashicorp/terraform-provider-azuread/issues/110))
* `data.azuread_service_principal` - export the `app_roles` property ([#110](https://github.com/hashicorp/terraform-provider-azuread/issues/110))

BUG FIXES:

* `azuread_application_password` - will now wait for replication on resource creation ([#118](https://github.com/hashicorp/terraform-provider-azuread/issues/118))
* `azuread_service_principal_password` - will now wait for replication on resource creation ([#117](https://github.com/hashicorp/terraform-provider-azuread/issues/117))

## 0.4.0 (June 06, 2019)

NOTES:

* Resource creation potentially could take longer after this release as the provider will now attempt to wait for replication like the az cli tool. 

FEATURES:

* **New Resource:** `azuread_application_password` ([#71](https://github.com/hashicorp/terraform-provider-azuread/issues/71))

IMPROVEMENTS:

* dependencies: upgrading to `v0.12.0` of `github.com/hashicorp/terraform` ([#82](https://github.com/hashicorp/terraform-provider-azuread/issues/82))
* `azuread_application` - support for the `group_membership_claims` property ([#78](https://github.com/hashicorp/terraform-provider-azuread/issues/78))
* `azuread_application` - now exports the `oauth2_permissions` property ([#79](https://github.com/hashicorp/terraform-provider-azuread/issues/79))
* `azuread_application` - now exports the `object_id` property ([#99](https://github.com/hashicorp/terraform-provider-azuread/issues/99))
* `azuread_application` - support for the `type` property enabling the creation of `native` applications ([#74](https://github.com/hashicorp/terraform-provider-azuread/issues/74))
* `azuread_application` - will now wait for replication by waiting for 10 successful reads after creation ([#93](https://github.com/hashicorp/terraform-provider-azuread/issues/93))
* `azuread_group` - will now wait for replication by waiting for 10 successful reads after creation ([#91](https://github.com/hashicorp/terraform-provider-azuread/issues/91))
* `azuread_group` - now exports the `object_id` property ([#99](https://github.com/hashicorp/terraform-provider-azuread/issues/99))
* `azuread_service_principal` - will now wait for replication by waiting for 10 successful reads after creation ([#93](https://github.com/hashicorp/terraform-provider-azuread/issues/93))
* `azuread_service_principal` - now exports the `object_id` property ([#99](https://github.com/hashicorp/terraform-provider-azuread/issues/99))
* `azuread_user` - will now wait for replication by waiting for 10 successful reads after creation ([#91](https://github.com/hashicorp/terraform-provider-azuread/issues/91))
* `azuread_user` - increase the maximum allowed length of `password` to 256 ([#81](https://github.com/hashicorp/terraform-provider-azuread/issues/81))
* `azuread_user` - now exports the `object_id` property ([#99](https://github.com/hashicorp/terraform-provider-azuread/issues/99))
* `data.azuread_application` - now exports the `group_membership_claims` property ([#78](https://github.com/hashicorp/terraform-provider-azuread/issues/78))
* `data.azuread_application` - now exports the `oauth2_permissions` property ([#79](https://github.com/hashicorp/terraform-provider-azuread/issues/79))

## 0.3.1 (April 18, 2019)

BUG FIXES:

* Release fixing metadata to register the provider as compatible with Terraform 0.12.

## 0.3.0 (April 18, 2019)

NOTES:

* This release includes a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#56](https://github.com/hashicorp/terraform-provider-azuread/issues/56))

BUG FIXES:

* `azuread_application` - the order of the `reply_urls` property no longer matters ([#61](https://github.com/hashicorp/terraform-provider-azuread/issues/61))

## 0.2.0 (March 12, 2019)

FEATURES:

* **New Data Source:** `azuread_domains` ([#27](https://github.com/hashicorp/terraform-provider-azuread/issues/27))
* **New Data Source:** `azuread_group` ([#14](https://github.com/hashicorp/terraform-provider-azuread/issues/14))
* **New Resource:** `azuread_group` ([#14](https://github.com/hashicorp/terraform-provider-azuread/issues/14))

IMPROVEMENTS:

* dependencies: switching to use Go Modules ([#26](https://github.com/hashicorp/terraform-provider-azuread/issues/26))
* dependencies: updating `github.com/Azure/azure-sdk-for-go` to v24.1.0 ([#25](https://github.com/hashicorp/terraform-provider-azuread/issues/25))
* dependencies: updating `github.com/Azure/go-autorest` to v11.2.8 ([#24](https://github.com/hashicorp/terraform-provider-azuread/issues/24))
* validation: adding validation to all fields ([#30](https://github.com/hashicorp/terraform-provider-azuread/issues/30))
* `azuread_application` - support for `required_resource_access` property ([#23](https://github.com/hashicorp/terraform-provider-azuread/issues/23))
* `azuread_service_principal` - support for the `tags` property ([#31](https://github.com/hashicorp/terraform-provider-azuread/issues/31))
* `azuread_service_principal_password` - support for realitive ends dates with the `end_date_relative` property ([#53](https://github.com/hashicorp/terraform-provider-azuread/issues/53))

BUG FIXES:

* `azuread_application` - correctly reading back the `reply_urls` property into state ([#21](https://github.com/hashicorp/terraform-provider-azuread/issues/21))


## 0.1.0 (January 09, 2019)

Initial release of the Azure Active Directory provider - featuring resources split out from the AzureRM Provider.

FEATURES:

* New Data Source: `azuread_application`
* New Data Source: `azuread_service_principal`
* New Resource: `azuread_application`
* New Resource: `azuread_service_principal`
* New Resource: `azuread_service_principal_password`
