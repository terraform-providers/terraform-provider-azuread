---
subcategory: "Applications"
---

# Resource: azuread_application_pre_authorized

Manages client applications that are pre-authorized with the specified permissions to access an application's APIs without requiring user consent.

## Example Usage

```terraform
resource "azuread_application" "authorized" {
  display_name = "example-authorized-app"
}

resource "azuread_application" "authorizer" {
  display_name = "example-authorizing-app"

  api {
    oauth2_permission_scope {
      admin_consent_description  = "Administer the application"
      admin_consent_display_name = "Administer"
      enabled                    = true
      id                         = "ced9c4c3-c273-4f0f-ac71-a20377b90f9c"
      type                       = "Admin"
      value                      = "administer"
    }

    oauth2_permission_scope {
      admin_consent_description  = "Access the application"
      admin_consent_display_name = "Access"
      enabled                    = true
      id                         = "2d5e07ca-664d-4d9b-ad61-ec07fd215213"
      type                       = "User"
      user_consent_description   = "Access the application"
      user_consent_display_name  = "Access"
      value                      = "user_impersonation"
    }
  }
}

resource "azuread_application_pre_authorized" "example" {
  application_object_id = azuread_application.authorizer.object_id
  authorized_app_id     = azuread_application.authorized.application_id
  permission_ids        = ["ced9c4c3-c273-4f0f-ac71-a20377b90f9c", "2d5e07ca-664d-4d9b-ad61-ec07fd215213"]
}
```

## Argument Reference

The following arguments are supported:

* `application_object_id` - (Required) The object ID of the application for which permissions are being authorized. Changing this field forces a new resource to be created.
* `authorized_application_id` - (Optional) The application ID (client ID) of the application being authorized. Changing this field forces a new resource to be created.
* `permission_ids` - (Required) A set of permission scope IDs required by the authorized application. 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*No additional attributes are exported*

## Import

Pre-authorized applications can be imported using the object ID of the authorizing application and the application ID of the application being authorized, e.g.

```shell
terraform import azuread_application_pre_authorized.example 00000000-0000-0000-0000-000000000000/preAuthorizedApplication/11111111-1111-1111-1111-111111111111
```

-> **NOTE:** This ID format is unique to Terraform and is composed of the authorizing application's object ID, the string "preAuthorizedApplication" and the authorized application's application ID (client ID) in the format `{ObjectId}/preAuthorizedApplication/{ApplicationId}`.
