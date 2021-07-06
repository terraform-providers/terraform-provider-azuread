package applications

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/terraform-provider-azuread/internal/helpers"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

func applicationDisableAppRoles(ctx context.Context, client *msgraph.ApplicationsClient, application *msgraph.Application, newRoles *[]msgraph.AppRole) error {
	if application.ID == nil {
		return fmt.Errorf("cannot use Application model with nil ID")
	}

	if newRoles == nil {
		newRoles = &[]msgraph.AppRole{}
	}

	app, status, err := client.Get(ctx, *application.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return fmt.Errorf("application with ID %q was not found", *application.ID)
		}

		return fmt.Errorf("retrieving Application with object ID %q: %+v", *application.ID, err)
	}

	var existingRoles []msgraph.AppRole
	if app.AppRoles != nil {
		existingRoles = *app.AppRoles
	}

	// Don't update if no changes to be made
	if reflect.DeepEqual(existingRoles, *newRoles) {
		return nil
	}

	// Identify any roles to be changed
	var disable bool
	for _, new := range *newRoles {
		if new.ID == nil || *new.ID == "" {
			return fmt.Errorf("new role provided with nil or empty ID")
		}
		for i, existing := range existingRoles {
			if existing.ID != nil && *existing.ID == *new.ID {
				if existing.IsEnabled != nil && *existing.IsEnabled && !reflect.DeepEqual(existing, new) {
					*existingRoles[i].IsEnabled = false
					disable = true
				}
				break
			}
		}
	}

	// Identify any roles to be removed
	for i, existing := range existingRoles {
		found := false
		for _, new := range *newRoles {
			if existing.ID != nil && *new.ID == *existing.ID {
				found = true
				break
			}
		}
		if !found {
			*existingRoles[i].IsEnabled = false
			disable = true
		}
	}

	if disable {
		// Disable any changed or removed roles
		properties := msgraph.Application{
			ID:       application.ID,
			AppRoles: &existingRoles,
		}
		if _, err := client.Update(ctx, properties); err != nil {
			return fmt.Errorf("disabling App Roles for Application with object ID %q: %+v", *application.ID, err)
		}

		// Wait for application manifest to reflect the disabled roles
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("context has no deadline")
		}
		timeout := time.Until(deadline)
		_, err = (&resource.StateChangeConf{
			Pending:    []string{"Waiting"},
			Target:     []string{"Disabled"},
			Timeout:    timeout,
			MinTimeout: 1 * time.Second,
			Refresh: func() (interface{}, string, error) {
				app, _, err := client.Get(ctx, *application.ID)
				if err != nil {
					return nil, "Error", fmt.Errorf("retrieving Application with object ID %q: %+v", *application.ID, err)
				}
				if app == nil || app.AppRoles == nil {
					return nil, "Error", fmt.Errorf("reading roles for Application with object ID %q: %+v", *application.ID, err)
				}
				actualRoles := *app.AppRoles
				for _, expectedRole := range existingRoles {
					if expectedRole.IsEnabled != nil && !*expectedRole.IsEnabled {
						for _, actualRole := range actualRoles {
							if expectedRole.ID != nil && actualRole.ID != nil && *expectedRole.ID == *actualRole.ID {
								if actualRole.IsEnabled != nil && *actualRole.IsEnabled {
									return actualRoles, "Waiting", nil
								}
								break
							}
						}
					}
				}
				return actualRoles, "Disabled", nil
			},
		}).WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf("waiting for App Roles to be disabled for Application with object ID %q: %+v", *application.ID, err)
		}
	}

	return nil
}

