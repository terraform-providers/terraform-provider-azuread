package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func userResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: userResourceCreate,
		ReadContext:   userResourceRead,
		UpdateContext: userResourceUpdate,
		DeleteContext: userResourceDelete,

		CustomizeDiff: userResourceCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			if _, err := uuid.ParseUUID(id); err != nil {
				return fmt.Errorf("specified ID (%q) is not valid: %s", id, err)
			}
			return nil
		}),

		Schema: map[string]*schema.Schema{
			"user_principal_name": {
				Description:      "The user principal name (UPN) of the user",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.StringIsEmailAddress,
			},

			"display_name": {
				Description:      "The name to display in the address book for the user",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"account_enabled": {
				Description: "Whether or not the account should be enabled",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},

			"age_group": {
				Description: "The age group of the user",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					string(msgraph.AgeGroupNone),
					string(msgraph.AgeGroupAdult),
					string(msgraph.AgeGroupMinor),
					string(msgraph.AgeGroupNotAdult),
				}, false),
			},

			"business_phones": {
				Description: "The telephone numbers for the user. Only one number can be set for this property. Read-only for users synced with Azure AD Connect",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"city": {
				Description: "The city in which the user is located",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"company_name": {
				Description: "The company name which the user is associated. This property can be useful for describing the company that an external user comes from",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"consent_provided_for_minor": {
				Description: "Whether consent has been obtained for minors",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					string(msgraph.ConsentProvidedForMinorNone),
					string(msgraph.ConsentProvidedForMinorDenied),
					string(msgraph.ConsentProvidedForMinorGranted),
					string(msgraph.ConsentProvidedForMinorNotRequired),
				}, false),
			},

			"country": {
				Description: "The country/region in which the user is located, e.g. `US` or `UK`",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"department": {
				Description: "The name for the department in which the user works",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"employee_id": {
				Description:  "The employee identifier assigned to the user by the organisation",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 16),
			},

			"force_password_change": {
				Description: "Whether the user is forced to change the password during the next sign-in. Only takes effect when also changing the password",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},

			"given_name": {
				Description: "The given name (first name) of the user",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"fax_number": {
				Description: "The fax number of the user",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"job_title": {
				Description: "The user’s job title",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"mail": {
				Description: "The SMTP address for the user. Cannot be unset.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"mail_nickname": {
				Description: "The mail alias for the user. Defaults to the user name part of the user principal name (UPN)",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"mobile_phone": {
				Description: "The primary cellular telephone number for the user",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"office_location": {
				Description: "The office location in the user's place of business",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"onpremises_immutable_id": {
				Description: "The value used to associate an on-premise Active Directory user account with their Azure AD user object. This must be specified if you are using a federated domain for the user's `user_principal_name` property when creating a new user account",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"other_mails": {
				Description: "Additional email addresses for the user",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"password": {
				Description:  "The password for the user. The password must satisfy minimum requirements as specified by the password policy. The maximum length is 256 characters. This property is required when creating a new user",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 256), // Currently the max length for AAD passwords is 256
			},

			"postal_code": {
				Description: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"preferred_language": {
				Description:      "The user's preferred language, in ISO 639-1 notation",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validate.ISO639Language,
			},

			"show_in_address_list": {
				Description: "Whether or not the Outlook global address list should include this user",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},

			"state": {
				Description: "The state or province in the user's address",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"street_address": {
				Description: "The street address of the user's place of business",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"surname": {
				Description: "The user's surname (family name or last name)",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"usage_location": {
				Description: "The usage location of the user. Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries. The usage location is a two letter country code (ISO standard 3166). Examples include: `NO`, `JP`, and `GB`. Cannot be reset to null once set",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"about_me": {
				Description: "A freeform field for the user to describe themselves",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"object_id": {
				Description: "The object ID of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"creation_type": {
				Description: "Indicates whether the user account was created as a regular school or work account (`null`), an external account (`Invitation`), a local account for an Azure Active Directory B2C tenant (`LocalAccount`) or self-service sign-up using email verification (`EmailVerified`)",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"external_user_state": {
				Description: "For an external user invited to the tenant, this property represents the invited user's invitation status",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"im_addresses": {
				Description: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"onpremises_distinguished_name": {
				Description: "The on-premise Active Directory distinguished name (DN) of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"onpremises_domain_name": {
				Description: "The on-premise FQDN (i.e. dnsDomainName) of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"onpremises_sam_account_name": {
				Description: "The on-premise SAM account name of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"onpremises_security_identifier": {
				Description: "The on-premise security identifier (SID) of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"onpremises_sync_enabled": {
				Description: "Whether this user is synchronized from an on-premises directory (true), no longer synchronized (false), or has never been synchronized (null)",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"onpremises_user_principal_name": {
				Description: "The on-premise user principal name of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"proxy_addresses": {
				Description: "Email addresses for the user that direct to the same mailbox",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"user_type": {
				Description: "The user type in the directory. Possible values are `Guest` or `Member`",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func userResourceCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if diff.Id() == "" && diff.Get("password").(string) == "" {
		return fmt.Errorf("`password` is required when creating a new user")
	}
	return nil
}

func userResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	upn := d.Get("user_principal_name").(string)
	mailNickName := d.Get("mail_nickname").(string)

	// Default mail nickname to the first part of the UPN (matches the portal)
	if mailNickName == "" {
		mailNickName = strings.Split(upn, "@")[0]
	}

	properties := msgraph.User{
		AccountEnabled:          utils.Bool(d.Get("account_enabled").(bool)),
		AgeGroup:                utils.NullableString(d.Get("age_group").(string)),
		City:                    utils.NullableString(d.Get("city").(string)),
		ConsentProvidedForMinor: utils.NullableString(d.Get("consent_provided_for_minor").(string)),
		CompanyName:             utils.NullableString(d.Get("company_name").(string)),
		Country:                 utils.NullableString(d.Get("country").(string)),
		Department:              utils.NullableString(d.Get("department").(string)),
		DisplayName:             utils.String(d.Get("display_name").(string)),
		EmployeeId:              utils.NullableString(d.Get("employee_id").(string)),
		FaxNumber:               utils.NullableString(d.Get("fax_number").(string)),
		GivenName:               utils.NullableString(d.Get("given_name").(string)),
		JobTitle:                utils.NullableString(d.Get("job_title").(string)),
		Mail:                    utils.NullableString(d.Get("mail").(string)),
		MailNickname:            utils.String(mailNickName),
		MobilePhone:             utils.NullableString(d.Get("mobile_phone").(string)),
		OfficeLocation:          utils.NullableString(d.Get("office_location").(string)),
		OtherMails:              tf.ExpandStringSlicePtr(d.Get("other_mails").(*schema.Set).List()),
		PostalCode:              utils.NullableString(d.Get("postal_code").(string)),
		PreferredLanguage:       utils.NullableString(d.Get("preferred_language").(string)),
		ShowInAddressList:       utils.Bool(d.Get("show_in_address_list").(bool)),
		State:                   utils.NullableString(d.Get("state").(string)),
		StreetAddress:           utils.NullableString(d.Get("street_address").(string)),
		Surname:                 utils.NullableString(d.Get("surname").(string)),
		UsageLocation:           utils.NullableString(d.Get("usage_location").(string)),
		UserPrincipalName:       utils.String(upn),

		PasswordProfile: &msgraph.UserPasswordProfile{
			ForceChangePasswordNextSignIn: utils.Bool(d.Get("force_password_change").(bool)),
			Password:                      utils.String(d.Get("password").(string)),
		},
	}

	if v, ok := d.GetOk("business_phones"); ok {
		properties.BusinessPhones = tf.ExpandStringSlicePtr(v.([]interface{}))
	}

	if v, ok := d.GetOk("onpremises_immutable_id"); ok {
		properties.OnPremisesImmutableId = utils.String(v.(string))
	}

	user, _, err := client.Create(ctx, properties)
	if err != nil {
		return tf.ErrorDiagF(err, "Creating user %q", upn)
	}

	if user.ID == nil || *user.ID == "" {
		return tf.ErrorDiagF(errors.New("API returned group with nil object ID"), "Bad API Response")
	}

	d.SetId(*user.ID)

	return userResourceRead(ctx, d, meta)
}

func userResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	properties := msgraph.User{
		ID:                      utils.String(d.Id()),
		AccountEnabled:          utils.Bool(d.Get("account_enabled").(bool)),
		AgeGroup:                utils.NullableString(d.Get("age_group").(string)),
		City:                    utils.NullableString(d.Get("city").(string)),
		CompanyName:             utils.NullableString(d.Get("company_name").(string)),
		ConsentProvidedForMinor: utils.NullableString(d.Get("consent_provided_for_minor").(string)),
		Country:                 utils.NullableString(d.Get("country").(string)),
		Department:              utils.NullableString(d.Get("department").(string)),
		DisplayName:             utils.String(d.Get("display_name").(string)),
		EmployeeId:              utils.NullableString(d.Get("employee_id").(string)),
		FaxNumber:               utils.NullableString(d.Get("fax_number").(string)),
		GivenName:               utils.NullableString(d.Get("given_name").(string)),
		JobTitle:                utils.NullableString(d.Get("job_title").(string)),
		MailNickname:            utils.String(d.Get("mail_nickname").(string)),
		MobilePhone:             utils.NullableString(d.Get("mobile_phone").(string)),
		OfficeLocation:          utils.NullableString(d.Get("office_location").(string)),
		OtherMails:              tf.ExpandStringSlicePtr(d.Get("other_mails").(*schema.Set).List()),
		PostalCode:              utils.NullableString(d.Get("postal_code").(string)),
		PreferredLanguage:       utils.NullableString(d.Get("preferred_language").(string)),
		ShowInAddressList:       utils.Bool(d.Get("show_in_address_list").(bool)),
		State:                   utils.NullableString(d.Get("state").(string)),
		StreetAddress:           utils.NullableString(d.Get("street_address").(string)),
		Surname:                 utils.NullableString(d.Get("surname").(string)),
		UsageLocation:           utils.NullableString(d.Get("usage_location").(string)),
	}

	if d.HasChange("password") {
		properties.PasswordProfile = &msgraph.UserPasswordProfile{
			ForceChangePasswordNextSignIn: utils.Bool(d.Get("force_password_change").(bool)),
			Password:                      utils.String(d.Get("password").(string)),
		}
	}

	if d.HasChange("business_phones") {
		properties.BusinessPhones = tf.ExpandStringSlicePtr(d.Get("business_phones").([]interface{}))
	}

	if d.HasChange("mail") {
		if mail := d.Get("mail").(string); mail != "" {
			properties.Mail = utils.NullableString(mail)
		}
	}

	if d.HasChange("onpremises_immutable_id") {
		properties.OnPremisesImmutableId = utils.String(d.Get("onpremises_immutable_id").(string))
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Could not update user with ID: %q", d.Id())
	}

	return userResourceRead(ctx, d, meta)
}

func userResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	objectId := d.Id()

	user, status, err := client.Get(ctx, objectId)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] User with Object ID %q was not found - removing from state!", objectId)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagF(err, "Retrieving user with object ID: %q", objectId)
	}

	tf.Set(d, "about_me", user.AboutMe)
	tf.Set(d, "account_enabled", user.AccountEnabled)
	tf.Set(d, "age_group", user.AgeGroup)
	tf.Set(d, "business_phones", user.BusinessPhones)
	tf.Set(d, "city", user.City)
	tf.Set(d, "company_name", user.CompanyName)
	tf.Set(d, "consent_provided_for_minor", user.ConsentProvidedForMinor)
	tf.Set(d, "country", user.Country)
	tf.Set(d, "creation_type", user.CreationType)
	tf.Set(d, "department", user.Department)
	tf.Set(d, "display_name", user.DisplayName)
	tf.Set(d, "employee_id", user.EmployeeId)
	tf.Set(d, "external_user_state", user.ExternalUserState)
	tf.Set(d, "fax_number", user.FaxNumber)
	tf.Set(d, "given_name", user.GivenName)
	tf.Set(d, "im_addresses", user.ImAddresses)
	tf.Set(d, "job_title", user.JobTitle)
	tf.Set(d, "mail", user.Mail)
	tf.Set(d, "mail_nickname", user.MailNickname)
	tf.Set(d, "mobile_phone", user.MobilePhone)
	tf.Set(d, "object_id", user.ID)
	tf.Set(d, "office_location", user.OfficeLocation)
	tf.Set(d, "onpremises_distinguished_name", user.OnPremisesDistinguishedName)
	tf.Set(d, "onpremises_domain_name", user.OnPremisesDomainName)
	tf.Set(d, "onpremises_immutable_id", user.OnPremisesImmutableId)
	tf.Set(d, "onpremises_sam_account_name", user.OnPremisesSamAccountName)
	tf.Set(d, "onpremises_security_identifier", user.OnPremisesSecurityIdentifier)
	tf.Set(d, "onpremises_sync_enabled", user.OnPremisesSyncEnabled)
	tf.Set(d, "onpremises_user_principal_name", user.OnPremisesUserPrincipalName)
	tf.Set(d, "other_mails", user.OtherMails)
	tf.Set(d, "postal_code", user.PostalCode)
	tf.Set(d, "preferred_language", user.PreferredLanguage)
	tf.Set(d, "proxy_addresses", user.ProxyAddresses)
	tf.Set(d, "show_in_address_list", user.ShowInAddressList)
	tf.Set(d, "state", user.State)
	tf.Set(d, "street_address", user.StreetAddress)
	tf.Set(d, "surname", user.Surname)
	tf.Set(d, "usage_location", user.UsageLocation)
	tf.Set(d, "user_principal_name", user.UserPrincipalName)
	tf.Set(d, "user_type", user.UserType)

	return nil
}

func userResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	_, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(fmt.Errorf("User was not found"), "id", "Retrieving user with object ID %q", d.Id())
		}

		return tf.ErrorDiagPathF(err, "id", "Retrieving user with object ID %q", d.Id())
	}

	status, err = client.Delete(ctx, d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Deleting user with object ID %q, got status %d", d.Id(), status)
	}

	return nil
}
