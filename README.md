# Terraform Provider for Azure Active Directory

**NOTE:** Version 1.0 and above of this provider requires Terraform 0.12 or later.

- [Terraform Website](https://www.terraform.io)
- [AzureAD Provider Documentation](https://terraform.io/docs/providers/azuread/)
- [AzureAD Provider Usage Examples](https://github.com/hashicorp/terraform-provider-azuread/tree/main/examples)
- [Slack Workspace for Contributors](https://terraform-azure.slack.com) ([Request Invite](https://join.slack.com/t/terraform-azure/shared_invite/enQtNDMzNjQ5NzcxMDc3LWNiY2ZhNThhNDgzNmY0MTM0N2MwZjE4ZGU0MjcxYjUyMzRmN2E5NjZhZmQ0ZTA1OTExMGNjYzA4ZDkwZDYxNDE))


## Usage Example

```
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

  # NOTE: Environment Variables can also be used for Service Principal authentication
  # Terraform also supports authenticating via the Azure CLI too.
  # See official docs for more info: https://registry.terraform.io/providers/hashicorp/azuread/latest/docs

  # client_id     = "..."
  # client_secret = "..."
  # tenant_id     = "..."
}

# Retrieve domain information
data "azuread_domains" "example" {
  only_initial = true
}

# Create an application
resource "azuread_application" "example" {
  name = "ExampleApp"
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

Further [usage documentation](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs) is available on the Terraform website.


## Developer Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x or later
- [Go](https://golang.org/doc/install) 1.16.x (to build the provider plugin)

If you're building on Windows, you will also need:
- [Git Bash for Windows](https://git-scm.com/download/win)
- [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

For *GNU32 Make*, make sure its bin path is added to your PATH environment variable.

For *Git Bash for Windows*, at the step of "Adjusting your PATH environment", please choose "Use Git and optional Unix tools from Windows Command Prompt".


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Clone the repository to: `$GOPATH/src/github.com/hashicorp/terraform-provider-azuread`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/hashicorp/terraform-provider-azuread
```

Change to the clone directory and run `make tools` to install the dependent tooling needed to test and build the provider.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make tools
...
$ make build
...
$ $GOPATH/bin/terraform-provider-azuread
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

The majority of tests in the provider are Acceptance Tests - which provisions real resources in Azure. It's possible to run the entire acceptance test suite by running `make testacc` - however it's likely you'll want to run a subset, which you can do using a prefix, by running:

```
make testacc TESTARGS='-run=TestAccApplication'
```

The following ENV variables must be set in your shell prior to running acceptance tests:
- ARM_CLIENT_ID
- ARM_CLIENT_SECRET
- ARM_TENANT_ID
- ARM_TEST_LOCATION
- ARM_TEST_LOCATION_ALT

*NOTE:* Acceptance tests create real resources, and may cost money to run.
