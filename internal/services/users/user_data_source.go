package users

import (
	"context"
	"errors"
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

func userDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: userDataSourceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"mail_nickname": {
				Description:      "The email alias of the user",
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"mail_nickname", "object_id", "user_principal_name"},
				Computed:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"object_id": {
				Description:      "The object ID of the user",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"mail_nickname", "object_id", "user_principal_name"},
				ValidateDiagFunc: validate.UUID,
			},

			"user_principal_name": {
				Description:      "The user principal name (UPN) of the user",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"mail_nickname", "object_id", "user_principal_name"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"account_enabled": {
				Description: "Whether or not the account is enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"age_group": {
				Description: "The age group of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"business_phones": {
				Description: "The telephone numbers for the user",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"city": {
				Description: "The city in which the user is located",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"company_name": {
				Description: "The company name which the user is associated. This property can be useful for describing the company that an external user comes from",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"consent_provided_for_minor": {
				Description: "Whether consent has been obtained for minors",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"country": {
				Description: "The country/region in which the user is located, e.g. `US` or `UK`",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"creation_type": {
				Description: "Indicates whether the user account was created as a regular school or work account (`null`), an external account (`Invitation`), a local account for an Azure Active Directory B2C tenant (`LocalAccount`) or self-service sign-up using email verification (`EmailVerified`)",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"department": {
				Description: "The name for the department in which the user works",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"display_name": {
				Description: "The display name of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"employee_id": {
				Description: "The employee identifier assigned to the user by the organisation",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"external_user_state": {
				Description: "For an external user invited to the tenant, this property represents the invited user's invitation status",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"fax_number": {
				Description: "The fax number of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"given_name": {
				Description: "The given name (first name) of the user",
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

			"job_title": {
				Description: "The user’s job title",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"mail": {
				Description: "The SMTP address for the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"mobile_phone": {
				Description: "The primary cellular telephone number for the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"office_location": {
				Description: "The office location in the user's place of business",
				Type:        schema.TypeString,
				Computed:    true,
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

			"onpremises_immutable_id": {
				Description: "The value used to associate an on-premise Active Directory user account with their Azure AD user object",
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

			"other_mails": {
				Description: "Additional email addresses for the user",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"postal_code": {
				Description: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"preferred_language": {
				Description: "The user's preferred language, in ISO 639-1 notation",
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

			"show_in_address_list": {
				Description: "Whether or not the Outlook global address list should include this user",
				Type:        schema.TypeBool,
				Computed:    true,
			},

			"state": {
				Description: "The state or province in the user's address",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"street_address": {
				Description: "The street address of the user's place of business",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"surname": {
				Description: "The user's surname (family name or last name)",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"usage_location": {
				Description: "The usage location of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"user_type": {
				Description: "The user type in the directory. Possible values are `Guest` or `Member`",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func userDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Users.UsersClient

	var user msgraph.User

	if upn, ok := d.Get("user_principal_name").(string); ok && upn != "" {
		query := odata.Query{
			Filter: fmt.Sprintf("userPrincipalName eq '%s'", upn),
		}
		users, _, err := client.List(ctx, query)
		if err != nil {
			return tf.ErrorDiagF(err, "Finding user with UPN: %q", upn)
		}
		if users == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}
		count := len(*users)
		if count > 1 {
			return tf.ErrorDiagPathF(nil, "user_principal_name", "More than one user found with UPN: %q", upn)
		} else if count == 0 {
			return tf.ErrorDiagPathF(err, "user_principal_name", "User with UPN %q was not found", upn)
		}
		user = (*users)[0]
	} else if objectId, ok := d.Get("object_id").(string); ok && objectId != "" {
		u, status, err := client.Get(ctx, objectId)
		if err != nil {
			if status == http.StatusNotFound {
				return tf.ErrorDiagPathF(nil, "object_id", "User not found with object ID: %q", objectId)
			}
			return tf.ErrorDiagF(err, "Retrieving user with object ID: %q", objectId)
		}
		if u == nil {
			return tf.ErrorDiagPathF(nil, "object_id", "User not found with object ID: %q", objectId)
		}
		user = *u
	} else if mailNickname, ok := d.Get("mail_nickname").(string); ok && mailNickname != "" {
		query := odata.Query{
			Filter: fmt.Sprintf("mailNickname eq '%s'", mailNickname),
		}
		users, _, err := client.List(ctx, query)
		if err != nil {
			return tf.ErrorDiagF(err, "Finding user with email alias: %q", mailNickname)
		}
		if users == nil {
			return tf.ErrorDiagF(errors.New("API returned nil result"), "Bad API Response")
		}
		count := len(*users)
		if count > 1 {
			return tf.ErrorDiagPathF(nil, "mail_nickname", "More than one user found with email alias: %q", upn)
		} else if count == 0 {
			return tf.ErrorDiagPathF(err, "mail_nickname", "User not found with email alias: %q", upn)
		}
		user = (*users)[0]
	} else {
		return tf.ErrorDiagF(nil, "One of `object_id`, `user_principal_name` or `mail_nickname` must be supplied")
	}

	if user.ID == nil {
		return tf.ErrorDiagF(errors.New("API returned user with nil object ID"), "Bad API Response")
	}

	d.SetId(*user.ID)

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
