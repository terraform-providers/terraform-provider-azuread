package msgraph

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/errors"
)

type AddIn struct {
	ID         *string          `json:"id,omitempty"`
	Properties *[]AddInKeyValue `json:"properties,omitempty"`
	Type       *string          `json:"type,omitempty"`
}

type AddInKeyValue struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type ApiPreAuthorizedApplication struct {
	AppId         *string   `json:"appId,omitempty"`
	PermissionIds *[]string `json:"permissionIds,omitempty"`
}

// Application describes an Application object.
type Application struct {
	ID                            *string                   `json:"id,omitempty"`
	AddIns                        *[]AddIn                  `json:"addIns,omitempty"`
	Api                           *ApplicationApi           `json:"api,omitempty"`
	AppId                         *string                   `json:"appId,omitempty"`
	AppRoles                      *[]AppRole                `json:"appRoles,omitempty"`
	CreatedDateTime               *time.Time                `json:"createdDateTime,omitempty"`
	DefaultRedirectUri            *string                   `json:"defaultRedirectUri,omitempty"`
	DeletedDateTime               *time.Time                `json:"deletedDateTime,omitempty"`
	DisabledByMicrosoftStatus     interface{}               `json:"disabledByMicrosoftStatus,omitempty"`
	DisplayName                   *string                   `json:"displayName,omitempty"`
	GroupMembershipClaims         *[]GroupMembershipClaim   `json:"groupMembershipClaims,omitempty"`
	IdentifierUris                *[]string                 `json:"identifierUris,omitempty"`
	Info                          *InformationalUrl         `json:"info,omitempty"`
	IsAuthorizationServiceEnabled *bool                     `json:"isAuthorizationServiceEnabled,omitempty"`
	IsDeviceOnlyAuthSupported     *bool                     `json:"isDeviceOnlyAuthSupported,omitempty"`
	IsFallbackPublicClient        *bool                     `json:"isFallbackPublicClient,omitempty"`
	IsManagementRestricted        *bool                     `json:"isManagementRestricted,omitempty"`
	KeyCredentials                *[]KeyCredential          `json:"keyCredentials,omitempty"`
	Oauth2RequirePostResponse     *bool                     `json:"oauth2RequirePostResponse,omitempty"`
	OnPremisesPublishing          *OnPremisesPublishing     `json:"onPremisePublishing,omitempty"`
	OptionalClaims                *OptionalClaims           `json:"optionalClaims,omitempty"`
	ParentalControlSettings       *ParentalControlSettings  `json:"parentalControlSettings,omitempty"`
	PasswordCredentials           *[]PasswordCredential     `json:"passwordCredentials,omitempty"`
	PublicClient                  *PublicClient             `json:"publicClient,omitempty"`
	PublisherDomain               *string                   `json:"publisherDomain,omitempty"`
	RequiredResourceAccess        *[]RequiredResourceAccess `json:"requiredResourceAccess,omitempty"`
	SignInAudience                SignInAudience            `json:"signInAudience,omitempty"`
	Tags                          *[]string                 `json:"tags,omitempty"`
	TokenEncryptionKeyId          *string                   `json:"tokenEncryptionKeyId,omitempty"`
	UniqueName                    *string                   `json:"uniqueName,omitempty"`
	VerifiedPublisher             *VerifiedPublisher        `json:"verifiedPublisher,omitempty"`
	Web                           *ApplicationWeb           `json:"web,omitempty"`

	Owners *[]string `json:"owners@odata.bind,omitempty"`
}

func (a Application) MarshalJSON() ([]byte, error) {
	var val *StringNullWhenEmpty
	if a.GroupMembershipClaims != nil {
		claims := make([]string, 0)
		for _, c := range *a.GroupMembershipClaims {
			claims = append(claims, string(c))
		}
		theClaims := StringNullWhenEmpty(strings.Join(claims, ","))
		val = &theClaims
	}
	type application Application
	return json.Marshal(&struct {
		GroupMembershipClaims *StringNullWhenEmpty `json:"groupMembershipClaims,omitempty"`
		*application
	}{
		GroupMembershipClaims: val,
		application:           (*application)(&a),
	})
}

func (a *Application) UnmarshalJSON(data []byte) error {
	type application Application
	app := &struct {
		GroupMembershipClaims *string `json:"groupMembershipClaims"`
		*application
	}{
		application: (*application)(a),
	}
	if err := json.Unmarshal(data, &app); err != nil {
		return err
	}
	if app.GroupMembershipClaims != nil {
		var groupMembershipClaims []GroupMembershipClaim
		for _, c := range strings.Split(*app.GroupMembershipClaims, ",") {
			groupMembershipClaims = append(groupMembershipClaims, GroupMembershipClaim(strings.TrimSpace(c)))
		}
		a.GroupMembershipClaims = &groupMembershipClaims
	}
	return nil
}

// AppendOwner appends a new owner object URI to the Owners slice.
func (a *Application) AppendOwner(endpoint environments.ApiEndpoint, apiVersion ApiVersion, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var owners []string
	if a.Owners != nil {
		owners = *a.Owners
	}
	owners = append(owners, val)
	a.Owners = &owners
}

