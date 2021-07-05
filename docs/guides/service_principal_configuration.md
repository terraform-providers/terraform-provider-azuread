---
page_title: "Configuring a User or Service Principal to manage Azure Active Directory"
subcategory: "Authentication"
---

# Configuring a User or Service Principal for managing Azure Active Directory

Terraform supports a number of different methods for authenticating to Azure:

* [Authenticating to Azure using the Azure CLI](azure_cli.html)
* [Authenticating to Azure using Managed Identity](managed_service_identity.html)
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

We recommend using either a Service Principal or Managed Identity when running Terraform non-interactively (such as when running Terraform in a CI server) - and authenticating using the Azure CLI when running Terraform locally.

## Creating a Service Principal

A Service Principal represents an application within Azure Active Directory whose properties and authentication tokens can be used as the `tenant_id`, `client_id` and `client_secret` fields needed by Terraform.

Depending on how the service principal authenticates to Azure it can be created and configured in a number of different ways: 
* [Authenticating to Azure using a Service Principal and a Client Certificate](service_principal_client_certificate.html)
* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.html)

---

## Azure Active Directory permissions

Now that you have created and authenticated an Application / Service Principal pair, you will need to grant some permissions to administer Azure Active Directory. You can choose either of the following methods to achieve similar results.

### Method 1: API roles (recommended for service principals)

This method involves granting API roles to your Application, and then granting consent for your Service Principal to access the APIs in its own capacity (i.e. not on behalf of a user).

Navigate to the [Azure Active Directory overview][azure-portal-aad-overview] within the [Azure Portal][azure-portal] and select the [App Registrations blade][azure-portal-aad-applications-blade]. Locate your registered Application and click on its display name to manage it.

Go to the API Permissions blade for the Application and click the "Add a permission" button. In the pane that opens, select "Microsoft Graph".

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

Depending on the configuration of your AAD tenant, you may also need to grant the Directory.Read.All and/or Directory.ReadWrite.All roles. If a resource you are using is not shown in the table above, consult the resource documentation.

After assigning permissions, you will need to grant consent for the service principal to utilise them. The easiest way to do this is by clicking the Grant Admin Consent button in the same API Permissions pane. You will need to be signed in to the Portal as a Global Administrator.

The Application now has the necessary permissions to administer your Azure Active Directory tenant.

### Method 2: Directory Roles (recommended for users, i.e. Azure CLI authentication)

With this method, you will assign directory roles to your User Principal, to grant the desired permissions to administer objects in your Azure Active Directory tenant. The following steps may need to be performed by an existing Global Administrator, if that is someone else.

Navigate to the [Azure Active Directory overview][azure-portal-aad-overview] within the [Azure Portal][azure-portal]. Go to the [Roles and Administrators blade][azure-portal-aad-roles-blade].

Locate the role you wish to assign and click on it. Consult the [documentation for administrator role permissions][admin-roles-docs] from Microsoft for more information about the available roles and the permissions they grant.

Click "Add assignments" and type the display name or user principal name of your User in the search box to locate it. If you know the Object ID of the User, verify that it is the same. Select it and click the "Add" button to assign the role.

The choice of which directory roles to assign will be specific to your organisation's security policy. Commonly used roles include:

Role                        | Description
--------------------------- | -----------------------------------------------------------------------------------------------------------------------
`Global Administrator`      | Effective superuser permissions to administer any object in your AAD tenant. Sometimes called `Company Administrator`.
`Global Reader`             | Commonly used in conjunction with other roles to allow reading, but not writing, of directory data.
`Application Administrator` | Create and manage applications, service principals (enterprise applications) and application proxy.
`Groups Administrator`      | Create and manage groups.
`User Administrator`        | Create and manage users _and_ groups.

Once the desired directory role has been assigned, you may need to obtain a new access token in order for the role to take effect. This can be performed by signing out and signing back in to the Azure CLI.

```shell
$ az logout
$ az login --allow-no-subscriptions
```

[admin-roles-docs]: https://docs.microsoft.com/en-us/azure/active-directory/users-groups-roles/directory-assign-admin-roles
[azure-portal]: https://portal.azure.com/
[azure-portal-aad-overview]: https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview
[azure-portal-aad-roles-blade]: https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RolesAndAdministrators
[azure-portal-aad-applications-blade]: https://portal.azure.com#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps
