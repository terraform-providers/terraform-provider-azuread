package client

import (
	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	AadClient *graphrbac.UsersClient // We set this as a UsersClient for now because configureClient does not support passing nil as autorest client.
	MsClient  *msgraph.InvitationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	aadClient := graphrbac.NewUsersClient(o.TenantID)
	msClient := msgraph.NewInvitationsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient, &aadClient.Client)

	return &Client{
		AadClient: &aadClient,
		MsClient:  msClient,
	}
}
