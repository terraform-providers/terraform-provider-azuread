---
subcategory: "Users"
---

# Data Source: azuread_user

Gets information about an Azure Active Directory user.

## Example Usage

```terraform
data "azuread_user" "example" {
  user_principal_name = "user@hashicorp.com"
}
```

## Argument Reference

The following arguments are supported:

* `mail_nickname` - (Optional) The email alias of the user.
* `object_id` - (Optional) The object ID of the user.
* `user_principal_name` - (Optional) The user principal name (UPN) of the user.

~> **NOTE:** One of `user_principal_name`, `object_id` or `mail_nickname` must be specified.

## Attributes Reference

The following attributes are exported:

* `account_enabled` - Whether or not the account is enabled.
* `city` - The city in which the user is located.
* `company_name` - The company name which the user is associated. This property can be useful for describing the company that an external user comes from.
* `country` - The country/region in which the user is located, e.g. `US` or `UK`.
* `department` - The name for the department in which the user works.
* `display_name` - The display name of the user.
* `given_name` - The given name (first name) of the user.
* `job_title` - The user’s job title.
* `mail_nickname` - The email alias of the user.
* `mail` - The primary email address of the user.
* `mobile_phone` - The primary cellular telephone number for the user.
* `office_location` - The office location in the user's place of business.
* `onpremises_immutable_id` - The value used to associate an on-premise Active Directory user account with their Azure AD user object.
* `onpremises_sam_account_name` - The on-premise SAM account name of the user.
* `onpremises_user_principal_name` - The on-premise user principal name of the user.
* `postal_code` - The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code.
* `state` - The state or province in the user's address.
* `street_address` - The street address of the user's place of business.
* `surname` - The user's surname (family name or last name).
* `usage_location` - The usage location of the user.
* `user_principal_name` - The user principal name (UPN) of the user.
* `user_type` - The user type in the directory. Possible values are `Guest` or `Member`.
