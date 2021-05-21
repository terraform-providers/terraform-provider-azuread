package provider

import (
	"github.com/hashicorp/terraform-provider-azuread/internal/services/applications"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/domains"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/groups"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/invitations"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/serviceprincipals"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/users"
)

func SupportedServices() []ServiceRegistration {
	return []ServiceRegistration{
		applications.Registration{},
		domains.Registration{},
		groups.Registration{},
		invitations.Registration{},
		serviceprincipals.Registration{},
		users.Registration{},
	}
}