func applicationDisableOauth2PermissionScopes(ctx context.Context, client *msgraph.ApplicationsClient, application *msgraph.Application, newScopes *[]msgraph.PermissionScope) error {
	if application.ID == nil {
		return fmt.Errorf("Cannot use Application model with nil ID")
	}

	if newScopes == nil {
		newScopes = &[]msgraph.PermissionScope{}
	}

	app, status, err := client.Get(ctx, *application.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return fmt.Errorf("application with ID %q was not found", *application.ID)
		}

		return fmt.Errorf("retrieving Application with object ID %q: %+v", *application.ID, err)
	}

	var existingScopes []msgraph.PermissionScope
	if app.Api != nil && app.Api.OAuth2PermissionScopes != nil {
		existingScopes = *app.Api.OAuth2PermissionScopes
	}

	// Don't update if no changes to be made
	if reflect.DeepEqual(existingScopes, *newScopes) {
		return nil
	}

	// Identify any scopes to be changed
	var disable bool
	for _, new := range *newScopes {
		if new.ID == nil || *new.ID == "" {
			return fmt.Errorf("new scope provided with nil or empty ID")
		}
		for i, existing := range existingScopes {
			if existing.ID != nil && *existing.ID == *new.ID {
				if existing.IsEnabled != nil && *existing.IsEnabled && !reflect.DeepEqual(existing, new) {
					*existingScopes[i].IsEnabled = false
					disable = true
				}
				break
			}
		}
	}

	// Identify any scopes to be removed
	for i, existing := range existingScopes {
		found := false
		for _, new := range *newScopes {
			if existing.ID != nil && *new.ID == *existing.ID {
				found = true
				break
			}
		}
		if !found {
			*existingScopes[i].IsEnabled = false
			disable = true
		}
	}

	if disable {
		// Disable any changed or removed scopes
		properties := msgraph.Application{
			ID: application.ID,
			Api: &msgraph.ApplicationApi{
				OAuth2PermissionScopes: &existingScopes,
			},
		}
		if _, err := client.Update(ctx, properties); err != nil {
			return fmt.Errorf("disabling OAuth2 Permission Scopes for Application with object ID %q: %+v", *application.ID, err)
		}

		// Wait for application manifest to reflect the disabled scopes
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("context has no deadline")
		}
		timeout := time.Until(deadline)
		_, err = (&resource.StateChangeConf{
			Pending:    []string{"Waiting"},
			Target:     []string{"Disabled"},
			Timeout:    timeout,
			MinTimeout: 1 * time.Second,
			Refresh: func() (interface{}, string, error) {
				app, _, err := client.Get(ctx, *application.ID)
				if err != nil {
					return nil, "Error", fmt.Errorf("retrieving Application with object ID %q: %+v", *application.ID, err)
				}
				if app == nil || app.Api == nil || app.Api.OAuth2PermissionScopes == nil {
					return nil, "Error", fmt.Errorf("reading scopes for Application with object ID %q: %+v", *application.ID, err)
				}
				actualScopes := *app.Api.OAuth2PermissionScopes
				for _, expectedScope := range existingScopes {
					if expectedScope.IsEnabled != nil && !*expectedScope.IsEnabled {
						for _, actualScope := range actualScopes {
							if expectedScope.ID != nil && actualScope.ID != nil && *expectedScope.ID == *actualScope.ID {
								if actualScope.IsEnabled != nil && *actualScope.IsEnabled {
									return actualScopes, "Waiting", nil
								}
								break
							}
						}
					}
				}
				return actualScopes, "Disabled", nil
			},
		}).WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf("waiting for OAuth2 Permission Scopes to be disabled for Application with object ID %q: %+v", *application.ID, err)
		}
	}

	return nil
}

func ApplicationFindAppRole(app *msgraph.Application, roleId string) (*msgraph.AppRole, error) {
	if app == nil || app.AppRoles == nil {
		return nil, nil
	}
	if roleId == "" {
		return nil, fmt.Errorf("specified role ID is empty")
	}
	for _, r := range *app.AppRoles {
		if r.ID == nil {
			continue
		}
		if *r.ID == roleId {
			return &r, nil
		}
	}
	return nil, nil
}

func ApplicationFindOAuth2PermissionScope(app *msgraph.Application, scopeId string) (*msgraph.PermissionScope, error) {
	if app == nil || app.Api == nil || app.Api.OAuth2PermissionScopes == nil {
		return nil, nil
	}
	if scopeId == "" {
		return nil, fmt.Errorf("specified scope ID is empty")
	}
	for _, s := range *app.Api.OAuth2PermissionScopes {
		if s.ID == nil {
			continue
		}
		if *s.ID == scopeId {
			return &s, nil
		}
	}
	return nil, nil
}

