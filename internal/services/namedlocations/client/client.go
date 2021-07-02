package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	MsClient *msgraph.NamedLocationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewNamedLocationsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient)

	return &Client{
		MsClient: msClient,
	}
}
