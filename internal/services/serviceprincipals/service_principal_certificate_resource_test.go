package serviceprincipals_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/serviceprincipals/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

// To create test certificates:
// openssl req -subj '/CN=hashicorptest/O=HashiCorp, Inc./ST=CA/C=US' -new -newkey rsa:2048 -sha256 -days 3650 -nodes -x509 -keyout server.key -out server.crt
// grep -v \\----- server.crt >server.b64
// cat server.b64 | base64 -d | xxd -p

const servicePrincipalCertificatePem string = `-----BEGIN CERTIFICATE-----
MIIDFDCCAfwCCQCvHp+vopfOOTANBgkqhkiG9w0BAQsFADBMMRYwFAYDVQQDDA1o
YXNoaWNvcnB0ZXN0MRgwFgYDVQQKDA9IYXNoaUNvcnAsIEluYy4xCzAJBgNVBAgM
AkNBMQswCQYDVQQGEwJVUzAeFw0yMTAzMDkxMTAyMTNaFw0zMTAzMDcxMTAyMTNa
MEwxFjAUBgNVBAMMDWhhc2hpY29ycHRlc3QxGDAWBgNVBAoMD0hhc2hpQ29ycCwg
SW5jLjELMAkGA1UECAwCQ0ExCzAJBgNVBAYTAlVTMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAlVmb5pmoASvZ5pxD6CEBiPYqADb7teCHV54RRwv1aJjS
eiPUW/1WNQooIQF0M0yzFdHmwx3HSoxCkQwwxVMAPsuqFJVabs/eAr41NpxQCncb
i+vKlbmaAWbaIdidxeUe1jXB2N0YXRCg7Ps8IGA0UochvRGypfciy4k6/xEfrrQP
FlrPeDeaurNUjJ4IotTBLzWNAX9nT1HKzvljYNg4A0PwuzPNOmgxUSpAeiPbDoQo
D/YcQUKWzBlW8qt9ZnuRMGNi6V2fnQeTLblfsheaavXyP11syJ9owz6mDffZELHd
SYC7j2EOqG+Pndd55MLOac8cF4D9Y91PkLKKjNIrWwIDAQABMA0GCSqGSIb3DQEB
CwUAA4IBAQBlVJLn17BFmigbqS8JIx0/RTbGokRoLKdg7SZAQJWn20jDtunSo+sp
ZzuZ4uS8WbgZ+SFD1rrQy3s0F9HssZFBwDGyn31z/sGjkwWpoAP65v1DCaNzmAsz
xMNijhYlShv61g2IEO9Q98bgBW9LNwmJRGnGxz0ufzeZuUr9IV9EjeoJCKPIbwJC
lab0Ty/kRC13JgNhHtNFwYVwK6NDt46IRsjxqWQ6bVakrEROlfuoY8sxUjunj+hB
2vZTkZKaPc0sFvUQjNHxHX4jMeTwCopQCo+qF3lPde+G7C1MNf30kDZlks++GLNs
0/0Ayfjh6JllWqW482dIIqMErl6s5DuK
-----END CERTIFICATE-----`

const servicePrincipalCertificateBase64 string = `MIIDFDCCAfwCCQCvHp+vopfOOTANBgkqhkiG9w0BAQsFADBMMRYwFAYDVQQD
DA1oYXNoaWNvcnB0ZXN0MRgwFgYDVQQKDA9IYXNoaUNvcnAsIEluYy4xCzAJ
BgNVBAgMAkNBMQswCQYDVQQGEwJVUzAeFw0yMTAzMDkxMTAyMTNaFw0zMTAz
MDcxMTAyMTNaMEwxFjAUBgNVBAMMDWhhc2hpY29ycHRlc3QxGDAWBgNVBAoM
D0hhc2hpQ29ycCwgSW5jLjELMAkGA1UECAwCQ0ExCzAJBgNVBAYTAlVTMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlVmb5pmoASvZ5pxD6CEB
iPYqADb7teCHV54RRwv1aJjSeiPUW/1WNQooIQF0M0yzFdHmwx3HSoxCkQww
xVMAPsuqFJVabs/eAr41NpxQCncbi+vKlbmaAWbaIdidxeUe1jXB2N0YXRCg
7Ps8IGA0UochvRGypfciy4k6/xEfrrQPFlrPeDeaurNUjJ4IotTBLzWNAX9n
T1HKzvljYNg4A0PwuzPNOmgxUSpAeiPbDoQoD/YcQUKWzBlW8qt9ZnuRMGNi
6V2fnQeTLblfsheaavXyP11syJ9owz6mDffZELHdSYC7j2EOqG+Pndd55MLO
ac8cF4D9Y91PkLKKjNIrWwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBlVJLn
17BFmigbqS8JIx0/RTbGokRoLKdg7SZAQJWn20jDtunSo+spZzuZ4uS8WbgZ
+SFD1rrQy3s0F9HssZFBwDGyn31z/sGjkwWpoAP65v1DCaNzmAszxMNijhYl
Shv61g2IEO9Q98bgBW9LNwmJRGnGxz0ufzeZuUr9IV9EjeoJCKPIbwJClab0
Ty/kRC13JgNhHtNFwYVwK6NDt46IRsjxqWQ6bVakrEROlfuoY8sxUjunj+hB
2vZTkZKaPc0sFvUQjNHxHX4jMeTwCopQCo+qF3lPde+G7C1MNf30kDZlks++
GLNs0/0Ayfjh6JllWqW482dIIqMErl6s5DuK`

const servicePrincipalCertificateHex string = `30820314308201fc020900af1e9fafa297ce39300d06092a864886f70d01010b0500304c3116301406035504030c0d6861736869636f72707465737431183016060355040a0c0f4861736869436f72702c20496e632e310b300906035504080c024341310b3009060355040613025553301e170d3231303330393131303231335a170d3331303330373131303231335a304c3116301406035504030c0d6861736869636f72707465737431183016060355040a0c0f4861736869436f72702c20496e632e310b300906035504080c024341310b300906035504061302555330820122300d06092a864886f70d01010105000382010f003082010a028201010095599be699a8012bd9e69c43e8210188f62a0036fbb5e087579e11470bf56898d27a23d45bfd56350a28210174334cb315d1e6c31dc74a8c42910c30c553003ecbaa14955a6ecfde02be35369c500a771b8bebca95b99a0166da21d89dc5e51ed635c1d8dd185d10a0ecfb3c206034528721bd11b2a5f722cb893aff111faeb40f165acf78379abab3548c9e08a2d4c12f358d017f674f51cacef96360d8380343f0bb33cd3a6831512a407a23db0e84280ff61c414296cc1956f2ab7d667b91306362e95d9f9d07932db95fb2179a6af5f23f5d6cc89f68c33ea60df7d910b1dd4980bb8f610ea86f8f9dd779e4c2ce69cf1c1780fd63dd4f90b28a8cd22b5b0203010001300d06092a864886f70d01010b05000382010100655492e7d7b0459a281ba92f09231d3f4536c6a244682ca760ed26404095a7db48c3b6e9d2a3eb29673b99e2e4bc59b819f92143d6bad0cb7b3417d1ecb19141c031b29f7d73fec1a39305a9a003fae6fd4309a373980b33c4c3628e16254a1bfad60d8810ef50f7c6e0056f4b3709894469c6c73d2e7f3799b94afd215f448dea0908a3c86f024295a6f44f2fe4442d772603611ed345c185702ba343b78e8846c8f1a9643a6d56a4ac444e95fba863cb31523ba78fe841daf65391929a3dcd2c16f5108cd1f11d7e2331e4f00a8a500a8faa17794f75ef86ec2d4c35fdf490366592cfbe18b36cd3fd00c9f8e1e899655aa5b8f3674822a304ae5eace43b8a`

type ServicePrincipalCertificateResource struct{}

func TestAccServicePrincipalCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_service_principal_certificate", "test")
	endDate := time.Now().AddDate(0, 3, 27).UTC().Format(time.RFC3339)
	r := ServicePrincipalCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, endDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_id").Exists(),
			),
		},
		data.ImportStep("encoding", "value"),
	})
}

func TestAccServicePrincipalCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_service_principal_certificate", "test")
	startDate := time.Now().AddDate(0, 0, 7).UTC().Format(time.RFC3339)
	endDate := time.Now().AddDate(0, 3, 27).UTC().Format(time.RFC3339)
	r := ServicePrincipalCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data, startDate, endDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_id").Exists(),
			),
		},
		data.ImportStep("encoding", "end_date_relative", "value"),
	})
}

func TestAccServicePrincipalCertificate_base64Cert(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_service_principal_certificate", "test")
	endDate := time.Now().AddDate(0, 3, 27).UTC().Format(time.RFC3339)
	r := ServicePrincipalCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.base64Cert(data, endDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_id").Exists(),
			),
		},
		data.ImportStep("encoding", "end_date_relative", "value"),
	})
}

func TestAccServicePrincipalCertificate_hexCert(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_service_principal_certificate", "test")
	endDate := time.Now().AddDate(0, 3, 27).UTC().Format(time.RFC3339)
	r := ServicePrincipalCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.hexCert(data, endDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_id").Exists(),
			),
		},
		data.ImportStep("encoding", "end_date_relative", "value"),
	})
}

func TestAccServicePrincipalCertificate_relativeEndDate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_service_principal_certificate", "test")
	r := ServicePrincipalCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.relativeEndDate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_id").Exists(),
				check.That(data.ResourceName).Key("end_date").Exists(),
			),
		},
		data.ImportStep("encoding", "end_date_relative", "value"),
	})
}

func TestAccServicePrincipalCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_service_principal_certificate", "test")
	endDate := time.Now().AddDate(0, 3, 27).UTC().Format(time.RFC3339)
	r := ServicePrincipalCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, endDate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("key_id").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport(data, endDate)),
	})
}

func (r ServicePrincipalCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.ServicePrincipals.ServicePrincipalsClient
	client.BaseClient.DisableRetries = true

	id, err := parse.CertificateID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing Service Principal Certificate ID: %v", err)
	}

	servicePrincipal, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Service Principal with object ID %q does not exist", id.ObjectId)
		}
		return nil, fmt.Errorf("failed to retrieve Service Principal with object ID %q: %+v", id.ObjectId, err)
	}

	if servicePrincipal.KeyCredentials != nil {
		for _, cred := range *servicePrincipal.KeyCredentials {
			if cred.KeyId != nil && *cred.KeyId == id.KeyId {
				return utils.Bool(true), nil
			}
		}
	}

	return nil, fmt.Errorf("Key Credential %q was not found for Service Principal %q", id.KeyId, id.ObjectId)
}

func (ServicePrincipalCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_application" "test" {
  display_name = "acctestServicePrincipal-%[1]d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}
`, data.RandomInteger)
}

func (r ServicePrincipalCertificateResource) basic(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_service_principal_certificate" "test" {
  service_principal_id = azuread_service_principal.test.id
  type                 = "AsymmetricX509Cert"
  end_date             = "%[2]s"
  value                = <<EOT
%[3]s
EOT
}
`, r.template(data), endDate, servicePrincipalCertificatePem)
}

func (r ServicePrincipalCertificateResource) complete(data acceptance.TestData, startDate, endDate string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_service_principal_certificate" "test" {
  service_principal_id = azuread_service_principal.test.id
  key_id               = "%[2]s"
  type                 = "AsymmetricX509Cert"
  start_date           = "%[3]s"
  end_date             = "%[4]s"
  encoding             = "pem"
  value                = <<EOT
%[5]s
EOT
}
`, r.template(data), data.RandomID, startDate, endDate, servicePrincipalCertificatePem)
}

func (r ServicePrincipalCertificateResource) base64Cert(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_service_principal_certificate" "test" {
  service_principal_id = azuread_service_principal.test.id
  type                 = "AsymmetricX509Cert"
  end_date             = "%[2]s"
  encoding             = "base64"
  value                = <<EOT
%[3]s
EOT
}
`, r.template(data), endDate, servicePrincipalCertificateBase64)
}

func (r ServicePrincipalCertificateResource) hexCert(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_service_principal_certificate" "test" {
  service_principal_id = azuread_service_principal.test.id
  type                 = "AsymmetricX509Cert"
  end_date             = "%[2]s"
  encoding             = "hex"
  value                = <<EOT
%[3]s
EOT
}
`, r.template(data), endDate, servicePrincipalCertificateHex)
}

func (r ServicePrincipalCertificateResource) relativeEndDate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_service_principal_certificate" "test" {
  service_principal_id = azuread_service_principal.test.id
  end_date_relative    = "2280h"
  type                 = "AsymmetricX509Cert"
  value                = <<EOT
%[2]s
EOT
}
`, r.template(data), servicePrincipalCertificatePem)
}

func (r ServicePrincipalCertificateResource) requiresImport(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_service_principal_certificate" "import" {
  service_principal_id = azuread_service_principal_certificate.test.service_principal_id
  key_id               = azuread_service_principal_certificate.test.key_id
  type                 = azuread_service_principal_certificate.test.type
  end_date             = azuread_service_principal_certificate.test.end_date
  value                = azuread_service_principal_certificate.test.value
}
`, r.basic(data, endDate))
}
