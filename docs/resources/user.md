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
* `age_group` - (Optional) The age group of the user. Supported values are `Adult`, `NotAdult` and `Minor`. Omit this property or specify a blank string to unset.
* `business_phones` - (Optional) A list of telephone numbers for the user. Only one number can be set for this property. Read-only for users synced with Azure AD Connect.
* `city` - (Optional) The city in which the user is located.
* `company_name` - (Optional) The company name which the user is associated. This property can be useful for describing the company that an external user comes from.
* `consent_provided_for_minor` - (Optional) Whether consent has been obtained for minors. Supported values are `Granted`, `Denied` and `NotRequired`. Omit this property or specify a blank string to unset.
* `country` - (Optional) The country/region in which the user is located, e.g. `US` or `UK`.
* `department` - (Optional) The name for the department in which the user works.
* `display_name` - (Required) The name to display in the address book for the user.
* `employee_id` - (Optional) The employee identifier assigned to the user by the organisation.
* `fax_number` - (Optional) The fax number of the user.
* `force_password_change` - (Optional) Whether the user is forced to change the password during the next sign-in. Only takes effect when also changing the password. Defaults to `false`.
* `given_name` - (Optional) The given name (first name) of the user.
* `job_title` - (Optional) The user’s job title.
* `mail` - (Optional) The SMTP address for the user. This property cannot be unset once specified.
* `mail_nickname` - (Optional) The mail alias for the user. Defaults to the user name part of the user principal name (UPN).
* `mobile_phone` - (Optional) The primary cellular telephone number for the user.
* `office_location` - (Optional) The office location in the user's place of business.
* `other_mails` - (Optional) A list of additional email addresses for the user.
* `onpremises_immutable_id` - (Optional) The value used to associate an on-premise Active Directory user account with their Azure AD user object. This must be specified if you are using a federated domain for the user's `user_principal_name` property when creating a new user account.
* `password` - (Optional) The password for the user. The password must satisfy minimum requirements as specified by the password policy. The maximum length is 256 characters. This property is required when creating a new user.
* `postal_code` - (Optional) The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code.
* `preferred_language` - (Optional) The user's preferred language, in ISO 639-1 notation.
* `show_in_address_list` - (Optional) Whether or not the Outlook global address list should include this user. Defaults to `true`.
* `state` - (Optional) The state or province in the user's address.
* `street_address` - (Optional) The street address of the user's place of business.
* `surname` - (Optional) The user's surname (family name or last name).
* `usage_location` - (Optional) The usage location of the user. Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries. The usage location is a two letter country code (ISO standard 3166). Examples include: `NO`, `JP`, and `GB`. Cannot be reset to null once set. 
* `user_principal_name` - (Required) The user principal name (UPN) of the user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `creation_type` - Indicates whether the user account was created as a regular school or work account (`null`), an external account (`Invitation`), a local account for an Azure Active Directory B2C tenant (`LocalAccount`) or self-service sign-up using email verification (`EmailVerified`).
* `external_user_state` - For an external user invited to the tenant, this property represents the invited user's invitation status. Possible values are `PendingAcceptance` or `Accepted`.
* `im_addresses` - A list of instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user.
* `object_id` - The object ID of the user.
* `onpremises_distinguished_name` - The on-premises distinguished name (DN) of the user, synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_domain_name` - The on-premises FQDN, also called dnsDomainName, synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_sam_account_name` - The on-premise SAM account name of the user.
* `onpremises_security_identifier` - The on-premises security identifier (SID), synchronised from the on-premises directory when Azure AD Connect is used.
* `onpremises_sync_enabled` - Whether this user is synchronised from an on-premises directory (`true`), no longer synchronised (`false`), or has never been synchronised (`null`).
* `onpremises_user_principal_name` - The on-premise user principal name of the user.
* `proxy_addresses` - List of email addresses for the user that direct to the same mailbox.
* `user_type` - The user type in the directory. Possible values are `Guest` or `Member`.

## Import

Users can be imported using their object ID, e.g.

```shell
terraform import azuread_user.my_user 00000000-0000-0000-0000-000000000000
```
