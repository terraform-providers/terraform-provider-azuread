package client

import (
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/common"
)

type Client struct {
	GroupsClient *msgraph.GroupsClient
}

func NewClient(o *common.ClientOptions) *Client {
	msClient := msgraph.NewGroupsClient(o.TenantID)
	o.ConfigureClient(&msClient.BaseClient)

	return &Client{
		GroupsClient: msClient,
	}
}
