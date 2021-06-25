package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	MsClient *msgraph.ConditionalAccessPolicyClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewConditionalAccessPolicyClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient, nil)

	return &Client{
		MsClient: msClient,
	}
}
