package applications

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func applicationDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: applicationDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Description:      "The application's object ID",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"application_id", "display_name", "object_id"},
				ValidateDiagFunc: validate.UUID,
			},

			"application_id": {
				Description:      "The Application ID (also called Client ID)",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"application_id", "display_name", "object_id"},
				ValidateDiagFunc: validate.UUID,
			},

			"disabled_by_microsoft": {
				Description: "Whether Microsoft has disabled the registered application",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"display_name": {
				Description:      "The display name for the application",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"application_id", "display_name", "object_id"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"api": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"known_client_applications": {
							Description: "Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"mapped_claims_enabled": {
							Description: "Allows an application to use claims mapping without specifying a custom signing key",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"oauth2_permission_scopes": {
							Description: "List of OAuth2 permission scopes published by the application",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The unique identifier of the delegated permission. Must be a valid UUID",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"admin_consent_description": {
										Description: "Delegated permission description that appears in all tenant-wide admin consent experiences, intended to be read by an administrator granting the permission on behalf of all users",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"admin_consent_display_name": {
										Description: "Display name for the delegated permission, intended to be read by an administrator granting the permission on behalf of all users",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"enabled": {
										Description: "Determines if the permission scope is enabled",
										Type:        schema.TypeBool,
										Computed:    true,
									},

									"type": {
										Description: "Whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator should be required for consent to the permissions. Possible values are `User` or `Admin`",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"user_consent_description": {
										Description: "Delegated permission description that appears in the end user consent experience, intended to be read by a user consenting on their own behalf",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"user_consent_display_name": {
										Description: "Display name for the delegated permission that appears in the end user consent experience",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"value": {
										Description: "The value that is used for the `scp` claim in OAuth 2.0 access tokens",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},

						"requested_access_token_version": {
							Description: "Specifies the access token version expected by this resource",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},

			"app_roles": {
				Description: "List of app roles published by the application",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The unique identifier of the app role",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"allowed_member_types": {
							Description: "Specifies whether this app role definition can be assigned to users and groups, or to other applications (that are accessing this application in a standalone scenario). Possible values are `User` or `Application`, or both",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"description": {
							Description: "Description of the app role that appears when the role is being assigned and, if the role functions as an application permissions, during the consent experiences",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"display_name": {
							Description: "Display name for the app role that appears during app role assignment and in consent experiences",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"enabled": {
							Description: "The unique identifier of the app role",
							Type:        schema.TypeBool,
							Computed:    true,
						},

						"value": {
							Description: "The value that is used for the `roles` claim in ID tokens and OAuth 2.0 access tokens that are authenticating an assigned service or user principal",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"app_role_ids": {
				Description: "Mapping of app role names to UUIDs",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"device_only_auth_enabled": {
				Description: "Specifies whether this application supports device authentication without a user.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"fallback_public_client_enabled": {
				Description: "The fallback application type as public client, such as an installed application running on a mobile device",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"group_membership_claims": {
				Description: "The `groups` claim issued in a user or OAuth 2.0 access token that the app expects",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"identifier_uris": {
				Description: "A list of user-defined URI(s) that uniquely identify a Web application within it's Azure AD tenant, or within a verified custom domain if the application is multi-tenant",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"logo_url": {
				Description: "CDN URL to the application's logo",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"marketing_url": {
				Description: "URL of the application's marketing page",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"oauth2_permission_scope_ids": {
				Description: "Mapping of OAuth2.0 permission scope names to UUIDs",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"oauth2_post_response_required": {
				Description: "Specifies whether, as part of OAuth 2.0 token requests, Azure AD allows POST requests, as opposed to GET requests.",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"optional_claims": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_token": schemaOptionalClaims(),
						"id_token":     schemaOptionalClaims(),
						"saml2_token":  schemaOptionalClaims(),
					},
				},
			},

			"owners": {
				Description: "A list of object IDs of principals that are assigned ownership of the application",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"privacy_statement_url": {
				Description: "URL of the application's privacy statement",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"public_client": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_uris": {
							Description: "The URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"publisher_domain": {
				Description: "The verified publisher domain for the application",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"required_resource_access": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_app_id": {
							Description: "The unique identifier for the resource that the application requires access to. This is the Application ID of the target application",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"resource_access": {
							Description: "A collection of `resource_access` blocks describing OAuth2.0 permission scopes and app roles that the application requires from the specified resource",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "The unique identifier for an app role or OAuth2 permission scope published by the resource application",
										Type:        schema.TypeString,
										Computed:    true,
									},

									"type": {
										Description: "Specifies whether the `id` property references an app role or an OAuth2 permission scope. Possible values are `Role` or `Scope`",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},

			"sign_in_audience": {
				Description: "The Microsoft account types that are supported for the current application",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"single_page_application": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redirect_uris": {
							Description: "The URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"support_url": {
				Description: "URL of the application's support page",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"terms_of_service_url": {
				Description: "URL of the application's terms of service statement",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"web": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"homepage_url": {
							Description: "Home page or landing page of the application",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"logout_url": {
							Description: "The URL that will be used by Microsoft's authorization service to sign out a user using front-channel, back-channel or SAML logout protocols",
							Type:        schema.TypeString,
							Computed:    true,
						},

						"redirect_uris": {
							Description: "A list of URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"implicit_grant": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_token_issuance_enabled": {
										Description: "Whether this web application can request an access token using OAuth 2.0 implicit flow",
										Type:        schema.TypeBool,
										Computed:    true,
									},

									"id_token_issuance_enabled": {
										Description: "Whether this web application can request an ID token using OAuth 2.0 implicit flow",
										Type:        schema.TypeBool,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func applicationDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient

	var app *msgraph.Application

	if objectId, ok := d.Get("object_id").(string); ok && objectId != "" {
		var status int
		var err error
		app, status, err = client.Get(ctx, objectId)
		if err != nil {
			if status == http.StatusNotFound {
				return tf.ErrorDiagPathF(nil, "object_id", "Application with object ID %q was not found", objectId)
			}

			return tf.ErrorDiagPathF(err, "object_id", "Retrieving Application with object ID %q", objectId)
		}
	} else {
		var fieldName, fieldValue string
		if applicationId, ok := d.Get("application_id").(string); ok && applicationId != "" {
			fieldName = "appId"
			fieldValue = applicationId
		} else if displayName, ok := d.Get("display_name").(string); ok && displayName != "" {
			fieldName = "displayName"
			fieldValue = displayName
		} else {
			return tf.ErrorDiagF(nil, "One of `object_id`, `application_id` or `displayName` must be specified")
		}

		filter := fmt.Sprintf("%s eq '%s'", fieldName, fieldValue)

		result, _, err := client.List(ctx, odata.Query{Filter: filter})
		if err != nil {
			return tf.ErrorDiagF(err, "Listing applications for filter %q", filter)
		}

		switch {
		case result == nil || len(*result) == 0:
			return tf.ErrorDiagF(fmt.Errorf("No applications found matching filter: %q", filter), "Application not found")
		case len(*result) > 1:
			return tf.ErrorDiagF(fmt.Errorf("Found multiple applications matching filter: %q", filter), "Multiple applications found")
		}

		app = &(*result)[0]
		switch fieldName {
		case "appId":
			if app.AppId == nil {
				return tf.ErrorDiagF(fmt.Errorf("nil AppID for applications matching filter: %q", filter), "Bad API Response")
			}
			if *app.AppId != fieldValue {
				return tf.ErrorDiagF(fmt.Errorf("AppID does not match (%q != %q) for applications matching filter: %q", *app.AppId, fieldValue, filter), "Bad API Response")
			}
		case "displayName":
			if app.DisplayName == nil {
				return tf.ErrorDiagF(fmt.Errorf("nil displayName for applications matching filter: %q", filter), "Bad API Response")
			}
			if *app.DisplayName != fieldValue {
				return tf.ErrorDiagF(fmt.Errorf("DisplayName does not match (%q != %q) for applications matching filter: %q", *app.DisplayName, fieldValue, filter), "Bad API Response")
			}
		}
	}

	if app == nil {
		return tf.ErrorDiagF(fmt.Errorf("app was unexpectedly nil"), "Application not found")
	}

	if app.ID == nil {
		return tf.ErrorDiagF(fmt.Errorf("Object ID returned for application is nil"), "Bad API Response")
	}

	d.SetId(*app.ID)

	tf.Set(d, "api", flattenApplicationApi(app.Api, true))
	tf.Set(d, "app_roles", flattenApplicationAppRoles(app.AppRoles))
	tf.Set(d, "app_role_ids", flattenApplicationAppRoleIDs(app.AppRoles))
	tf.Set(d, "application_id", app.AppId)
	tf.Set(d, "device_only_auth_enabled", app.IsDeviceOnlyAuthSupported)
	tf.Set(d, "disabled_by_microsoft", fmt.Sprintf("%v", app.DisabledByMicrosoftStatus))
	tf.Set(d, "display_name", app.DisplayName)
	tf.Set(d, "fallback_public_client_enabled", app.IsFallbackPublicClient)
	tf.Set(d, "group_membership_claims", flattenApplicationGroupMembershipClaims(app.GroupMembershipClaims))
	tf.Set(d, "identifier_uris", tf.FlattenStringSlicePtr(app.IdentifierUris))
	tf.Set(d, "oauth2_post_response_required", app.Oauth2RequirePostResponse)
	tf.Set(d, "object_id", app.ID)
	tf.Set(d, "optional_claims", flattenApplicationOptionalClaims(app.OptionalClaims))
	tf.Set(d, "public_client", flattenApplicationPublicClient(app.PublicClient))
	tf.Set(d, "publisher_domain", app.PublisherDomain)
	tf.Set(d, "required_resource_access", flattenApplicationRequiredResourceAccess(app.RequiredResourceAccess))
	tf.Set(d, "sign_in_audience", app.SignInAudience)
	tf.Set(d, "single_page_application", flattenApplicationSpa(app.Spa))
	tf.Set(d, "web", flattenApplicationWeb(app.Web))

	if app.Api != nil {
		tf.Set(d, "oauth2_permission_scope_ids", flattenApplicationOAuth2PermissionScopeIDs(app.Api.OAuth2PermissionScopes))
	}

	if app.Info != nil {
		tf.Set(d, "logo_url", app.Info.LogoUrl)
		tf.Set(d, "marketing_url", app.Info.MarketingUrl)
		tf.Set(d, "privacy_statement_url", app.Info.PrivacyStatementUrl)
		tf.Set(d, "support_url", app.Info.SupportUrl)
		tf.Set(d, "terms_of_service_url", app.Info.TermsOfServiceUrl)
	}

	owners, _, err := client.ListOwners(ctx, *app.ID)
	if err != nil {
		return tf.ErrorDiagPathF(err, "owners", "Could not retrieve owners for application with object ID %q", *app.ID)
	}
	tf.Set(d, "owners", owners)

	return nil
}
