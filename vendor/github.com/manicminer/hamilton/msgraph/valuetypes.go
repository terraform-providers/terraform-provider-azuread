package msgraph

import "encoding/json"

// StringNullWhenEmpty is a string type that marshals its JSON representation as null when set to its zero value.
// Can be used with a pointer reference with the `omitempty` tag to omit a field when the pointer is nil, but send a
// JSON null value when the string is empty.
type StringNullWhenEmpty string

func (s StringNullWhenEmpty) MarshalJSON() ([]byte, error) {
	if s == "" {
		return []byte("null"), nil
	}
	return json.Marshal(string(s))
}

type ApplicationExtensionDataType string

const (
	ApplicationExtensionDataTypeBinary       ApplicationExtensionDataType = "Binary"
	ApplicationExtensionDataTypeBoolean      ApplicationExtensionDataType = "Boolean"
	ApplicationExtensionDataTypeDateTime     ApplicationExtensionDataType = "DateTime"
	ApplicationExtensionDataTypeInteger      ApplicationExtensionDataType = "Integer"
	ApplicationExtensionDataTypeLargeInteger ApplicationExtensionDataType = "LargeInteger"
	ApplicationExtensionDataTypeString       ApplicationExtensionDataType = "String"
)

type ApplicationExtensionTargetObject string

const (
	ApplicationExtensionTargetObjectApplication  ApplicationExtensionTargetObject = "Application"
	ApplicationExtensionTargetObjectDevice       ApplicationExtensionTargetObject = "Device"
	ApplicationExtensionTargetObjectGroup        ApplicationExtensionTargetObject = "Group"
	ApplicationExtensionTargetObjectOrganization ApplicationExtensionTargetObject = "Organization"
	ApplicationExtensionTargetObjectUser         ApplicationExtensionTargetObject = "User"
)

type AppRoleAllowedMemberType string

const (
	AppRoleAllowedMemberTypeApplication AppRoleAllowedMemberType = "Application"
	AppRoleAllowedMemberTypeUser        AppRoleAllowedMemberType = "User"
)

type BodyType string

const (
	BodyTypeText BodyType = "text"
	BodyTypeHtml BodyType = "html"
)

type ExtensionSchemaTargetType string

const (
	ExtensionSchemaTargetTypeAdministrativeUnit ExtensionSchemaTargetType = "AdministrativeUnit"
	ExtensionSchemaTargetTypeContact            ExtensionSchemaTargetType = "Contact"
	ExtensionSchemaTargetTypeDevice             ExtensionSchemaTargetType = "Device"
	ExtensionSchemaTargetTypeEvent              ExtensionSchemaTargetType = "Event"
	ExtensionSchemaTargetTypePost               ExtensionSchemaTargetType = "Post"
	ExtensionSchemaTargetTypeGroup              ExtensionSchemaTargetType = "Group"
	ExtensionSchemaTargetTypeMessage            ExtensionSchemaTargetType = "Message"
	ExtensionSchemaTargetTypeOrganization       ExtensionSchemaTargetType = "Organization"
	ExtensionSchemaTargetTypeUser               ExtensionSchemaTargetType = "User"
)

type ExtensionSchemaPropertyDataType string

const (
	ExtensionSchemaPropertyDataBinary   ExtensionSchemaPropertyDataType = "Binary"
	ExtensionSchemaPropertyDataBoolean  ExtensionSchemaPropertyDataType = "Boolean"
	ExtensionSchemaPropertyDataDateTime ExtensionSchemaPropertyDataType = "DateTime"
	ExtensionSchemaPropertyDataInteger  ExtensionSchemaPropertyDataType = "Integer"
	ExtensionSchemaPropertyDataString   ExtensionSchemaPropertyDataType = "String"
)

type GroupType string

const (
	GroupTypeUnified GroupType = "Unified"
)

type GroupMembershipClaim string

const (
	GroupMembershipClaimAll              GroupMembershipClaim = "All"
	GroupMembershipClaimNone             GroupMembershipClaim = "None"
	GroupMembershipClaimApplicationGroup GroupMembershipClaim = "ApplicationGroup"
	GroupMembershipClaimDirectoryRole    GroupMembershipClaim = "DirectoryRole"
	GroupMembershipClaimSecurityGroup    GroupMembershipClaim = "SecurityGroup"
)

type KeyCredentialType string

const (
	KeyCredentialTypeAsymmetricX509Cert  KeyCredentialType = "AsymmetricX509Cert"
	KeyCredentialTypeX509CertAndPassword KeyCredentialType = "X509CertAndPassword"
)

type KeyCredentialUsage string

const (
	KeyCredentialUsageSign   KeyCredentialUsage = "Sign"
	KeyCredentialUsageVerify KeyCredentialUsage = "Verify"
)

type PermissionScopeType string

const (
	PermissionScopeTypeAdmin PermissionScopeType = "Admin"
	PermissionScopeTypeUser  PermissionScopeType = "User"
)

type ResourceAccessType string

const (
	ResourceAccessTypeRole  ResourceAccessType = "Role"
	ResourceAccessTypeScope ResourceAccessType = "Scope"
)

type SignInAudience string

const (
	SignInAudienceAzureADMyOrg                       SignInAudience = "AzureADMyOrg"
	SignInAudienceAzureADMultipleOrgs                SignInAudience = "AzureADMultipleOrgs"
	SignInAudienceAzureADandPersonalMicrosoftAccount SignInAudience = "AzureADandPersonalMicrosoftAccount"
	SignInAudiencePersonalMicrosoftAccount           SignInAudience = "PersonalMicrosoftAccount"
)
