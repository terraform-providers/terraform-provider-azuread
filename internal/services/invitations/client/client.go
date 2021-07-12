package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	MsClient *msgraph.InvitationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewInvitationsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient)

	return &Client{
		MsClient: msClient,
	}
}
