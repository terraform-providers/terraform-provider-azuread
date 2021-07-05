package helpers

import (
	"github.com/manicminer/hamilton/msgraph"
)

func ApplicationFlattenAppRoles(in *[]msgraph.AppRole) (result []map[string]interface{}) {
	if in == nil {
		return
	}

	for _, role := range *in {
		roleId := ""
		if role.ID != nil {
			roleId = *role.ID
		}
		allowedMemberTypes := make([]interface{}, 0)
		if v := role.AllowedMemberTypes; v != nil {
			for _, m := range *v {
				allowedMemberTypes = append(allowedMemberTypes, m)
			}
		}
		description := ""
		if role.Description != nil {
			description = *role.Description
		}
		displayName := ""
		if role.DisplayName != nil {
			displayName = *role.DisplayName
		}
		enabled := false
		if role.IsEnabled != nil && *role.IsEnabled {
			enabled = true
		}
		value := ""
		if role.Value != nil {
			value = *role.Value
		}
		result = append(result, map[string]interface{}{
			"id":                   roleId,
			"allowed_member_types": allowedMemberTypes,
			"description":          description,
			"display_name":         displayName,
			"enabled":              enabled,
			"value":                value,
		})
	}

	return //nolint:nakedret
}

func ApplicationFlattenOAuth2PermissionScopes(in *[]msgraph.PermissionScope) (result []map[string]interface{}) {
	if in == nil {
		return
	}

	for _, p := range *in {
		adminConsentDescription := ""
		if v := p.AdminConsentDescription; v != nil {
			adminConsentDescription = *v
		}

		adminConsentDisplayName := ""
		if v := p.AdminConsentDisplayName; v != nil {
			adminConsentDisplayName = *v
		}

		id := ""
		if v := p.ID; v != nil {
			id = *v
		}

		userConsentDescription := ""
		if v := p.UserConsentDescription; v != nil {
			userConsentDescription = *v
		}

		userConsentDisplayName := ""
		if v := p.UserConsentDisplayName; v != nil {
			userConsentDisplayName = *v
		}

		value := ""
		if v := p.Value; v != nil {
			value = *v
		}

		enabled := p.IsEnabled != nil && *p.IsEnabled

		result = append(result, map[string]interface{}{
			"admin_consent_description":  adminConsentDescription,
			"admin_consent_display_name": adminConsentDisplayName,
			"id":                         id,
			"enabled":                    enabled,
			"type":                       p.Type,
			"user_consent_description":   userConsentDescription,
			"user_consent_display_name":  userConsentDisplayName,
			"value":                      value,
		})
	}

	return //nolint:nakedret
}
