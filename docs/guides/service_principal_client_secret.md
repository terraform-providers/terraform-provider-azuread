---
page_title: "Authenticating via a Service Principal and a Client Secret"
subcategory: "Authentication"
---

# Authenticating using a Service Principal with a Client Secret

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Identity](managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* Authenticating to Azure using a Service Principal and a Client Secret (covered in this guide)

---

We recommend using either a Service Principal or Managed Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

Once you have configured a Service Principal as described in this guide, you should follow the [Configuring a Service Principal for managing Azure Active Directory](service_principal_configuration.html) guide to grant the Service Principal necessary permissions to create and modify Azure Active Directory objects such as users and groups.

---

## Setting up an Application and Service Principal

A Service Principal is a security principal within Azure Active Directory which can be granted permissions to manage objects in Azure Active Directory. To authenticate with a Service Principal, you will need to create an Application object within Azure Active Directory, which you will use as a means of authentication, either using a [Client Certificate](service_principal_client_certificate.html) or a Client Secret (which is documented in this guide). This can be done using the Azure Portal.

This guide will cover how to create an Application and linked Service Principal, and then how to generate a Client Secret for the Application so that it can be used for authentication.

---

## Creating the Application and Service Principal

We're going to create the Application in the Azure Portal - to do this navigate to the [**Azure Active Directory** overview][azure-portal-aad-overview] within the [Azure Portal][azure-portal] - then select the [**App Registrations** blade][azure-portal-applications-blade]. Click the **New registration** button at the top to add a new Application within Azure Active Directory. On this page, set the following values then press **Create**:

- **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
- **Supported Account Types** - this should be set to "Accounts in this organisational directory only (single tenant)"
- **Redirect URI** - you should choose "Web" in for the URI type. the actual value can be left blank

At this point the newly created Azure Active Directory application should be visible on-screen - if it's not, navigate to the [**App Registration** blade][azure-portal-applications-blade] and select the Azure Active Directory application.

At the top of this page, you'll need to take note of the "Application (client) ID" and the "Directory (tenant) ID", which you can use for the values of `client_id` and `tenant_id` respectively.

### Generating a Client Secret for the Azure Active Directory Application

Now that the Azure Active Directory Application exists we can create a Client Secret which can be used for authentication - to do this select **Certificates & secrets**. This screen displays the Certificates and Client Secrets (i.e. passwords) which are associated with this Azure Active Directory Application.

Click the "New client secret" button, then enter a short description, choose an expiry period and click "Add". Once the Client Secret has been generated it will be displayed on screen - _the secret is only displayed once_ so **be sure to copy it now** (otherwise you will need to regenerate a new one). This is the `client_secret` you will need.

---

## Configuring Terraform to use the Client Certificate

Now we have obtained the necessary credentials, it's possible to configure Terraform in a few different ways.

## Environment Variables

Our recommended approach is storing the credentials as Environment Variables, for example:

```bash
# sh
$ export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_CLIENT_SECRET="MyCl1eNtSeCr3t"
$ export ARM_TENANT_ID="10000000-2000-3000-4000-500000000000"

# PowerShell
$env:ARM_CLIENT_ID = 00000000-0000-0000-0000-000000000000
$env:ARM_CLIENT_SECRET = 03MWPmH2W3i008UVcucO1E1vifY_bR
$env:ARM_TENANT_ID = 10000000-2000-3000-4000-500000000000
```

At this point running either `terraform plan` or `terraform apply` should allow Terraform to authenticate using the Client Secret.

Next you should follow the [Configuring a Service Principal for managing Azure Active Directory](service_principal_configuration.html) guide to grant the Service Principal necessary permissions to create and modify Azure Active Directory objects such as users and groups.

## Provider Block

It's also possible to configure these variables either directly, or from variables, in your provider block, like so:

~> We recommend not defining these variables in-line since they could easily be checked into Source Control.

```hcl
variable "client_secret" {}

provider "azuread" {
  client_id     = "00000000-0000-0000-0000-000000000000"
  client_secret = var.client_secret
  tenant_id     = "10000000-2000-3000-4000-500000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../index.html#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to authenticate using the Client Secret.

Next you should follow the [Configuring a Service Principal for managing Azure Active Directory](service_principal_configuration.html) guide to grant the Service Principal necessary permissions to create and modify Azure Active Directory objects such as users and groups.

[azure-portal]: https://portal.azure.com/
[azure-portal-aad-overview]: https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview
[azure-portal-applications-blade]: https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview
