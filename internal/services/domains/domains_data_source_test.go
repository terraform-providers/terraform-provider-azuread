package domains_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
)

type DomainsDataSource struct{}

func TestAccDomainsDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_domains", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: DomainsDataSource{}.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("domains.0.domain_name").Exists(),
				check.That(data.ResourceName).Key("domains.0.authentication_type").Exists(),
				check.That(data.ResourceName).Key("domains.0.admin_managed").Exists(),
				check.That(data.ResourceName).Key("domains.0.default").Exists(),
				check.That(data.ResourceName).Key("domains.0.initial").Exists(),
				check.That(data.ResourceName).Key("domains.0.root").Exists(),
				check.That(data.ResourceName).Key("domains.0.supported_services.#").Exists(),
				check.That(data.ResourceName).Key("domains.0.verified").Exists(),
			),
		},
	})
}

func TestAccDomainsDataSource_onlyDefault(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_domains", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: DomainsDataSource{}.onlyDefault(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("domains.0.domain_name").Exists(),
				check.That(data.ResourceName).Key("domains.0.admin_managed").Exists(),
				check.That(data.ResourceName).Key("domains.0.default").HasValue("true"),
				check.That(data.ResourceName).Key("domains.0.initial").Exists(),
				check.That(data.ResourceName).Key("domains.0.root").Exists(),
				check.That(data.ResourceName).Key("domains.0.supported_services.#").Exists(),
				check.That(data.ResourceName).Key("domains.0.verified").Exists(),
			),
		},
	})
}

func TestAccDomainsDataSource_onlyInitial(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_domains", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: DomainsDataSource{}.onlyInitial(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("domains.0.domain_name").Exists(),
				check.That(data.ResourceName).Key("domains.0.admin_managed").Exists(),
				check.That(data.ResourceName).Key("domains.0.default").Exists(),
				check.That(data.ResourceName).Key("domains.0.initial").HasValue("true"),
				check.That(data.ResourceName).Key("domains.0.root").Exists(),
				check.That(data.ResourceName).Key("domains.0.supported_services.#").Exists(),
				check.That(data.ResourceName).Key("domains.0.verified").Exists(),
			),
		},
	})
}

func TestAccDomainsDataSource_supportsServices(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azuread_domains", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: DomainsDataSource{}.supportsServices(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("domains.0.domain_name").Exists(),
				check.That(data.ResourceName).Key("domains.0.admin_managed").Exists(),
				check.That(data.ResourceName).Key("domains.0.default").Exists(),
				check.That(data.ResourceName).Key("domains.0.initial").Exists(),
				check.That(data.ResourceName).Key("domains.0.root").Exists(),
				check.That(data.ResourceName).Key("domains.0.supported_services.#").Exists(),
				check.That(data.ResourceName).Key("domains.0.verified").Exists(),
			),
		},
	})
}

func (DomainsDataSource) basic() string {
	return `data "azuread_domains" "test" {}`
}

func (DomainsDataSource) onlyDefault() string {
	return `
data "azuread_domains" "test" {
  only_default = true
}
`
}

func (DomainsDataSource) onlyInitial() string {
	return `
data "azuread_domains" "test" {
  only_initial = true
}
`
}

func (DomainsDataSource) supportsServices() string {
	return `
data "azuread_domains" "test" {
  supports_services = ["Email", "OfficeCommunicationsOnline"]
}
`
}
