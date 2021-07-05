package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	UsersClient *msgraph.UsersClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewUsersClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient)

	return &Client{
		UsersClient: msClient,
	}
}
