---
subcategory: "Users"
---

# Resource: azuread_user

Manages a user within Azure Active Directory.

## Example Usage

```terraform
resource "azuread_user" "example" {
  user_principal_name = "jdoe@hashicorp.com"
  display_name        = "J. Doe"
  mail_nickname       = "jdoe"
  password            = "SecretP@sswd99!"
}
```

## Argument Reference

The following arguments are supported:

* `account_enabled` - (Optional) Whether or not the account should be enabled.
* `city` - (Optional) The city in which the user is located.
* `company_name` - (Optional) The company name which the user is associated. This property can be useful for describing the company that an external user comes from.
* `country` - (Optional) The country/region in which the user is located, e.g. `US` or `UK`.
* `department` - (Optional) The name for the department in which the user works.
* `display_name` - (Required) The name to display in the address book for the user.
* `force_password_change` - (Optional) Whether the user is forced to change the password during the next sign-in. Only takes effect when also changing the password. Defaults to `false`.
* `given_name` - (Optional) The given name (first name) of the user.
* `job_title` - (Optional) The user’s job title.
* `mail_nickname` - (Optional) The mail alias for the user. Defaults to the user name part of the user principal name (UPN).
* `mobile_phone` - (Optional) The primary cellular telephone number for the user.
* `office_location` - (Optional) The office location in the user's place of business.
* `onpremises_immutable_id` - (Optional) The value used to associate an on-premise Active Directory user account with their Azure AD user object. This must be specified if you are using a federated domain for the user's `user_principal_name` property when creating a new user account.
* `password` - (Optional) The password for the user. The password must satisfy minimum requirements as specified by the password policy. The maximum length is 256 characters. This property is required when creating a new user.
* `postal_code` - (Optional) The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code.
* `state` - (Optional) The state or province in the user's address.
* `street_address` - (Optional) The street address of the user's place of business.
* `surname` - (Optional) The user's surname (family name or last name).
* `usage_location` - (Optional) The usage location of the user. Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries. The usage location is a two letter country code (ISO standard 3166). Examples include: `NO`, `JP`, and `GB`. Cannot be reset to null once set. 
* `user_principal_name` - (Required) The user principal name (UPN) of the user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mail` - The primary email address of the user.
* `object_id` - The object ID of the user.
* `onpremises_sam_account_name` - The on-premise SAM account name of the user.
* `onpremises_user_principal_name` - The on-premise user principal name of the user.
* `user_type` - The user type in the directory. Possible values are `Guest` or `Member`.

## Import

Users can be imported using their object ID, e.g.

```shell
terraform import azuread_user.my_user 00000000-0000-0000-0000-000000000000
```