// AppendAppRole adds a new AppRole to an Application, checking to see if it already exists.
func (a *Application) AppendAppRole(role AppRole) error {
	if role.ID == nil {
		return goerrors.New("ID of new role is nil")
	}

	cap := 1
	if a.AppRoles != nil {
		cap += len(*a.AppRoles)
	}

	newRoles := make([]AppRole, 1, cap)
	newRoles[0] = role

	for _, v := range *a.AppRoles {
		if v.ID != nil && *v.ID == *role.ID {
			return &errors.AlreadyExistsError{Obj: "AppRole", Id: *role.ID}
		}
		newRoles = append(newRoles, v)
	}

	a.AppRoles = &newRoles
	return nil
}

// RemoveAppRole removes an AppRole from an Application.
func (a *Application) RemoveAppRole(role AppRole) error {
	if role.ID == nil {
		return goerrors.New("ID of role is nil")
	}

	if a.AppRoles == nil {
		return goerrors.New("no roles to remove")
	}

	appRoles := make([]AppRole, 0, len(*a.AppRoles))
	for _, v := range *a.AppRoles {
		if v.ID == nil || *v.ID != *role.ID {
			appRoles = append(appRoles, v)
		}
	}

	if len(appRoles) == len(*a.AppRoles) {
		return goerrors.New("could not find role to remove")
	}

	a.AppRoles = &appRoles
	return nil
}

// UpdateAppRole amends an existing AppRole defined in an Application.
func (a *Application) UpdateAppRole(role AppRole) error {
	if role.ID == nil {
		return goerrors.New("ID of role is nil")
	}

	if a.AppRoles == nil {
		return goerrors.New("no roles to update")
	}

	appRoles := *a.AppRoles
	for i, v := range appRoles {
		if v.ID != nil && *v.ID == *role.ID {
			appRoles[i] = role
			break
		}
	}

	a.AppRoles = &appRoles
	return nil
}

type ApplicationApi struct {
	AcceptMappedClaims          *bool                          `json:"acceptMappedClaims,omitempty"`
	KnownClientApplications     *[]string                      `json:"knownClientApplications,omitempty"`
	OAuth2PermissionScopes      *[]PermissionScope             `json:"oauth2PermissionScopes,omitempty"`
	PreAuthorizedApplications   *[]ApiPreAuthorizedApplication `json:"preAuthorizedApplications,omitempty"`
	RequestedAccessTokenVersion *int32                         `json:"requestedAccessTokenVersion,omitempty"`
}

// AppendOAuth2PermissionScope adds a new ApplicationOAuth2PermissionScope to an ApplicationApi, checking to see if it already exists.
func (a *ApplicationApi) AppendOAuth2PermissionScope(scope PermissionScope) error {
	if scope.ID == nil {
		return goerrors.New("ID of new scope is nil")
	}

	cap := 1
	if a.OAuth2PermissionScopes != nil {
		cap += len(*a.OAuth2PermissionScopes)
	}

	newScopes := make([]PermissionScope, 1, cap)
	newScopes[0] = scope

	for _, v := range *a.OAuth2PermissionScopes {
		if v.ID != nil && *v.ID == *scope.ID {
			return &errors.AlreadyExistsError{Obj: "OAuth2PermissionScope", Id: *scope.ID}
		}
		newScopes = append(newScopes, v)
	}

	a.OAuth2PermissionScopes = &newScopes
	return nil
}

// RemoveOAuth2PermissionScope removes an ApplicationOAuth2PermissionScope from an ApplicationApi.
func (a *ApplicationApi) RemoveOAuth2PermissionScope(scope PermissionScope) error {
	if scope.ID == nil {
		return goerrors.New("ID of scope is nil")
	}

	if a.OAuth2PermissionScopes == nil {
		return goerrors.New("no scopes to remove")
	}

	apiScopes := make([]PermissionScope, 0, len(*a.OAuth2PermissionScopes))
	for _, v := range *a.OAuth2PermissionScopes {
		if v.ID == nil || *v.ID != *scope.ID {
			apiScopes = append(apiScopes, v)
		}
	}

	if len(apiScopes) == len(*a.OAuth2PermissionScopes) {
		return goerrors.New("could not find scope to remove")
	}

	a.OAuth2PermissionScopes = &apiScopes
	return nil
}

// UpdateOAuth2PermissionScope amends an existing ApplicationOAuth2PermissionScope defined in an ApplicationApi.
func (a *ApplicationApi) UpdateOAuth2PermissionScope(scope PermissionScope) error {
	if scope.ID == nil {
		return goerrors.New("ID of scope is nil")
	}

	if a.OAuth2PermissionScopes == nil {
		return goerrors.New("no scopes to update")
	}

	apiScopes := *a.OAuth2PermissionScopes
	for i, v := range apiScopes {
		if v.ID != nil && *v.ID == *scope.ID {
			apiScopes[i] = scope
			break
		}
	}

	a.OAuth2PermissionScopes = &apiScopes
	return nil
}

type ApplicationEnforcedRestrictionsSessionControl struct {
	IsEnabled *bool `json:"isEnabled,omitempty"`
}

type ApplicationWeb struct {
	HomePageUrl           *StringNullWhenEmpty   `json:"homePageUrl,omitempty"`
	ImplicitGrantSettings *ImplicitGrantSettings `json:"implicitGrantSettings,omitempty"`
	LogoutUrl             *StringNullWhenEmpty   `json:"logoutUrl,omitempty"`
	RedirectUris          *[]string              `json:"redirectUris,omitempty"`
}

type AppRole struct {
	ID                 *string                     `json:"id,omitempty"`
	AllowedMemberTypes *[]AppRoleAllowedMemberType `json:"allowedMemberTypes,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	DisplayName        *string                     `json:"displayName,omitempty"`
	IsEnabled          *bool                       `json:"isEnabled,omitempty"`
	Origin             *string                     `json:"origin,omitempty"`
	Value              *string                     `json:"value,omitempty"`
}

type AppRoleAssignment struct {
	Id                   *string    `json:"id,omitempty"`
	DeletedDateTime      *time.Time `json:"deletedDateTime,omitempty"`
	AppRoleId            *string    `json:"appRoleId,omitempty"`
	CreatedDateTime      *time.Time `json:"createdDateTime,omitempty"`
	PrincipalDisplayName *string    `json:"principalDisplayName,omitempty"`
	PrincipalId          *string    `json:"principalId,omitempty"`
	PrincipalType        *string    `json:"principalType,omitempty"`
	ResourceDisplayName  *string    `json:"resourceDisplayName,omitempty"`
	ResourceId           *string    `json:"resourceId,omitempty"`
}

type BaseNamedLocation struct {
	ODataType        *string    `json:"@odata.type,omitempty"`
	ID               *string    `json:"id,omitempty"`
	DisplayName      *string    `json:"displayName,omitempty"`
	CreatedDateTime  *time.Time `json:"createdDateTime,omitempty"`
	ModifiedDateTime *time.Time `json:"modifiedDateTime,omitempty"`
}

type CloudAppSecurityControl struct {
	IsEnabled            *bool   `json:"isEnabled,omitempty"`
	CloudAppSecurityType *string `json:"cloudAppSecurityType,omitempty"`
}

// ConditionalAccessPolicy describes an Conditional Access Policy object.
type ConditionalAccessPolicy struct {
	Conditions       *ConditionalAccessConditionSet    `json:"conditions,omitempty"`
	CreatedDateTime  *time.Time                        `json:"createdDateTime,omitempty"`
	DisplayName      *string                           `json:"displayName,omitempty"`
	GrantControls    *ConditionalAccessGrantControls   `json:"grantControls,omitempty"`
	ID               *string                           `json:"id,omitempty"`
	ModifiedDateTime *time.Time                        `json:"modifiedDateTime,omitempty"`
	SessionControls  *ConditionalAccessSessionControls `json:"sessionControls,omitempty"`
	State            *string                           `json:"state,omitempty"`
}

type ConditionalAccessConditionSet struct {
	Applications     *ConditionalAccessApplications `json:"applications,omitempty"`
	Users            *ConditionalAccessUsers        `json:"users,omitempty"`
	ClientAppTypes   *[]string                      `json:"clientAppTypes,omitempty"`
	Locations        *ConditionalAccessLocations    `json:"locations,omitempty"`
	Platforms        *ConditionalAccessPlatforms    `json:"platforms,omitempty"`
	SignInRiskLevels *[]string                      `json:"signInRiskLevels,omitempty"`
	UserRiskLevels   *[]string                      `json:"userRiskLevels,omitempty"`
}

type ConditionalAccessApplications struct {
	IncludeApplications *[]string `json:"includeApplications,omitempty"`
	ExcludeApplications *[]string `json:"excludeApplications,omitempty"`
	IncludeUserActions  *[]string `json:"includeUserActions,omitempty"`
}

type ConditionalAccessUsers struct {
	IncludeUsers  *[]string `json:"includeUsers,omitempty"`
	ExcludeUsers  *[]string `json:"excludeUsers,omitempty"`
	IncludeGroups *[]string `json:"includeGroups,omitempty"`
	ExcludeGroups *[]string `json:"excludeGroups,omitempty"`
	IncludeRoles  *[]string `json:"includeRoles,omitempty"`
	ExcludeRoles  *[]string `json:"excludeRoles,omitempty"`
}

type ConditionalAccessLocations struct {
	IncludeLocations *[]string `json:"includeLocations,omitempty"`
	ExcludeLocations *[]string `json:"excludeLocations,omitempty"`
}

type ConditionalAccessPlatforms struct {
	IncludePlatforms *[]string `json:"includePlatforms,omitempty"`
	ExcludePlatforms *[]string `json:"excludePlatforms,omitempty"`
}

type ConditionalAccessGrantControls struct {
	Operator                    *string   `json:"operator,omitempty"`
	BuiltInControls             *[]string `json:"builtInControls,omitempty"`
	CustomAuthenticationFactors *[]string `json:"customAuthenticationFactors,omitempty"`
	TermsOfUse                  *[]string `json:"termsOfUse,omitempty"`
}

type ConditionalAccessSessionControls struct {
	ApplicationEnforcedRestrictions *ApplicationEnforcedRestrictionsSessionControl `json:"applicationEnforcedRestrictions,omitempty"`
	CloudAppSecurity                *CloudAppSecurityControl                       `json:"cloudAppSecurity,omitempty"`
	PersistentBrowser               *PersistentBrowserSessionControl               `json:"persistentBrowser,omitempty"`
	SignInFrequency                 *SignInFrequencySessionControl                 `json:"signInFrequency,omitempty"`
}

// CountryNamedLocation describes an Country Named Location object.
type CountryNamedLocation struct {
	*BaseNamedLocation
	CountriesAndRegions               *[]string `json:"countriesAndRegions,omitempty"`
	IncludeUnknownCountriesAndRegions *bool     `json:"includeUnknownCountriesAndRegions,omitempty"`
}

// DirectoryRoleTemplate describes a Directory Role Template.
type DirectoryRoleTemplate struct {
	ID              *string    `json:"id,omitempty"`
	DeletedDateTime *time.Time `json:"deletedDateTime,omitempty"`
	Description     *string    `json:"description,omitempty"`
	DisplayName     *string    `json:"displayName,omitempty"`
}

type DirectoryRole struct {
	ID             *string `json:"id,omitempty"`
	Description    *string `json:"description,omitempty"`
	DisplayName    *string `json:"displayName,omitempty"`
	RoleTemplateId *string `json:"roleTemplateId,omitempty"`

	Members *[]string `json:"-"`
}

// AppendMember appends a new member object URI to the Members slice.
func (d *DirectoryRole) AppendMember(endpoint environments.ApiEndpoint, apiVersion ApiVersion, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var members []string
	if d.Members != nil {
		members = *d.Members
	}
	members = append(members, val)
	d.Members = &members
}

// Domain describes a Domain object.
type Domain struct {
	ID                               *string   `json:"id,omitempty"`
	AuthenticationType               *string   `json:"authenticationType,omitempty"`
	IsAdminManaged                   *bool     `json:"isAdminManaged,omitempty"`
	IsDefault                        *bool     `json:"isDefault,omitempty"`
	IsInitial                        *bool     `json:"isInitial,omitempty"`
	IsRoot                           *bool     `json:"isRoot,omitempty"`
	IsVerified                       *bool     `json:"isVerified,omitempty"`
	PasswordNotificationWindowInDays *int      `json:"passwordNotificationWindowInDays,omitempty"`
	PasswordValidityPeriodInDays     *int      `json:"passwordValidityPeriodInDays,omitempty"`
	SupportedServices                *[]string `json:"supportedServices,omitempty"`

	State *DomainState `json:"state,omitempty"`
}

type DomainState struct {
	LastActionDateTime *time.Time `json:"lastActionDateTime,omitempty"`
	Operation          *string    `json:"operation,omitempty"`
	Status             *string    `json:"status,omitempty"`
}

type EmailAddress struct {
	Address *string `json:"address,omitempty"`
	Name    *string `json:"name,omitempty"`
}

// Group describes a Group object.
type Group struct {
	ID                            *string                             `json:"id,omitempty"`
	AllowExternalSenders          *string                             `json:"allowExternalSenders,omitempty"`
	AssignedLabels                *[]GroupAssignedLabel               `json:"assignedLabels,omitempty"`
	AssignedLicenses              *[]GroupAssignedLicense             `json:"assignLicenses,omitempty"`
	AutoSubscribeNewMembers       *bool                               `json:"autoSubscribeNewMembers,omitempty"`
	Classification                *string                             `json:"classification,omitempty"`
	CreatedDateTime               *time.Time                          `json:"createdDateTime,omitempty"`
	DeletedDateTime               *time.Time                          `json:"deletedDateTime,omitempty"`
	Description                   *StringNullWhenEmpty                `json:"description,omitempty"`
	DisplayName                   *string                             `json:"displayName,omitempty"`
	ExpirationDateTime            *time.Time                          `json:"expirationDateTime,omitempty"`
	GroupTypes                    []GroupType                         `json:"groupTypes,omitempty"`
	HasMembersWithLicenseErrors   *bool                               `json:"hasMembersWithLicenseErrors,omitempty"`
	HideFromAddressLists          *bool                               `json:"hideFromAddressLists,omitempty"`
	HideFromOutlookClients        *bool                               `json:"hideFromOutlookClients,omitempty"`
	IsSubscribedByMail            *bool                               `json:"isSubscribedByMail,omitempty"`
	LicenseProcessingState        *string                             `json:"licenseProcessingState,omitempty"`
	Mail                          *string                             `json:"mail,omitempty"`
	MailEnabled                   *bool                               `json:"mailEnabled,omitempty"`
	MailNickname                  *string                             `json:"mailNickname,omitempty"`
	MembershipRule                *string                             `json:"membershipRule,omitempty"`
	MembershipRuleProcessingState *string                             `json:"membershipRuleProcessingState,omitempty"`
	OnPremisesDomainName          *string                             `json:"onPremisesDomainName,omitempty"`
	OnPremisesLastSyncDateTime    *time.Time                          `json:"onPremisesLastSyncDateTime,omitempty"`
	OnPremisesNetBiosName         *string                             `json:"onPremisesNetBiosName,omitempty"`
	OnPremisesProvisioningErrors  *[]GroupOnPremisesProvisioningError `json:"onPremisesProvisioningErrors,omitempty"`
	OnPremisesSamAccountName      *string                             `json:"onPremisesSamAccountName,omitempty"`
	OnPremisesSecurityIdentifier  *string                             `json:"onPremisesSecurityIdentifier,omitempty"`
	OnPremisesSyncEnabled         *bool                               `json:"onPremisesSyncEnabled,omitempty"`
	PreferredDataLocation         *string                             `json:"preferredDataLocation,omitempty"`
	PreferredLanguage             *string                             `json:"preferredLanguage,omitempty"`
	ProxyAddresses                *[]string                           `json:"proxyAddresses,omitempty"`
	RenewedDateTime               *time.Time                          `json:"renewedDateTime,omitempty"`
	SecurityEnabled               *bool                               `json:"securityEnabled,omitempty"`
	SecurityIdentifier            *string                             `json:"securityIdentifier,omitempty"`
	Theme                         *string                             `json:"theme,omitempty"`
	UnseenCount                   *int                                `json:"unseenCount,omitempty"`
	Visibility                    *string                             `json:"visibility,omitempty"`
	IsAssignableToRole            *bool                               `json:"isAssignableToRole,omitempty"`

	Members *[]string `json:"members@odata.bind,omitempty"`
	Owners  *[]string `json:"owners@odata.bind,omitempty"`
}

// AppendMember appends a new member object URI to the Members slice.
func (g *Group) AppendMember(endpoint environments.ApiEndpoint, apiVersion ApiVersion, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var members []string
	if g.Members != nil {
		members = *g.Members
	}
	members = append(members, val)
	g.Members = &members
}

// AppendOwner appends a new owner object URI to the Owners slice.
func (g *Group) AppendOwner(endpoint environments.ApiEndpoint, apiVersion ApiVersion, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var owners []string
	if g.Owners != nil {
		owners = *g.Owners
	}
	owners = append(owners, val)
	g.Owners = &owners
}

type GroupAssignedLabel struct {
	LabelId     *string `json:"labelId,omitempty"`
	DisplayName *string `json:"displayNanme,omitempty"`
}

type GroupAssignedLicense struct {
	DisabledPlans *[]string `json:"disabledPlans,omitempty"`
	SkuId         *string   `json:"skuId,omitempty"`
}

type GroupOnPremisesProvisioningError struct {
	Category             *string   `json:"category,omitempty"`
	OccurredDateTime     time.Time `json:"occurredDateTime,omitempty"`
	PropertyCausingError *string   `json:"propertyCausingError,omitempty"`
	Value                *string   `json:"value,omitempty"`
}

type IdentityProvider struct {
	ODataType    *string `json:"@odata.type,omitempty"`
	ID           *string `json:"id,omitempty"`
	ClientId     *string `json:"clientId,omitempty"`
	ClientSecret *string `json:"clientSecret,omitempty"`
	Type         *string `json:"identityProviderType,omitempty"`
	Name         *string `json:"displayName,omitempty"`
}

type ImplicitGrantSettings struct {
	EnableAccessTokenIssuance *bool `json:"enableAccessTokenIssuance,omitempty"`
	EnableIdTokenIssuance     *bool `json:"enableIdTokenIssuance,omitempty"`
}

type InformationalUrl struct {
	LogoUrl             *string `json:"logoUrl,omitempty"`
	MarketingUrl        *string `json:"marketingUrl"`
	PrivacyStatementUrl *string `json:"privacyStatementUrl"`
	SupportUrl          *string `json:"supportUrl"`
	TermsOfServiceUrl   *string `json:"termsOfServiceUrl"`
}

// Invitation describes a Invitation object.
type Invitation struct {
	ID                      *string `json:"id,omitempty"`
	InvitedUserDisplayName  *string `json:"invitedUserDisplayName,omitempty"`
	InvitedUserEmailAddress *string `json:"invitedUserEmailAddress,omitempty"`
	SendInvitationMessage   *bool   `json:"sendInvitationMessage,omitempty"`
	InviteRedirectURL       *string `json:"inviteRedirectUrl,omitempty"`
	InviteRedeemURL         *string `json:"inviteRedeemUrl,omitempty"`
	Status                  *string `json:"status,omitempty"`
	InvitedUserType         *string `json:"invitedUserType,omitempty"`

	InvitedUserMessageInfo *InvitedUserMessageInfo `json:"invitedUserMessageInfo,omitempty"`
	InvitedUser            *User                   `json:"invitedUser,omitempty"`
}

type InvitedUserMessageInfo struct {
	CCRecipients          *[]Recipient `json:"ccRecipients,omitempty"`
	CustomizedMessageBody *string      `json:"customizedMessageBody,omitempty"`
	MessageLanguage       *string      `json:"messageLanguage,omitempty"`
}

// IPNamedLocation describes an IP Named Location object.
type IPNamedLocation struct {
	*BaseNamedLocation
	IPRanges  *[]IPNamedLocationIPRange `json:"ipRanges,omitempty"`
	IsTrusted *bool                     `json:"isTrusted,omitempty"`
}

type IPNamedLocationIPRange struct {
	CIDRAddress *string `json:"cidrAddress,omitempty"`
}

type ItemBody struct {
	Content     *string   `json:"content,omitempty"`
	ContentType *BodyType `json:"contentType,omitempty"`
}

type KerberosSignOnSettings struct {
	ServicePrincipalName       *string `json:"kerberosServicePrincipalName,omitempty"`
	SignOnMappingAttributeType *string `jsonL:"kerberosSignOnMappingAttributeType,omitempty"`
}

// KeyCredential describes a key (certificate) credential for an object.
type KeyCredential struct {
	CustomKeyIdentifier *string            `json:"customKeyIdentifier,omitempty"`
	DisplayName         *string            `json:"displayName,omitempty"`
	EndDateTime         *time.Time         `json:"endDateTime,omitempty"`
	KeyId               *string            `json:"keyId,omitempty"`
	StartDateTime       *time.Time         `json:"startDateTime,omitempty"`
	Type                KeyCredentialType  `json:"type"`
	Usage               KeyCredentialUsage `json:"usage"`
	Key                 *string            `json:"key,omitempty"`
}

type MailMessage struct {
	Message *Message `json:"message,omitempty"`
}

// Me describes the authenticated user.
type Me struct {
	ID                *string `json:"id"`
	DisplayName       *string `json:"displayName"`
	UserPrincipalName *string `json:"userPrincipalName"`
}

type Message struct {
	ID            *string      `json:"id,omitempty"`
	Subject       *string      `json:"subject,omitempty"`
	Body          *ItemBody    `json:"body,omitempty"`
	From          *Recipient   `json:"from,omitempty"`
	ToRecipients  *[]Recipient `json:"toRecipients,omitempty"`
	CcRecipients  *[]Recipient `json:"ccRecipients,omitempty"`
	BccRecipients *[]Recipient `json:"bccRecipients,omitempty"`
}

type NamedLocation interface{}

type OnPremisesPublishing struct {
	AlternateUrl                  *string `json:"alternateUrl,omitempty"`
	ApplicationServerTimeout      *string `json:"applicationServerTimeout,omitempty"`
	ApplicationType               *string `json:"applicationType,omitempty"`
	ExternalAuthenticationType    *string `json:"externalAuthenticationType,omitempty"`
	ExternalUrl                   *string `json:"externalUrl,omitempty"`
	InternalUrl                   *string `json:"internalUrl,omitempty"`
	IsHttpOnlyCookieEnabled       *bool   `json:"isHttpOnlyCookieEnabled,omitempty"`
	IsOnPremPublishingEnabled     *bool   `json:"isOnPremPublishingEnabled,omitempty"`
	IsPersistentCookieEnabled     *bool   `json:"isPersistentCookieEnabled,omitempty"`
	IsSecureCookieEnabled         *bool   `json:"isSecureCookieEnabled,omitempty"`
	IsTranslateHostHeaderEnabled  *bool   `json:"isTranslateHostHeaderEnabled,omitempty"`
	IsTranslateLinksInBodyEnabled *bool   `json:"isTranslateLinksInBodyEnabled,omitempty"`

	SingleSignOnSettings                     *OnPremisesPublishingSingleSignOn                             `json:"singleSignOnSettings,omitempty"`
	VerifiedCustomDomainCertificatesMetadata *OnPremisesPublishingVerifiedCustomDomainCertificatesMetadata `json:"verifiedCustomDomainCertificatesMetadata,omitempty"`
	VerifiedCustomDomainKeyCredential        *KeyCredential                                                `json:"verifiedCustomDomainKeyCredential,omitempty"`
	VerifiedCustomDomainPasswordCredential   *PasswordCredential                                           `json:"verifiedCustomDomainPasswordCredential,omitempty"`
}

type OnPremisesPublishingSingleSignOn struct {
	KerberosSignOnSettings *KerberosSignOnSettings `json:"kerberosSignOnSettings,omitempty"`
	SingleSignOnMode       *string                 `json:"singleSignOnMode,omitempty"`
}

type OnPremisesPublishingVerifiedCustomDomainCertificatesMetadata struct {
	ExpiryDate  *time.Time `json:"expiryDate,omitempty"`
	IssueDate   *time.Time `json:"issueDate,omitempty"`
	IssuerName  *string    `json:"issuerName,omitempty"`
	SubjectName *string    `json:"subjectName,omitempty"`
	Thumbprint  *string    `json:"thumbprint,omitempty"`
}

type OptionalClaim struct {
	AdditionalProperties *[]string `json:"additionalProperties,omitempty"`
	Essential            *bool     `json:"essential,omitempty"`
	Name                 *string   `json:"name,omitempty"`
	Source               *string   `json:"source,omitempty"`
}

type OptionalClaims struct {
	AccessToken *[]OptionalClaim `json:"accessToken,omitempty"`
	IdToken     *[]OptionalClaim `json:"idToken,omitempty"`
	Saml2Token  *[]OptionalClaim `json:"saml2Token,omitempty"`
}

type ParentalControlSettings struct {
	CountriesBlockedForMinors *[]string `json:"countriesBlockedForMinors,omitempty"`
	LegalAgeGroupRule         *string   `json:"legalAgeGroupRule,omitempty"`
}

// PasswordCredential describes a password credential for an object.
type PasswordCredential struct {
	CustomKeyIdentifier *string    `json:"customKeyIdentifier,omitempty"`
	DisplayName         *string    `json:"displayName,omitempty"`
	EndDateTime         *time.Time `json:"endDateTime,omitempty"`
	Hint                *string    `json:"hint,omitempty"`
	KeyId               *string    `json:"keyId,omitempty"`
	SecretText          *string    `json:"secretText,omitempty"`
	StartDateTime       *time.Time `json:"startDateTime,omitempty"`
}

type PasswordSingleSignOnSettings struct {
	Fields *[]SingleSignOnField `json:"fields,omitempty"`
}

type PermissionScope struct {
	ID                      *string             `json:"id,omitempty"`
	AdminConsentDescription *string             `json:"adminConsentDescription,omitempty"`
	AdminConsentDisplayName *string             `json:"adminConsentDisplayName,omitempty"`
	IsEnabled               *bool               `json:"isEnabled,omitempty"`
	Type                    PermissionScopeType `json:"type,omitempty"`
	UserConsentDescription  *string             `json:"userConsentDescription,omitempty"`
	UserConsentDisplayName  *string             `json:"userConsentDisplayName,omitempty"`
	Value                   *string             `json:"value,omitempty"`
}

type PersistentBrowserSessionControl struct {
	IsEnabled *bool   `json:"isEnabled,omitempty"`
	Mode      *string `json:"mode,omitempty"`
}

type PublicClient struct {
	RedirectUris *[]string `json:"redirectUris,omitempty"`
}

type Recipient struct {
	EmailAddress *EmailAddress `json:"emailAddress,omitempty"`
}

type RequiredResourceAccess struct {
	ResourceAccess *[]ResourceAccess `json:"resourceAccess,omitempty"`
	ResourceAppId  *string           `json:"resourceAppId,omitempty"`
}

type ResourceAccess struct {
	ID   *string            `json:"id,omitempty"`
	Type ResourceAccessType `json:"type,omitempty"`
}

type SamlSingleSignOnSettings struct {
	RelayState *string `json:"relayState,omitempty"`
}

// ServicePrincipal describes a Service Principal object.
type ServicePrincipal struct {
	ID                                  *string                       `json:"id,omitempty"`
	AccountEnabled                      *bool                         `json:"accountEnabled,omitempty"`
	AddIns                              *[]AddIn                      `json:"addIns,omitempty"`
	AlternativeNames                    *[]string                     `json:"alternativeNames,omitempty"`
	AppDisplayName                      *string                       `json:"appDisplayName,omitempty"`
	AppId                               *string                       `json:"appId,omitempty"`
	ApplicationTemplateId               *string                       `json:"applicationTemplateId,omitempty"`
	AppOwnerOrganizationId              *string                       `json:"appOwnerOrganizationId,omitempty"`
	AppRoleAssignmentRequired           *bool                         `json:"appRoleAssignmentRequired,omitempty"`
	AppRoles                            *[]AppRole                    `json:"appRoles,omitempty"`
	DeletedDateTime                     *time.Time                    `json:"deletedDateTime,omitempty"`
	DisplayName                         *string                       `json:"displayName,omitempty"`
	Homepage                            *string                       `json:"homepage,omitempty"`
	Info                                *InformationalUrl             `json:"info,omitempty"`
	KeyCredentials                      *[]KeyCredential              `json:"keyCredentials,omitempty"`
	LoginUrl                            *string                       `json:"loginUrl,omitempty"`
	LogoutUrl                           *string                       `json:"logoutUrl,omitempty"`
	NotificationEmailAddresses          *[]string                     `json:"notificationEmailAddresses,omitempty"`
	PasswordCredentials                 *[]PasswordCredential         `json:"passwordCredentials,omitempty"`
	PasswordSingleSignOnSettings        *PasswordSingleSignOnSettings `json:"passwordSingleSignOnSettings,omitempty"`
	PreferredSingleSignOnMode           *string                       `json:"preferredSingleSignOnMode,omitempty"`
	PreferredTokenSigningKeyEndDateTime *time.Time                    `json:"preferredTokenSigningKeyEndDateTime,omitempty"`
	PublishedPermissionScopes           *[]PermissionScope            `json:"publishedPermissionScopes,omitempty"`
	ReplyUrls                           *[]string                     `json:"replyUrls,omitempty"`
	SamlSingleSignOnSettings            *SamlSingleSignOnSettings     `json:"samlSingleSignOnSettings,omitempty"`
	ServicePrincipalNames               *[]string                     `json:"servicePrincipalNames,omitempty"`
	ServicePrincipalType                *string                       `json:"servicePrincipalType,omitempty"`
	SignInAudience                      SignInAudience                `json:"signInAudience,omitempty"`
	Tags                                *[]string                     `json:"tags,omitempty"`
	TokenEncryptionKeyId                *string                       `json:"tokenEncryptionKeyId,omitempty"`
	VerifiedPublisher                   *VerifiedPublisher            `json:"verifiedPublisher,omitempty"`

	Owners *[]string `json:"owners@odata.bind,omitempty"`
}

// AppendOwner appends a new owner object URI to the Owners slice.
func (a *ServicePrincipal) AppendOwner(endpoint string, apiVersion string, id string) {
	val := fmt.Sprintf("%s/%s/directoryObjects/%s", endpoint, apiVersion, id)
	var owners []string
	if a.Owners != nil {
		owners = *a.Owners
	}
	owners = append(owners, val)
	a.Owners = &owners
}

type SignInFrequencySessionControl struct {
	IsEnabled *bool   `json:"isEnabled,omitempty"`
	Type      *string `json:"type,omitempty"`
	Value     *int32  `json:"value,omitempty"`
}

type SingleSignOnField struct {
	CustomizedLabel *string `json:"customizedLabel,omitempty"`
	DefaultLabel    *string `json:"defaultLabel,omitempty"`
	FieldId         *string `json:"fieldId,omitempty"`
	Type            *string `json:"type,omitempty"`
}

// User describes a User object.
type User struct {
	ID                              *string              `json:"id,omitempty"`
	AboutMe                         *string              `json:"aboutMe,omitempty"`
	AccountEnabled                  *bool                `json:"accountEnabled,omitempty"`
	BusinessPhones                  *[]string            `json:"businessPhones,omitempty"`
	City                            *StringNullWhenEmpty `json:"city,omitempty"`
	CompanyName                     *StringNullWhenEmpty `json:"companyName,omitempty"`
	Country                         *StringNullWhenEmpty `json:"country,omitempty"`
	CreatedDateTime                 *time.Time           `json:"createdDateTime,omitempty"`
	CreationType                    *string              `json:"creationType,omitempty"`
	DeletedDateTime                 *time.Time           `json:"deletedDateTime,omitempty"`
	Department                      *StringNullWhenEmpty `json:"department,omitempty"`
	DisplayName                     *string              `json:"displayName,omitempty"`
	EmployeeHireDate                *time.Time           `json:"employeeHireDate,omitempty"`
	EmployeeId                      *string              `json:"employeeId,omitempty"`
	EmployeeType                    *string              `json:"employeeType,omitempty"`
	ExternalUserState               *string              `json:"externalUserState,omitempty"`
	FaxNumber                       *string              `json:"faxNumber,omitempty"`
	GivenName                       *StringNullWhenEmpty `json:"givenName,omitempty"`
	ImAddresses                     *[]string            `json:"imAddresses,omitempty"`
	Interests                       *[]string            `json:"interests,omitempty"`
	IsManagementRestricted          *bool                `json:"isManagementRestricted,omitempty"`
	IsResourceAccount               *bool                `json:"isResourceAccount,omitempty"`
	JobTitle                        *StringNullWhenEmpty `json:"jobTitle,omitempty"`
	Mail                            *string              `json:"mail,omitempty"`
	MailNickname                    *string              `json:"mailNickname,omitempty"`
	MobilePhone                     *StringNullWhenEmpty `json:"mobilePhone,omitempty"`
	MySite                          *string              `json:"mySite,omitempty"`
	OfficeLocation                  *StringNullWhenEmpty `json:"officeLocation,omitempty"`
	OnPremisesDistinguishedName     *string              `json:"onPremisesDistinguishedName,omitempty"`
	OnPremisesDomainName            *string              `json:"onPremisesDomainName,omitempty"`
	OnPremisesImmutableId           *string              `json:"onPremisesImmutableId,omitempty"`
	OnPremisesLastSyncDateTime      *string              `json:"onPremisesLastSyncDateTime,omitempty"`
	OnPremisesSamAccountName        *string              `json:"onPremisesSamAccountName,omitempty"`
	OnPremisesSecurityIdentifier    *string              `json:"onPremisesSecurityIdentifier,omitempty"`
	OnPremisesSyncEnabled           *bool                `json:"onPremisesSyncEnabled,omitempty"`
	OnPremisesUserPrincipalName     *string              `json:"onPremisesUserPrincipalName,omitempty"`
	OtherMails                      *[]string            `json:"otherMails,omitempty"`
	PasswordPolicies                *string              `json:"passwordPolicies,omitempty"`
	PastProjects                    *[]string            `json:"pastProjects,omitempty"`
	PostalCode                      *StringNullWhenEmpty `json:"postalCode,omitempty"`
	PreferredDataLocation           *string              `json:"preferredDataLocation,omitempty"`
	PreferredLanguage               *string              `json:"preferredLanguage,omitempty"`
	PreferredName                   *string              `json:"preferredName,omitempty"`
	ProxyAddresses                  *[]string            `json:"proxyAddresses,omitempty"`
	RefreshTokensValidFromDateTime  *time.Time           `json:"refreshTokensValidFromDateTime,omitempty"`
	Responsibilities                *[]string            `json:"responsibilities,omitempty"`
	Schools                         *[]string            `json:"schools,omitempty"`
	ShowInAddressList               *bool                `json:"showInAddressList,omitempty"`
	SignInSessionsValidFromDateTime *time.Time           `json:"signInSessionsValidFromDateTime,omitempty"`
	Skills                          *[]string            `json:"skills,omitempty"`
	State                           *StringNullWhenEmpty `json:"state,omitempty"`
	StreetAddress                   *StringNullWhenEmpty `json:"streetAddress,omitempty"`
	Surname                         *StringNullWhenEmpty `json:"surname,omitempty"`
	UsageLocation                   *StringNullWhenEmpty `json:"usageLocation,omitempty"`
	UserPrincipalName               *string              `json:"userPrincipalName,omitempty"`
	UserType                        *string              `json:"userType,omitempty"`

	PasswordProfile *UserPasswordProfile `json:"passwordProfile,omitempty"`
}

type UserPasswordProfile struct {
	ForceChangePasswordNextSignIn        *bool   `json:"forceChangePasswordNextSignIn,omitempty"`
	ForceChangePasswordNextSignInWithMfa *bool   `json:"forceChangePasswordNextSignInWithMfa,omitempty"`
	Password                             *string `json:"password,omitempty"`
}

type VerifiedPublisher struct {
	AddedDateTime       *time.Time `json:"addedDateTime,omitempty"`
	DisplayName         *string    `json:"displayName,omitempty"`
	VerifiedPublisherId *string    `json:"verifiedPublisherId,omitempty"`
}
