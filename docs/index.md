# Azure Active Directory Provider

The Azure Provider can be used to configure infrastructure in [Azure Active Directory](https://azure.microsoft.com/en-us/services/active-directory/) using the [Microsoft Graph](https://docs.microsoft.com/en-us/graph/overview) API. Documentation regarding the [Data Sources](https://www.terraform.io/docs/language/data-sources/index.html) and [Resources](https://www.terraform.io/docs/language/resources/index.html) supported by the Azure Active Directory Provider can be found in the navigation to the left.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the [changelog](https://github.com/hashicorp/terraform-provider-azuread/blob/main/CHANGELOG.md) for version information and release notes.

## Example Usage

```hcl
# Configure Terraform
terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.0.0"
    }
  }
}

# Configure the Azure Active Directory Provider
provider "azuread" {
  tenant_id = "00000000-0000-0000-0000-000000000000"
}

# Retrieve domain information
data "azuread_domains" "example" {
  only_initial = true
}

# Create an application
resource "azuread_application" "example" {
  display_name = "ExampleApp"
}

# Create a service principal
resource "azuread_service_principal" "example" {
  application_id = azuread_application.example.application_id
}

# Create a user
resource "azuread_user" "example" {
  user_principal_name = "ExampleUser@${data.azuread_domains.example.domains.0.domain_name}"
  display_name        = "Example User"
  password            = "..."
}
```

## Authenticating to Azure Active Directory

Terraform supports a number of different methods for authenticating to Azure Active Directory:

* [Authenticating to Azure Active Directory using the Azure CLI](guides/azure_cli.html)
* [Authenticating to Azure Active Directory using Managed Identity](guides/managed_service_identity.html)
* [Authenticating to Azure Active Directory using a Service Principal and a Client Certificate](guides/service_principal_client_certificate.html)
* [Authenticating to Azure Active Directory using a Service Principal and a Client Secret](guides/service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Identity when running Terraform non-interactively (such as when running Terraform in a CI/CD pipeline), and authenticating using the Azure CLI when running Terraform locally.

## Features and Bug Requests

Bugs and feature requests can be reported on the [GitHub issues tracker](https://github.com/hashicorp/terraform-provider-azuread/issues). Please avoid "me too" or "+1" comments. Instead, use a thumbs up [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) on enhancement requests. Provider maintainers will often prioritise work based on the number of thumbs on an issue.

Community input is appreciated on outstanding issues! We love to hear what use cases you have for new features, and want to provide the best possible experience for you using the Azure Active Directory provider.

If you have a bug or feature request without an existing issue:

* if an existing resource or field is working in an unexpected way, [file a bug](https://github.com/hashicorp/terraform-provider-azuread/issues/new?template=bug.md).
* if you'd like the provider to support a new resource or field, [file an enhancement/feature request](https://github.com/hashicorp/terraform-provider-azuread/issues/new?template=enhancement.md).

The provider maintainers will often use the assignee field on an issue to mark who is working on it.

* An issue assigned to an individual maintainer indicates that maintainer is working on the issue
* If you're interested in working on an issue please leave a comment in that issue

---

If you have configuration questions, or general questions about using the provider, try checking out:

* [Terraform's community resources](https://www.terraform.io/docs/extend/community/index.html)
* [HashiCorp support](https://support.hashicorp.com) for Terraform Enterprise customers

## Argument Reference

The following arguments are supported:

* `client_id` - (Optional) The Client ID which should be used when authenticating as a service principal. This can also be sourced from the `ARM_CLIENT_ID` Environment Variable.
* `environment` - (Optional) The Cloud Environment which be used. Possible values are `global`, `germany`, `china`, `usgovernmentl4` and `usgovernmentl5`. Defaults to `global`. This can also be sourced from the `ARM_ENVIRONMENT` environment variable.
* `tenant_id` - (Optional) The Tenant ID which should be used. This can also be sourced from the `ARM_TENANT_ID` Environment Variable.

---

When authenticating as a Service Principal using a Client Certificate, the following fields can be set:

* `client_certificate_password` - (Optional) The password associated with the Client Certificate. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PASSWORD` Environment Variable.
* `client_certificate_path` - (Optional) The path to the Client Certificate associated with the Service Principal which should be used. This can also be sourced from the `ARM_CLIENT_CERTIFICATE_PATH` Environment Variable.

More information on [how to configure a Service Principal using a Client Certificate can be found in this guide](guides/service_principal_client_certificate.html).

---

When authenticating as a Service Principal using a Client Secret, the following fields can be set:

* `client_secret` - (Optional) The Client Secret which should be used. This can also be sourced from the `ARM_CLIENT_SECRET` Environment Variable.

More information on [how to configure a Service Principal using a Client Secret can be found in this guide](guides/service_principal_client_secret.html).

---

When authenticating using Managed Identity, the following fields can be set:

* `msi_endpoint` - (Optional) The path to a custom endpoint for Managed Identity - in most circumstances this should be detected automatically. This can also be sourced from the `ARM_MSI_ENDPOINT` Environment Variable.
* `use_msi` - (Optional) Should Managed Identity be used for authentication? This can also be sourced from the `ARM_USE_MSI` Environment Variable. Defaults to `false`.

More information on [how to configure a Service Principal using Managed Identity can be found in this guide](guides/managed_service_identity.html).

---

For Azure CLI authentication, the following fields can be set:

* `use_cli` - (Optional) Should Azure CLI be used for authentication? This can also be sourced from the `ARM_USE_CLI` Environment Variable. Defaults to `true`.

---

## Advanced Usage

For more advanced scenarios, the following additional arguments are supported:

* `disable_terraform_partner_id` - (Optional) Disable sending the Terraform Partner ID if a custom `partner_id` isn't specified. The default Partner ID allows Microsoft to better understand the usage of Terraform and does not give HashiCorp any direct access to usage information. This can also be sourced from the `ARM_DISABLE_TERRAFORM_PARTNER_ID` environment variable. Defaults to `false`.

* `partner_id` - (Optional) A GUID/UUID that is [registered](https://docs.microsoft.com/azure/marketplace/azure-partner-customer-usage-attribution#register-guids-and-offers) with Microsoft to facilitate partner resource usage attribution. This can also be sourced from the `ARM_PARTNER_ID` Environment Variable.

It's also possible to use multiple Provider blocks within a single Terraform configuration, for example to work with resources across multiple Azure Active Directory Environments - more information can be found [in the documentation for Providers](https://www.terraform.io/docs/configuration/providers.html#multiple-provider-instances).
