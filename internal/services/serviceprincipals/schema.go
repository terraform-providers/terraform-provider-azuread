package serviceprincipals

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func schemaAppRolesComputed() *schema.Schema {
	return &schema.Schema{
		Description: "",
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
	}
}

func schemaOauth2PermissionScopesComputed() *schema.Schema {
	return &schema.Schema{
		Description: "",
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
	}
}