func applicationFindByName(ctx context.Context, client *msgraph.ApplicationsClient, displayName string) (*[]msgraph.Application, error) {
	query := odata.Query{
		Filter: fmt.Sprintf("displayName eq '%s'", displayName),
	}
	apps, _, err := client.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to list Applications with filter %q: %+v", query.Filter, err)
	}

	result := make([]msgraph.Application, 0)
	if apps != nil {
		for _, app := range *apps {
			if app.DisplayName != nil && *app.DisplayName == displayName {
				result = append(result, app)
			}
		}
	}

	return &result, nil
}

func applicationSetOwners(ctx context.Context, client *msgraph.ApplicationsClient, application *msgraph.Application, desiredOwners []string) error {
	if application.ID == nil {
		return fmt.Errorf("Cannot use Application model with nil ID")
	}

	owners, _, err := client.ListOwners(ctx, *application.ID)
	if err != nil {
		return fmt.Errorf("retrieving owners for Application with object ID %q: %+v", *application.ID, err)
	}

	existingOwners := *owners
	ownersForRemoval := utils.Difference(existingOwners, desiredOwners)
	ownersToAdd := utils.Difference(desiredOwners, existingOwners)

	if ownersToAdd != nil {
		for _, m := range ownersToAdd {
			application.AppendOwner(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, m)
		}

		if _, err := client.AddOwners(ctx, application); err != nil {
			return fmt.Errorf("adding owners to Application with object ID %q: %+v", *application.ID, err)
		}
	}

	if ownersForRemoval != nil {
		if _, err = client.RemoveOwners(ctx, *application.ID, &ownersForRemoval); err != nil {
			return fmt.Errorf("removing owner from Application with object ID %q: %+v", *application.ID, err)
		}
	}

	return nil
}

func applicationValidateRolesScopes(appRoles, oauth2Permissions []interface{}) error {
	var ids, values []string

	for _, roleRaw := range appRoles {
		role := roleRaw.(map[string]interface{})
		if id := role["id"].(string); id != "" {
			ids = append(ids, id)
		}
		if val := role["value"].(string); val != "" {
			values = append(values, val)
		}
	}

	for _, scopeRaw := range oauth2Permissions {
		scope := scopeRaw.(map[string]interface{})
		if id := scope["id"].(string); id != "" {
			ids = append(ids, id)
		}
		if val := scope["value"].(string); val != "" {
			values = append(values, val)
		}
	}

	encounteredIds := make([]string, 0)
	for _, id := range ids {
		for _, en := range encounteredIds {
			if en == id {
				return fmt.Errorf("validation failed: duplicate ID found: %q", id)
			}
		}
		encounteredIds = append(encounteredIds, id)
	}

	encounteredValues := make([]string, 0)
	for _, val := range values {
		for _, en := range encounteredValues {
			if en == val {
				return fmt.Errorf("validation failed: duplicate value found: %q", val)
			}
		}
		encounteredValues = append(encounteredValues, val)
	}

	return nil
}

func expandApplicationApi(input []interface{}) (result *msgraph.ApplicationApi) {
	result = &msgraph.ApplicationApi{
		AcceptMappedClaims:          utils.Bool(false),
		KnownClientApplications:     &[]string{},
		OAuth2PermissionScopes:      &[]msgraph.PermissionScope{},
		RequestedAccessTokenVersion: utils.Int32(int32(1)),
	}

	if len(input) == 0 {
		return
	}

	in := input[0].(map[string]interface{})
	if v, ok := in["mapped_claims_enabled"]; ok {
		result.AcceptMappedClaims = utils.Bool(v.(bool))
	}
	if v, ok := in["known_client_applications"]; ok {
		result.KnownClientApplications = tf.ExpandStringSlicePtr(v.(*schema.Set).List())
	}
	result.OAuth2PermissionScopes = expandApplicationOAuth2PermissionScope(in["oauth2_permission_scope"].(*schema.Set).List())
	if v, ok := in["requested_access_token_version"]; ok {
		result.RequestedAccessTokenVersion = utils.Int32(int32(v.(int)))
	}

	return
}

