---
page_title: "Authenticating via Managed Identity"
subcategory: "Authentication"
---

# Authenticating using Managed Identity

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* Authenticating to Azure using Managed Identity (covered in this guide)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

Once you have configured a Service Principal as described in this guide, you should follow the [Configuring a Service Principal for managing Azure Active Directory](service_principal_configuration.html) guide to grant the Service Principal necessary permissions to create and modify Azure Active Directory objects such as users and groups.

## What is a managed identity?

[Managed identities][azure-managed-identities] for Azure resources can be used to authenticate to Azure Active Directory. There are two types of managed identities: system-assigned and user-assigned. This article is based on system-assigned managed identities.

Managed identities work in conjunction with Microsoft Graph, Azure AD, and the Azure Instance Metadata Service (IMDS). Azure resources that support managed identities expose an internal IMDS endpoint that the client can use to request an access token. No credentials are stored on the VM, and the only additional information needed to bootstrap the Terraform connection to Azure is the Tenant ID.

Azure AD creates a tenant-wide security principal when you configure an Azure resource to use a system-assigned managed identity. The configuration process is described in more detail below.  The lifecycle of a system-assigned identity is tied to the resource it is enabled for: it is created when the resource is created and it is automatically removed when the resource is deleted.

Not all Azure services support managed identities, and availability varies by region. Configuration details vary slightly among services. For more information, see [Services that support managed identities for Azure resources][azure-managed-identities-services].

When using a managed identity, you can only manage resources in the tenant where the corresponding service principal is homed. If you need to manage multiple tenants from the same location, we suggest instead using Service Principal Authentication with a [client certificate](service_principal_client_certificate.html) or [client secret](service_principal_client_secret.html) so that you can specify different credentials for each tenant.

## Configuring a VM to use a system-assigned managed identity

The (simplified) Terraform configuration below configures a Virtual Machine with a system-assigned identity, and then outputs the Object ID of the corresponding Service Principal:

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_linux_virtual_machine" "example" {
  name = "test-vm"

  # ...

  identity = {
    type = "SystemAssigned"
  }
}

output "example_msi_object_id" {
  value = azurerm_linux_virtual_machine.test.identity.0.principal_id
}
```

Refer to the [azurerm_linux_virtual_machine][azurerm_linux_virtual_machine] and [azurerm_windows_virtual_machine][azurerm_windows_virtual_machine] documentation for more information on how to use these resources to launch a new virtual machine.

The implicitly created Service Principal should have the same or similar name as your virtual machine. At this point you will need to assign permissions to access Azure Active Directory to create and modify Azure Active Directory objects such as users and groups. See the [Configuring a Service Principal for managing Azure Active Directory][azuread-service-principal-permissions] guide for more information.

## Configuring Managed Identity in Terraform

At this point we assume that managed identity is configured on the resource (e.g. virtual machine) being used, that permissions have been granted, and that you are running Terraform on that resource.

Terraform can be configured to use managed identity for authentication in one of two ways: using Environment Variables or by defining the fields within the Provider block.

You can configure Terraform to use Managed Identity by setting the Environment Variable `ARM_USE_MSI` to `true`; as shown below:

```shell
# sh
export ARM_USE_MSI=true ARM_TENANT_ID=00000000-0000-0000-0000-000000000000

# PowerShell
$env:ARM_USE_MSI = true
$env:ARM_TENANT_ID = 00000000-0000-0000-0000-000000000000
```

Note that when using managed identity for authentication, the tenant ID must also be specified.

-> **Using a Custom MSI Endpoint?** In the unlikely event you're using a custom endpoint for Managed Identity - this can be configured using the `ARM_MSI_ENDPOINT` Environment Variable - however this shouldn't need to be configured in regular use.

See the main provider documentation for more information on [the fields supported in the Provider block][azuread-provider-fields].

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using Managed Identity.

---

It's also possible to enable Managed Identity within the Provider Block:

```hcl
provider "azuread" {
  use_msi   = true
  tenant_id = "00000000-0000-0000-0000-000000000000"
}
```

Remember when using managed identity for authentication, the tenant ID must also be specified.

See the main provider documentation for more information on [the fields supported in the Provider block][azuread-provider-fields].

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using Managed Identity.

Next you should follow the [Configuring a Service Principal for managing Azure Active Directory][azuread-service-principal-configuration] guide to grant the Service Principal necessary permissions to create and modify Azure Active Directory objects such as users and groups.


[azure-managed-identities]: https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview
[azure-managed-identities-services]: https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/services-support-managed-identities
[azuread-provider-fields]: https://registry.terraform.io/providers/hashicorp/azuread/latest/docs#argument-reference
[azuread-service-principal-configuration]: https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration
[azuread-service-principal-permissions]: https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration#method-1-api-roles-recommended-for-service-principals
[azurerm_linux_virtual_machine]: https://www.terraform.io/docs/providers/azurerm/r/linux_virtual_machine.html
[azurerm_windows_virtual_machine]: https://www.terraform.io/docs/providers/azurerm/r/windows_virtual_machine.html
