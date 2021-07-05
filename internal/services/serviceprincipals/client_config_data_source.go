package serviceprincipals

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
)

func clientConfigDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: clientConfigDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"client_id": {
				Description: "The client ID (application ID) linked to the authenticated principal, or the application used for delegated authentication",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"tenant_id": {
				Description: "The tenant ID of the authenticated principal",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"object_id": {
				Description: "The object ID of the authenticated principal",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func clientConfigDataSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client)
	d.SetId(fmt.Sprintf("%s-%s-%s", client.TenantID, client.ClientID, client.Claims.ObjectId))
	tf.Set(d, "tenant_id", client.TenantID)
	tf.Set(d, "client_id", client.ClientID)
	tf.Set(d, "object_id", client.Claims.ObjectId)
	return nil
}
