package msgraph

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
)

func ApplicationFindByName(ctx context.Context, client *msgraph.ApplicationsClient, displayName string) (*msgraph.Application, error) {
	filter := fmt.Sprintf("displayName eq '%s'", displayName)
	result, _, err := client.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to list Applications with filter %q: %+v", filter, err)
	}

	if result != nil {
		for _, app := range *result {
			if app.DisplayName != nil && *app.DisplayName == displayName {
				return &app, nil
			}
		}
	}

	return nil, nil
}

func ApplicationFlattenAppRoles(in *[]msgraph.AppRole) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	appRoles := make([]map[string]interface{}, 0)
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
		appRoles = append(appRoles, map[string]interface{}{
			"id":                   roleId,
			"allowed_member_types": allowedMemberTypes,
			"description":          description,
			"display_name":         displayName,
			"enabled":              enabled,
			"is_enabled":           enabled, // TODO: remove in v2.0
			"value":                value,
		})
	}

	return appRoles
}

func ApplicationFlattenOAuth2Permissions(in *[]msgraph.PermissionScope) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0)
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

		enabled := false
		if p.IsEnabled != nil && *p.IsEnabled {
			enabled = true
		}

		permType := ""
		if v := p.Type; v != nil {
			permType = *v
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

		result = append(result, map[string]interface{}{
			"admin_consent_description":  adminConsentDescription,
			"admin_consent_display_name": adminConsentDisplayName,
			"id":                         id,
			"is_enabled":                 enabled,
			"type":                       permType,
			"user_consent_description":   userConsentDescription,
			"user_consent_display_name":  userConsentDisplayName,
			"value":                      value,
		})
	}

	return result
}

func ApplicationSetAppRoles(ctx context.Context, client *msgraph.ApplicationsClient, application *msgraph.Application, newRoles *[]msgraph.AppRole) error {
	if application.ID == nil {
		return fmt.Errorf("cannot use Application model with nil ID")
	}

	if newRoles == nil {
		newRoles = &[]msgraph.AppRole{}
	}

	// Roles must be disabled before they can be edited or removed.
	// Since we cannot match them by ID, we have to disable all the roles, and replace them in one pass.
	app, status, err := client.Get(ctx, *application.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return fmt.Errorf("application with ID %q was not found", *application.ID)
		}

		return fmt.Errorf("retrieving Application with object ID %q: %+v", *application.ID, err)
	}

	// don't update if no changes to be made
	if app.AppRoles != nil && reflect.DeepEqual(*app.AppRoles, *newRoles) {
		return nil
	}

	// first disable any existing roles
	if app.AppRoles != nil && len(*app.AppRoles) > 0 {
		properties := msgraph.Application{
			ID:       application.ID,
			AppRoles: app.AppRoles,
		}

		for _, role := range *properties.AppRoles {
			role.IsEnabled = utils.Bool(false)
		}

		if _, err := client.Update(ctx, properties); err != nil {
			return fmt.Errorf("disabling App Roles for Application with object ID %q: %+v", *application.ID, err)
		}
	}

	// then set the new roles
	properties := msgraph.Application{
		ID:       application.ID,
		AppRoles: newRoles,
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return fmt.Errorf("setting App Roles for Application with object ID %q: %+v", *application.ID, err)
	}

	return nil
}

func ApplicationSetOAuth2PermissionScopes(ctx context.Context, client *msgraph.ApplicationsClient, application *msgraph.Application, newScopes *[]msgraph.PermissionScope) error {
	if application.ID == nil {
		return fmt.Errorf("Cannot use Application model with nil ID")
	}

	if newScopes == nil {
		newScopes = &[]msgraph.PermissionScope{}
	}

	// OAuth2 Permission Scopes must be disabled before they can be edited or removed.
	// Since we cannot match them by ID, we have to disable all the scopes, and replace them in one pass.
	app, status, err := client.Get(ctx, *application.ID)
	if err != nil {
		if status == http.StatusNotFound {
			return fmt.Errorf("application with ID %q was not found", *application.ID)
		}

		return fmt.Errorf("retrieving Application with object ID %q: %+v", *application.ID, err)
	}

	// don't update if no changes to be made
	if app.Api != nil && app.Api.OAuth2PermissionScopes != nil && reflect.DeepEqual(*app.Api.OAuth2PermissionScopes, *newScopes) {
		return nil
	}

	// first disable any existing scopes
	if app.Api != nil && app.Api.OAuth2PermissionScopes != nil && len(*app.Api.OAuth2PermissionScopes) > 0 {
		properties := msgraph.Application{
			ID: application.ID,
			Api: &msgraph.ApplicationApi{
				OAuth2PermissionScopes: app.Api.OAuth2PermissionScopes,
			},
		}

		for _, scope := range *properties.Api.OAuth2PermissionScopes {
			scope.IsEnabled = utils.Bool(false)
		}

		if _, err := client.Update(ctx, properties); err != nil {
			return fmt.Errorf("disabling OAuth2 Permission Scopes for Application with object ID %q: %+v", *application.ID, err)
		}
	}

	// then set the new scopes
	properties := msgraph.Application{
		ID: application.ID,
		Api: &msgraph.ApplicationApi{
			OAuth2PermissionScopes: newScopes,
		},
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return fmt.Errorf("setting OAuth2 Permission Scopes for Application with object ID %q: %+v", *application.ID, err)
	}

	return nil
}

func ApplicationSetOwners(ctx context.Context, client *msgraph.ApplicationsClient, application *msgraph.Application, desiredOwners []string) error {
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

	if ownersForRemoval != nil {
		if _, err = client.RemoveOwners(ctx, *application.ID, &ownersForRemoval); err != nil {
			return fmt.Errorf("removing owner from Application with object ID %q: %+v", *application.ID, err)
		}
	}

	if ownersToAdd != nil {
		for _, m := range ownersToAdd {
			application.AppendOwner(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, m)
		}

		if _, err := client.AddOwners(ctx, application); err != nil {
			return fmt.Errorf("adding owners to Application with object ID %q: %+v", *application.ID, err)
		}
	}
	return nil
}

func AppRoleFindById(app *msgraph.Application, roleId string) (*msgraph.AppRole, error) {
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

func OAuth2PermissionFindById(app *msgraph.Application, scopeId string) (*msgraph.PermissionScope, error) {
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