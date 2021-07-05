package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	ServicePrincipalsClient *msgraph.ServicePrincipalsClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewServicePrincipalsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient)

	return &Client{
		ServicePrincipalsClient: msClient,
	}
}
