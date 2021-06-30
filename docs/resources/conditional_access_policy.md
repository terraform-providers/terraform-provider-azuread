---
subcategory: "Conditional Access"
---

# Resource: azuread_conditional_access_policy

Manages a Conditional Access Policy within Azure Active Directory.

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to `Policy.Read.All` and `Policy.ReadWrite.ConditionalAccess` within the `Windows Azure Active Directory` API.

## Example Usage

```terraform
resource "azuread_conditional_access_policy" "example" {
  display_name = "example policy"
  state        = "disabled"

  conditions {
    applications {
      included_applications = ["All"]
      excluded_applications = ["00000004-0000-0ff1-ce00-000000000000"]
    }
    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
    client_app_types = ["all"]
    locations {
      included_locations = ["All"]
      excluded_locations = ["AllTrusted"]
    }
    platforms {
      included_platforms = ["android"]
      excluded_platforms = ["iOS"]
    }
    sign_in_risk_levels = ["medium"]
    user_risk_levels    = ["medium"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }

  session_controls {
    application_enforced_restrictions {
      enabled = true
    }
    cloud_app_security {
      enabled                 = true
      cloud_app_security_type = "monitorOnly"
    }
    sign_in_frequency {
      enabled = true
      type    = "hours"
      value   = 10
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The friendly name for this Conditional Access Policy.
* `state` - (Required) Specifies the state of the policy object. Possible values are: `enabled`, `disabled` and `enabledForReportingButNotEnforced`

* `conditions` - (Required) A `conditions` block as documented below, which specifies the rules that must be met for the policy to apply.
* `grant_controls` - (Required) A `grant_controls` block as documented below, which specifies the grant controls that must be fulfilled to pass the policy.
* `session_controls` - (Optional) A `session_controls` block as documented below, which specifies the session controls that are enforced after sign-in.

---

`conditions` block supports the following:

* `applications` - (Optional) An `applications` block as documented below, which specifies applications and user actions included in and excluded from the policy.
* `users` - (Optional) A `users` block as documented below, which specifies users, groups, and roles included in and excluded from the policy.
* `client_app_types` - (Optional) A list of client application types included in the policy. Possible values are: `all`, `browser`, `mobileAppsAndDesktopClients`, `exchangeActiveSync`, `easSupported` and `other`.
* `locations` - (Optional) a `locations` block as documented below, which specifies locations included in and excluded from the policy.
* `platforms` - (Optional) a `platforms` block as documented below, which specifies platforms included in and excluded from the policy.
* `sign_in_risk_levels` - (Optional) A list of sign-in risk levels included in the policy. Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`.
* `user_risk_levels` - (Optional) A list of user risk levels included in the policy. Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`.

`applications` block supports the following:

* `included_applications` - (Optional) A list of application IDs the policy applies to, unless explicitly excluded (in `excluded_applications`). Can also be set to `All`.
* `excluded_applications` - (Optional) A list of application IDs explicitly excluded from the policy.
* `included_user_actions` - (Optional) A list of user actions to include. Supported values are `urn:user:registersecurityinfo` and `urn:user:registerdevice`.

`users` block supports the following:

* `included_users` - (Optional) A list of user IDs in scope of policy unless explicitly excluded, or `None` or `All` or `GuestsOrExternalUsers`.
* `excluded_users` - (Optional) A list of user IDs excluded from scope of policy and/or `GuestsOrExternalUsers`.
* `included_groups` - (Optional) A list of group IDs in scope of policy unless explicitly excluded, or `All`.
* `excluded_groups` - (Optional) A list of group IDs excluded from scope of policy.
* `included_roles` - (Optional) A list of role IDs in scope of policy unless explicitly excluded, or `All`.
* `excluded_roles` - (Optional) A list of role IDs excluded from scope of policy.

`locations` block supports the following:

* `included_locations` - (Optional) A list of location IDs in scope of policy unless explicitly excluded. Can also be set to `All`, or `AllTrusted`.
* `excluded_locations` - (Optional) A list of location IDs excluded from scope of policy.

`platforms` block supports the following:

* `included_platforms` - (Optional) A list of platforms the policy applies to, unless explicitly excluded. Possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`.
* `excluded_platforms` - (Optional) A list of platforms explicitly excluded from the policy. Possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`.

---

`grant_controls` block supports the following:

* `operator` - (Required) Defines the relationship of the grant controls. Possible values are: `AND`, `OR`.
* `built_in_controls` - (Required) List of built-in controls required by the policy. Possible values are: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`, `unknownFutureValue`.
* `custom_authentication_factors` - (Optional) List of custom controls IDs required by the policy.
* `terms_of_use` - (Optional) List of terms of use IDs required by the policy.

---

`session_controls` block supports the following:

-> **NOTE:** Only Office 365, Exchange Online and Sharepoint Online support `application_enforced_restrictions`
-> **NOTE:** `persistent_browser` session only works correctly when `included_applications` is set to `All`.

* `application_enforced_restrictions` - (Optional) An `application_enforced_restrictions` block as defined below, which enables session control to enforce application restrictions.
* `cloud_app_security` - (Optional) A `cloud_app_security` block as defined below, which is used to enforce cloud app security checks.
* `persistent_browser` - (Optional) A `persistent_browser` block as defined below, which defines whether to persist cookies or not. All apps should be selected for this session control to work correctly.
* `sign_in_frequency` - (Optional) A `sign_in_frequency` block as defined below, which controls the frequency users will have to sign in.

`application_enforced_restrictions` block supports the following:

* `enabled` - (Optional) enables application enforced restrictions.

`cloud_app_security` block supports the following:

* `enabled` - (Optional) enables conditional access application control.
* `cloud_app_security_type` - (Optional) Defines which policy to use. Possible values are: `mcasConfigured`, `monitorOnly`, `blockDownloads` or `unknownFutureValue`.

`persistent_browser` block supports the following:

* `enabled` - (Optional) enables persistent browser session.
* `mode` - (Optional) Defines which persistent browser mode is enabled. Possible values are: `always`, `never`.

`sign_in_frequency` block supports the following:

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Conditional Access Policy.

## Import

Azure Active Directory Conditional Access Policies can be imported using the `id`, e.g.

```shell
terraform import azuread_conditional_access_policy.my_location 00000000-0000-0000-0000-000000000000
```
