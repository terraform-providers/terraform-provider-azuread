package conditionalaccesspolicies_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azuread/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

type ConditionalAccessPolicyResource struct{}

func TestAccConditionalAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("display_name").HasValue(fmt.Sprintf("acctest-CONPOLICY-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("disabled"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccConditionalAccessPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuread_conditional_access_policy", "test")
	r := ConditionalAccessPolicyResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ConditionalAccessPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	var id *string

	app, status, err := clients.ConditionalAccessPolicies.MsClient.Get(ctx, state.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return nil, fmt.Errorf("Conditional Access Policy with ID %q does not exist", state.ID)
		}
		return nil, fmt.Errorf("failed to retrieve Conditional Access Policy with ID %q: %+v", state.ID, err)
	}
	id = app.ID

	return utils.Bool(id != nil && *id == state.ID), nil
}

func (ConditionalAccessPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    applications {
      included_applications = ["All"]
    }
    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
    client_app_types = ["browser"]
    locations {
      included_locations = ["All"]
    }
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["block"]
  }
}
`, data.RandomInteger)
}

func (ConditionalAccessPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuread_conditional_access_policy" "test" {
  display_name = "acctest-CONPOLICY-%[1]d"
  state        = "disabled"

  conditions {
    applications {
      included_applications = ["All"]
      excluded_applications = ["00000004-0000-0ff1-ce00-000000000000"]
    }
    users {
      included_users = ["All"]
      excluded_users = ["GuestsOrExternalUsers"]
    }
    client_app_types = ["all"]
    locations {
      included_locations = ["All"]
      excluded_locations = ["AllTrusted"]
    }
    platforms {
      included_platforms = ["android"]
      excluded_platforms = ["iOS"]
    }
    sign_in_risk_levels = ["medium"]
    user_risk_levels    = ["medium"]
  }

  grant_controls {
    operator          = "OR"
    built_in_controls = ["mfa"]
  }

  session_controls {
    application_enforced_restrictions {
      enabled = true
    }
    cloud_app_security {
      enabled                 = true
      cloud_app_security_type = "monitorOnly"
    }
    sign_in_frequency {
      enabled = true
      type    = "hours"
      value   = 10
    }
  }
}
`, data.RandomInteger)
}
