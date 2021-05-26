---
subcategory: "Invitations"
---

# Resource: azuread_invitation

Manages an invitation of a guest user within Azure Active Directory.

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to `User.ReadWrite.All` within the `Microsoft Graph` API.

## Example Usage

*Basic example*

```terraform
resource "azuread_invitation" "example" {
  user_email_address = "jdoe@hashicorp.com"
  redirect_url       = "https://portal.azure.com"
}
```

*Invitation with custom email*

```terraform
resource "azuread_invitation" "example" {
  user_display_name  = "Bob Bobson"
  user_email_address = "bbobson@hashicorp.com"
  redirect_url       = "https://portal.azure.com"

  send_invitation_message = true

  user_message_info {
    cc_recipients           = ["aaliceberg@hashicorp.com"]
    customized_message_body = "Hello there! You are invited to join my Azure tenant !"
    message_language        = "en-US"
  }
}
```

## Argument Reference

The following arguments are supported:

* `redirect_url` - (Required) URL the user should be redirected to once the invitation is redeemed.
* `send_invitation_message` - (Optional) If `true`, an email wille be sent to the user being invited. Must be set to `true` if a `user_message_info` block is specified. Defaults to `false`.
* `user_display_name` - (Optional) Display name of the user being invited.
* `user_email_address` - (Required) Email address of the user being invited.
* `user_message_info` - (Optional) A `user_message_info` block as documented below, which configures the message being sent to the invited user. `send_invitation_message` must be set to `true` if this block is specified.

---

`user_message_info` block supports the following:

* `cc_recipients` - (Optional) Additional recipients the invitation message should be sent to. Currently only 1 additional recipient is supported by Azure.
* `customized_message_body` - (Optional) Customized message body you want to send if you don't want the default message.
* `message_language` - (Optional) Language the message will be sent in. The value specified must be in ISO 639 format. Defaults to `en-US`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the invitation.
* `redeem_url` - URL the user can use to redeem the invitation.
* `user_id` - Object ID of the invited user.