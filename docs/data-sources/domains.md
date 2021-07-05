---
subcategory: "Domains"
---

# Data Source: azuread_domains

Use this data source to access information about existing Domains within Azure Active Directory.

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to `Directory.Read.All` within the `Windows Azure Active Directory` API.

## Example Usage

```terraform
data "azuread_domains" "aad_domains" {}

output "domains" {
  value = data.azuread_domains.aad_domains.domains
}
```

## Argument Reference

* `admin_managed` - (Optional) Set to `true` to only return domains whose DNS is managed by Microsoft 365. Defaults to `false`.
* `include_unverified` - (Optional) Set to `true` if unverified Azure AD domains should be included. Defaults to `false`.
* `only_default` - (Optional) Set to `true` to only return the default domain.
* `only_initial` - (Optional) Set to `true` to only return the initial domain, which is your primary Azure Active Directory tenant domain. Defaults to `false`.
* `only_root` - (Optional) Set to `true` to only return verified root domains. Excludes subdomains and unverified domains.
* `supports_services` - (Optional) A list of supported services that must be supported by a domain. Possible values include `Email`, `Sharepoint`, `EmailInternalRelayOnly`, `OfficeCommunicationsOnline`, `SharePointDefaultDomain`, `FullRedelegation`, `SharePointPublic`, `OrgIdAuthentication`, `Yammer` and `Intune`.

~> **NOTE:** If `include_unverified` is set to `true` you cannot specify `only_default` or `only_initial`. Additionally, you cannot combine `only_default` with `only_initial`.

## Attributes Reference

* `domains` - A list of tenant domains. Each `domain` object provides the attributes documented below.

`domain` object exports the following:

* `admin_managed` - Whether the DNS for the domain is managed by Microsoft 365.
* `authentication_type` - The authentication type of the domain. Possible values include `Managed` or `Federated`.
* `domain_name` - The name of the domain.
* `default` - Whether this is the default domain that is used for user creation.
* `initial` - Whether this is the initial domain created by Azure Active Directory.
* `root` - Whether the domain is a verified root domain (not a subdomain).
* `verified` - Whether the domain has completed domain ownership verification.
* `supported_services` - A list of capabilities / services supported by the domain. Possible values include `Email`, `Sharepoint`, `EmailInternalRelayOnly`, `OfficeCommunicationsOnline`, `SharePointDefaultDomain`, `FullRedelegation`, `SharePointPublic`, `OrgIdAuthentication`, `Yammer` and `Intune`.