---
subcategory: "Base"
---

# Data Source: azuread_client_config

Use this data source to access the configuration of the AzureAD provider.

## Example Usage

```hcl
data "azuread_client_config" "current" {
}

output "account_id" {
  value = data.azuread_client_config.current.client_id
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

* `client_id` - The client ID (application ID) linked to the authenticated principal, or the application used for delegated authentication.
* `object_id` - The object ID of the authenticated principal.
* `tenant_id` - The tenant ID of the authenticated principal.
