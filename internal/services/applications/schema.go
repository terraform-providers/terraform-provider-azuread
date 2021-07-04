package applications

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func schemaOptionalClaims() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "The name of the optional claim",
					Type:        schema.TypeString,
					Required:    true,
				},

				"source": {
					Description: "The source of the claim. If `source` is absent, the claim is a predefined optional claim. If `source` is `user`, the value of `name` is the extension property from the user object",
					Type:        schema.TypeString,
					Optional:    true,
					ValidateFunc: validation.StringInSlice(
						[]string{"user"},
						false,
					),
				},

				"essential": {
					Description: "Whether the claim specified by the client is necessary to ensure a smooth authorization experience",
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
				},

				"additional_properties": {
					Description: "List of additional properties of the claim. If a property exists in this list, it modifies the behaviour of the optional claim",
					Type:        schema.TypeList,
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
						ValidateFunc: validation.StringInSlice(
							[]string{
								"dns_domain_and_sam_account_name",
								"emit_as_roles",
								"include_externally_authenticated_upn",
								"include_externally_authenticated_upn_without_hash",
								"netbios_domain_and_sam_account_name",
								"on_premise_security_identifier",
								"sam_account_name",
								"use_guid",
							},
							false,
						),
					},
				},
			},
		},
	}
}