func expandApplicationAppRoles(input []interface{}) *[]msgraph.AppRole {
	result := make([]msgraph.AppRole, 0)

	if len(input) == 0 {
		return &result
	}

	for _, appRoleRaw := range input {
		appRole := appRoleRaw.(map[string]interface{})

		var allowedMemberTypes []msgraph.AppRoleAllowedMemberType
		for _, allowedMemberType := range appRole["allowed_member_types"].(*schema.Set).List() {
			allowedMemberTypes = append(allowedMemberTypes, allowedMemberType.(string))
		}

		newAppRole := msgraph.AppRole{
			ID:                 utils.String(appRole["id"].(string)),
			AllowedMemberTypes: &allowedMemberTypes,
			Description:        utils.String(appRole["description"].(string)),
			DisplayName:        utils.String(appRole["display_name"].(string)),
			IsEnabled:          utils.Bool(appRole["enabled"].(bool)),
		}

		if v, ok := appRole["value"]; ok {
			newAppRole.Value = utils.String(v.(string))
		}

		result = append(result, newAppRole)
	}

	return &result
}

func expandApplicationGroupMembershipClaims(in []interface{}) *[]msgraph.GroupMembershipClaim {
	result := make([]msgraph.GroupMembershipClaim, 0)
	if len(in) == 0 {
		return &result
	}
	for _, claimRaw := range in {
		result = append(result, claimRaw.(string))
	}
	return &result
}

func expandApplicationImplicitGrantSettings(input []interface{}) *msgraph.ImplicitGrantSettings {
	var enableAccessTokenIssuance, enableIdTokenIssuance bool

	if input != nil || len(input) > 0 {
		in := input[0].(map[string]interface{})
		enableAccessTokenIssuance = in["access_token_issuance_enabled"].(bool)
		enableIdTokenIssuance = in["id_token_issuance_enabled"].(bool)
	}

	return &msgraph.ImplicitGrantSettings{
		EnableAccessTokenIssuance: utils.Bool(enableAccessTokenIssuance),
		EnableIdTokenIssuance:     utils.Bool(enableIdTokenIssuance),
	}
}

func expandApplicationOAuth2PermissionScope(in []interface{}) *[]msgraph.PermissionScope {
	result := make([]msgraph.PermissionScope, 0)

	for _, raw := range in {
		oauth2Permissions := raw.(map[string]interface{})

		result = append(result,
			msgraph.PermissionScope{
				AdminConsentDescription: utils.String(oauth2Permissions["admin_consent_description"].(string)),
				AdminConsentDisplayName: utils.String(oauth2Permissions["admin_consent_display_name"].(string)),
				ID:                      utils.String(oauth2Permissions["id"].(string)),
				IsEnabled:               utils.Bool(oauth2Permissions["enabled"].(bool)),
				Type:                    oauth2Permissions["type"].(string),
				UserConsentDescription:  utils.String(oauth2Permissions["user_consent_description"].(string)),
				UserConsentDisplayName:  utils.String(oauth2Permissions["user_consent_display_name"].(string)),
				Value:                   utils.String(oauth2Permissions["value"].(string)),
			},
		)
	}

	return &result
}

func expandApplicationOptionalClaims(in []interface{}) *msgraph.OptionalClaims {
	result := msgraph.OptionalClaims{}

	if len(in) == 0 || in[0] == nil {
		return &result
	}

	optionalClaims := in[0].(map[string]interface{})

	result.AccessToken = expandApplicationOptionalClaim(optionalClaims["access_token"].([]interface{}))
	result.IdToken = expandApplicationOptionalClaim(optionalClaims["id_token"].([]interface{}))
	result.Saml2Token = expandApplicationOptionalClaim(optionalClaims["saml2_token"].([]interface{}))

	return &result
}

func expandApplicationOptionalClaim(in []interface{}) *[]msgraph.OptionalClaim {
	result := make([]msgraph.OptionalClaim, 0)

	for _, optionalClaimRaw := range in {
		optionalClaim := optionalClaimRaw.(map[string]interface{})

		additionalProps := make([]string, 0)
		if props, ok := optionalClaim["additional_properties"]; ok && props != nil {
			for _, prop := range props.([]interface{}) {
				additionalProps = append(additionalProps, prop.(string))
			}
		}

		newClaim := msgraph.OptionalClaim{
			Name:                 utils.String(optionalClaim["name"].(string)),
			Essential:            utils.Bool(optionalClaim["essential"].(bool)),
			AdditionalProperties: &additionalProps,
		}

		if source, ok := optionalClaim["source"].(string); ok && source != "" {
			newClaim.Source = &source
		}

		result = append(result, newClaim)
	}

	return &result
}

