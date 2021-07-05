---
subcategory: "Groups"
---

# Data Source: azuread_groups

Gets Object IDs or Display Names for multiple Azure Active Directory groups.

## Example Usage

```terraform
data "azuread_groups" "groups" {
  display_names = ["group-a", "group-b"]
}
```

## Argument Reference

The following arguments are supported:

* `display_names` - (Optional) The display names of the groups.
* `object_ids` - (Optional) The object IDs of the groups.

~> **NOTE:** Either `display_names` or `object_ids` should be specified. Either of these _may_ be specified as an empty list, in which case no results will be returned.

## Attributes Reference

The following attributes are exported:

* `display_names` - The display names of the groups.
* `object_ids` - The object IDs of the groups.
