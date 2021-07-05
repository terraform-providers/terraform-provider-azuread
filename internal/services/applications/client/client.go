package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	ApplicationsClient *msgraph.ApplicationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewApplicationsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient)

	return &Client{
		ApplicationsClient: msClient,
	}
}