func expandApplicationPublicClient(input []interface{}) (result *msgraph.PublicClient) {
	result = &msgraph.PublicClient{
		RedirectUris: &[]string{},
	}

	if len(input) == 0 {
		return
	}

	in := input[0].(map[string]interface{})
	result.RedirectUris = tf.ExpandStringSlicePtr(in["redirect_uris"].(*schema.Set).List())

	return
}

func expandApplicationRequiredResourceAccess(in []interface{}) *[]msgraph.RequiredResourceAccess {
	result := make([]msgraph.RequiredResourceAccess, 0)

	for _, raw := range in {
		requiredResourceAccess := raw.(map[string]interface{})

		result = append(result, msgraph.RequiredResourceAccess{
			ResourceAppId: utils.String(requiredResourceAccess["resource_app_id"].(string)),
			ResourceAccess: expandApplicationResourceAccess(
				requiredResourceAccess["resource_access"].([]interface{}),
			),
		})
	}

	return &result
}

func expandApplicationResourceAccess(in []interface{}) *[]msgraph.ResourceAccess {
	result := make([]msgraph.ResourceAccess, 0)

	for _, resourceAccessRaw := range in {
		resourceAccess := resourceAccessRaw.(map[string]interface{})

		result = append(result, msgraph.ResourceAccess{
			ID:   utils.String(resourceAccess["id"].(string)),
			Type: resourceAccess["type"].(string),
		})
	}

	return &result
}

func expandApplicationSpa(input []interface{}) (result *msgraph.ApplicationSpa) {
	result = &msgraph.ApplicationSpa{
		RedirectUris: &[]string{},
	}

	if len(input) == 0 {
		return
	}

	in := input[0].(map[string]interface{})
	result.RedirectUris = tf.ExpandStringSlicePtr(in["redirect_uris"].(*schema.Set).List())

	return
}

func expandApplicationWeb(input []interface{}) (result *msgraph.ApplicationWeb) {
	result = &msgraph.ApplicationWeb{
		HomePageUrl:           utils.NullableString(""),
		ImplicitGrantSettings: expandApplicationImplicitGrantSettings(nil),
		LogoutUrl:             utils.NullableString(""),
		RedirectUris:          &[]string{},
	}

	if len(input) == 0 {
		return
	}

	in := input[0].(map[string]interface{})
	result.HomePageUrl = utils.NullableString(in["homepage_url"].(string))
	result.LogoutUrl = utils.NullableString(in["logout_url"].(string))
	result.ImplicitGrantSettings = expandApplicationImplicitGrantSettings(in["implicit_grant"].([]interface{}))
	result.RedirectUris = tf.ExpandStringSlicePtr(in["redirect_uris"].(*schema.Set).List())

	return
}

func flattenApplicationApi(in *msgraph.ApplicationApi, dataSource bool) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	mappedClaims := false
	if in.AcceptMappedClaims != nil {
		mappedClaims = *in.AcceptMappedClaims
	}

	scopesKey := "oauth2_permission_scope"
	if dataSource {
		scopesKey = "oauth2_permission_scopes"
	}

	accessTokenVersion := 1
	if in.RequestedAccessTokenVersion != nil {
		accessTokenVersion = int(*in.RequestedAccessTokenVersion)
	}

	return []map[string]interface{}{{
		"known_client_applications":      tf.FlattenStringSlicePtr(in.KnownClientApplications),
		"mapped_claims_enabled":          mappedClaims,
		scopesKey:                        flattenApplicationOAuth2PermissionScopes(in.OAuth2PermissionScopes),
		"requested_access_token_version": accessTokenVersion,
	}}
}

func flattenApplicationAppRoleIDs(in *[]msgraph.AppRole) map[string]string {
	return helpers.ApplicationFlattenAppRoleIDs(in)
}

