package namedlocations_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type NamedLocationResource struct{}

func TestAccNamedLocation_basicIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_named_location", "test")
	r := NamedLocationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNamedLocation_completeIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_named_location", "test")
	r := NamedLocationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNamedLocation_updateIP(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_named_location", "test")
	r := NamedLocationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicIP(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNamedLocation_basicCountry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_named_location", "test")
	r := NamedLocationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCountry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNamedLocation_completeCountry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_named_location", "test")
	r := NamedLocationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeCountry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNamedLocation_updateCountry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_named_location", "test")
	r := NamedLocationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCountry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeCountry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicCountry(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r NamedLocationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	namedLocation, status, err := clients.NamedLocations.MsClient.Get(ctx, state.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Named Location with object ID %q does not exist", state.ID)
		}
		return nil, fmt.Errorf("failed to retrieve Named Location with object ID %q: %+v", state.ID, err)
	}
	ipnl, ok1 := (*namedLocation).(msgraph.IPNamedLocation)
	cnl, ok2 := (*namedLocation).(msgraph.CountryNamedLocation)
	if ok1 {
		return utils.Bool(ipnl.ID != nil && *ipnl.ID == state.ID), nil
	}
	if ok2 {
		return utils.Bool(cnl.ID != nil && *cnl.ID == state.ID), nil
	}
	return nil, fmt.Errorf("Unable to match object ID %q to a known type", state.ID)
}

func (NamedLocationResource) basicIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_named_location" "test" {
  display_name        = "acctestNLIP-%[1]d"
  ip {
    ip_ranges = [
      "1.1.1.1/32",
      "2.2.2.2/32",
    ]
  }
}
`, data.RandomInteger)
}

func (NamedLocationResource) completeIP(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_named_location" "test" {
  display_name        = "acctestNLIP-%[1]d"
  ip {
    ip_ranges = [
      "1.1.1.1/32",
      "2.2.2.2/32",
      "3.3.3.3/32",
      "64:ff9b::/96",
    ]
	trusted = true
  }
}
`, data.RandomInteger)
}

func (NamedLocationResource) basicCountry(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_named_location" "test" {
  display_name        = "acctestNLC-%[1]d"
  country {
    countries_and_regions = [
      "GB",
      "US",
    ]
  }
}
`, data.RandomInteger)
}

func (NamedLocationResource) completeCountry(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_named_location" "test" {
  display_name        = "acctestNLC-%[1]d"
  country {
    countries_and_regions = [
      "GB",
      "US",
      "JP",
    ]
	include_unknown_countries_and_regions = true
  }
}
`, data.RandomInteger)
}
