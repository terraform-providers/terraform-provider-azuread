---
subcategory: "Applications"
---

# Resource: azuread_application_certificate

Manages a certificate associated with an application within Azure Active Directory. These are also referred to as client certificates during authentication.

## Example Usage

*Using a PEM certificate*

```terraform
resource "azuread_application" "example" {
  display_name = "example"
}

resource "azuread_application_certificate" "example" {
  application_object_id = azuread_application.example.id
  type                  = "AsymmetricX509Cert"
  value                 = file("cert.pem")
  end_date              = "2021-05-01T01:02:03Z"
}
```

*Using a DER certificate*

```terraform
resource "azuread_application" "example" {
  display_name = "example"
}

resource "azuread_application_certificate" "example" {
  application_object_id = azuread_application.example.id
  type                  = "AsymmetricX509Cert"
  encoding              = "base64"
  value                 = base64encode(file("cert.der"))
  end_date              = "2021-05-01T01:02:03Z"
}
```

### Using a certificate from Azure Key Vault

```terraform
resource "azurerm_key_vault_certificate" "example" {
  name         = "generated-cert"
  key_vault_id = azurerm_key_vault.example.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      extended_key_usage = ["1.3.6.1.5.5.7.3.2"]

      key_usage = [
        "dataEncipherment",
        "digitalSignature",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject_alternative_names {
        dns_names = ["internal.contoso.com", "domain.hello.world"]
      }

      subject            = "CN=${azuread_application.example.name}"
      validity_in_months = 12
    }
  }
}

resource "azuread_application" "example" {
  display_name = "example"
}

resource "azuread_application_certificate" "example" {
  application_object_id = azuread_application.example.id
  type                  = "AsymmetricX509Cert"
  encoding              = "hex"
  value                 = azurerm_key_vault_certificate.example.certificate_data
  end_date              = azurerm_key_vault_certificate.example.certificate_attribute[0].expires
  start_date            = azurerm_key_vault_certificate.example.certificate_attribute[0].not_before
}
```

## Argument Reference

The following arguments are supported:

* `application_object_id` - (Required) The object ID of the application for which this certificate should be created. Changing this field forces a new resource to be created.
* `encoding` - (Optional) Specifies the encoding used for the supplied certificate data. Must be one of `pem`, `base64` or `hex`. Defaults to `pem`.

-> **NOTE:** The `hex` encoding option is useful for consuming certificate data from the [azurerm_key_vault_certificate](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault_certificate) resource.

* `end_date` - (Optional) The end date until which the certificate is valid, formatted as an RFC3339 date string (e.g. `2018-01-01T01:02:03Z`). If omitted, the API will decide a suitable expiry date, which is typically around 2 years from the start date. Changing this field forces a new resource to be created.
* `end_date_relative` - (Optional) A relative duration for which the certificate is valid until, for example `240h` (10 days) or `2400h30m`. Changing this field forces a new resource to be created.

~> **NOTE:** One of `end_date` or `end_date_relative` must be set. The maximum duration is enforced by Azure AD.

* `key_id` - (Optional) A UUID used to uniquely identify this certificate. If omitted, a random UUID will be automatically generated. Changing this field forces a new resource to be created.
* `start_date` - (Optional) The start date from which the certificate is valid, formatted as an RFC3339 date string (e.g. `2018-01-01T01:02:03Z`). If this isn't specified, the current date and time are used.  Changing this field forces a new resource to be created.
* `type` - (Required) The type of key/certificate. Must be one of `AsymmetricX509Cert` or `Symmetric`. Changing this fields forces a new resource to be created.
* `value` - (Required) The certificate data, which can be PEM encoded, base64 encoded DER or hexadecimal encoded DER. See also the `encoding` argument.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*No additional attributes are exported*

## Import

Certificates can be imported using the object ID of the associated application and the key ID of the certificate credential, e.g.

```shell
terraform import azuread_application_certificate.test 00000000-0000-0000-0000-000000000000/certificate/11111111-1111-1111-1111-111111111111
```

-> **NOTE:** This ID format is unique to Terraform and is composed of the application's object ID, the string "certificate" and the certificate's key ID in the format `{ObjectId}/certificate/{CertificateKeyId}`.