func flattenApplicationAppRoles(in *[]msgraph.AppRole) []map[string]interface{} {
	return helpers.ApplicationFlattenAppRoles(in)
}

func flattenApplicationImplicitGrant(in *msgraph.ImplicitGrantSettings) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	accessToken := false
	if in.EnableAccessTokenIssuance != nil {
		accessToken = *in.EnableAccessTokenIssuance
	}
	idToken := false
	if in.EnableIdTokenIssuance != nil {
		idToken = *in.EnableIdTokenIssuance
	}

	return []map[string]interface{}{{
		"access_token_issuance_enabled": accessToken,
		"id_token_issuance_enabled":     idToken,
	}}
}

func flattenApplicationOAuth2PermissionScopeIDs(in *[]msgraph.PermissionScope) map[string]string {
	return helpers.ApplicationFlattenOAuth2PermissionScopeIDs(in)
}

func flattenApplicationOAuth2PermissionScopes(in *[]msgraph.PermissionScope) []map[string]interface{} {
	return helpers.ApplicationFlattenOAuth2PermissionScopes(in)
}

func flattenApplicationOptionalClaims(in *msgraph.OptionalClaims) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{{
		"access_token": flattenApplicationOptionalClaim(in.AccessToken),
		"id_token":     flattenApplicationOptionalClaim(in.IdToken),
		"saml2_token":  flattenApplicationOptionalClaim(in.Saml2Token),
	}}
}

func flattenApplicationOptionalClaim(in *[]msgraph.OptionalClaim) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	optionalClaims := make([]interface{}, 0)
	for _, claim := range *in {
		optionalClaim := map[string]interface{}{
			"name":                  claim.Name,
			"essential":             claim.Essential,
			"source":                "",
			"additional_properties": []string{},
		}

		if claim.Source != nil {
			optionalClaim["source"] = *claim.Source
		}

		if claim.AdditionalProperties != nil && len(*claim.AdditionalProperties) > 0 {
			optionalClaim["additional_properties"] = *claim.AdditionalProperties
		}

		optionalClaims = append(optionalClaims, optionalClaim)
	}

	return optionalClaims
}

func flattenApplicationPublicClient(in *msgraph.PublicClient) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{{
		"redirect_uris": tf.FlattenStringSlicePtr(in.RedirectUris),
	}}
}

func flattenApplicationRequiredResourceAccess(in *[]msgraph.RequiredResourceAccess) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0)
	for _, requiredResourceAccess := range *in {
		resourceAppId := ""
		if requiredResourceAccess.ResourceAppId != nil {
			resourceAppId = *requiredResourceAccess.ResourceAppId
		}

		result = append(result, map[string]interface{}{
			"resource_app_id": resourceAppId,
			"resource_access": flattenApplicationResourceAccess(requiredResourceAccess.ResourceAccess),
		})
	}

	return result
}

func flattenApplicationResourceAccess(in *[]msgraph.ResourceAccess) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	accesses := make([]interface{}, 0)
	for _, resourceAccess := range *in {
		access := make(map[string]interface{})
		if resourceAccess.ID != nil {
			access["id"] = *resourceAccess.ID
		}
		access["type"] = resourceAccess.Type
		accesses = append(accesses, access)
	}

	return accesses
}

func flattenApplicationSpa(in *msgraph.ApplicationSpa) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	return []map[string]interface{}{{
		"redirect_uris": tf.FlattenStringSlicePtr(in.RedirectUris),
	}}
}

func flattenApplicationWeb(in *msgraph.ApplicationWeb) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	homepageUrl := ""
	if in.HomePageUrl != nil {
		homepageUrl = string(*in.HomePageUrl)
	}
	logoutUrl := ""
	if in.LogoutUrl != nil {
		logoutUrl = string(*in.LogoutUrl)
	}

	return []map[string]interface{}{{
		"homepage_url":   homepageUrl,
		"logout_url":     logoutUrl,
		"redirect_uris":  tf.FlattenStringSlicePtr(in.RedirectUris),
		"implicit_grant": flattenApplicationImplicitGrant(in.ImplicitGrantSettings),
	}}
}
