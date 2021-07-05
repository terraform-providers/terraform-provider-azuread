package domains

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
)

func domainsDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: domainsDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"admin_managed": {
				Description: "Set to `true` to only return domains whose DNS is managed by Microsoft 365",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			"include_unverified": {
				Description:   "Set to `true` if unverified Azure AD domains should be included",
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"only_default", "only_initial"}, // default or initial domains have to be verified
			},

			"only_default": {
				Description:   "Set to `true` to only return the default domain",
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"only_initial", "only_root"},
			},

			"only_initial": {
				Description:   "Set to `true` to only return the initial domain, which is your primary Azure Active Directory tenant domain",
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"only_default", "only_root"},
			},

			"only_root": {
				Description:   "Set to `true` to only return verified root domains. Excludes subdomains and unverified domains",
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"only_default", "only_initial"},
			},

			"supports_services": {
				Description: "A list of supported services that must be supported by a domain",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"domains": {
				Description: "A list of tenant domains",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Description: "The name of the domain",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"authentication_type": {
							Description: "The authentication type of the domain. Possible values include `Managed` or `Federated`",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"admin_managed": {
							Description: "Whether the DNS for the domain is managed by Microsoft 365",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"default": {
							Description: "Whether this is the default domain that is used for user creation",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"initial": {
							Description: "Whether this is the initial domain created by Azure Active Directory",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"root": {
							Description: "Whether the domain is a verified root domain (not a subdomain)",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"verified": {
							Description: "Whether the domain has completed domain ownership verification",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"supported_services": {
							Description: "A list of capabilities / services supported by the domain. Possible values include `Email`, `Sharepoint`, `EmailInternalRelayOnly`, `OfficeCommunicationsOnline`, `SharePointDefaultDomain`, `FullRedelegation`, `SharePointPublic`, `OrgIdAuthentication`, `Yammer` and `Intune`",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func domainsDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Domains.DomainsClient

	result, _, err := client.List(ctx)
	if err != nil {
		return tf.ErrorDiagF(err, "Could not list domains")
	}

	adminManaged := d.Get("admin_managed").(bool)
	onlyDefault := d.Get("only_default").(bool)
	onlyInitial := d.Get("only_initial").(bool)
	onlyRoot := d.Get("only_root").(bool)
	includeUnverified := d.Get("include_unverified").(bool)
	supportsServices := d.Get("supports_services").([]interface{})

	var domains []interface{}
	var domainNames []string
	if result != nil {
		for _, v := range *result {
			if adminManaged && v.IsAdminManaged != nil && !*v.IsAdminManaged {
				continue
			}
			if onlyDefault && v.IsDefault != nil && !*v.IsDefault {
				continue
			}
			if onlyInitial && v.IsInitial != nil && !*v.IsInitial {
				continue
			}
			if onlyRoot && v.IsRoot != nil && !*v.IsRoot {
				continue
			}
			if !includeUnverified && v.IsVerified != nil && !*v.IsVerified {
				continue
			}
			if len(supportsServices) > 0 && v.SupportedServices != nil {
				supported := 0
				for _, serviceNeeded := range supportsServices {
					for _, serviceSupported := range *v.SupportedServices {
						if serviceNeeded.(string) == serviceSupported {
							supported++
							break
						}
					}
				}
				if supported < len(supportsServices) {
					continue
				}
			}

			if v.ID != nil {
				domainNames = append(domainNames, *v.ID)

				domains = append(domains, map[string]interface{}{
					"admin_managed":       v.IsAdminManaged,
					"authentication_type": v.AuthenticationType,
					"default":             v.IsDefault,
					"domain_name":         v.ID,
					"initial":             v.IsInitial,
					"root":                v.IsRoot,
					"supported_services":  v.SupportedServices,
					"verified":            v.IsVerified,
				})
			}
		}
	}

	if len(domains) == 0 {
		return tf.ErrorDiagF(err, "No domains found for the provided filters")
	}

	// Generate a unique ID based on result
	h := sha1.New()
	if _, err := h.Write([]byte(strings.Join(domainNames, "/"))); err != nil {
		return tf.ErrorDiagF(err, "Unable to compute hash for domain names")
	}

	d.SetId(fmt.Sprintf("domains#%s#%s", client.BaseClient.TenantId, base64.URLEncoding.EncodeToString(h.Sum(nil))))
	tf.Set(d, "domains", domains)

	return nil
}
