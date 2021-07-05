package clients

import (
	"context"
	"fmt"

	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
	applications "github.com/hashicorp/terraform-provider-azuread/internal/services/applications/client"
	domains "github.com/hashicorp/terraform-provider-azuread/internal/services/domains/client"
	groups "github.com/hashicorp/terraform-provider-azuread/internal/services/groups/client"
	invitations "github.com/hashicorp/terraform-provider-azuread/internal/services/invitations/client"
	serviceprincipals "github.com/hashicorp/terraform-provider-azuread/internal/services/serviceprincipals/client"
	users "github.com/hashicorp/terraform-provider-azuread/internal/services/users/client"
)

// Client contains the handles to all the specific Azure AD resource classes' respective clients
type Client struct {
	Environment environments.Environment
	TenantID    string
	ClientID    string
	Claims      auth.Claims

	TerraformVersion string

	StopContext context.Context

	Applications      *applications.Client
	Domains           *domains.Client
	Groups            *groups.Client
	Invitations       *invitations.Client
	ServicePrincipals *serviceprincipals.Client
	Users             *users.Client
}

func (client *Client) build(ctx context.Context, o *common.ClientOptions) error {
	client.StopContext = ctx

	client.Applications = applications.NewClient(o)
	client.Domains = domains.NewClient(o)
	client.Groups = groups.NewClient(o)
	client.Invitations = invitations.NewClient(o)
	client.ServicePrincipals = serviceprincipals.NewClient(o)
	client.Users = users.NewClient(o)

	// Acquire an access token upfront so we can decode and populate the JWT claims
	token, err := o.Authorizer.Token()
	if err != nil {
		return fmt.Errorf("unable to obtain access token: %v", err)
	}
	client.Claims, err = auth.ParseClaims(token)
	if err != nil {
		return fmt.Errorf("unable to parse claims in access token: %v", err)
	}
	if client.Claims.ObjectId == "" {
		return fmt.Errorf("parsing claims in access token: oid claim is empty")
	}

	return nil
}
