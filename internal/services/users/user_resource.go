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

			"job_title": {
				Description: "The user’s job title",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"mail": {
				Description: "The primary email address of the user",
				Type:        schema.TypeString,
				Optional:    true,
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

			"onpremises_sam_account_name": {
				Description: "The on-premise SAM account name of the user",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"onpremises_user_principal_name": {
				Description: "The on-premise user principal name of the user",
				Type:        schema.TypeString,
				Computed:    true,
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

			"street_address": {
				Description: "The street address of the user's place of business",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"state": {
				Description: "The state or province in the user's address",
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

			"object_id": {
				Description: "The object ID of the user",
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
		AccountEnabled:    utils.Bool(d.Get("account_enabled").(bool)),
		City:              utils.NullableString(d.Get("city").(string)),
		CompanyName:       utils.NullableString(d.Get("company_name").(string)),
		Country:           utils.NullableString(d.Get("country").(string)),
		Department:        utils.NullableString(d.Get("department").(string)),
		DisplayName:       utils.String(d.Get("display_name").(string)),
		GivenName:         utils.NullableString(d.Get("given_name").(string)),
		JobTitle:          utils.NullableString(d.Get("job_title").(string)),
		MailNickname:      utils.String(mailNickName),
		MobilePhone:       utils.NullableString(d.Get("mobile_phone").(string)),
		OfficeLocation:    utils.NullableString(d.Get("office_location").(string)),
		PostalCode:        utils.NullableString(d.Get("postal_code").(string)),
		State:             utils.NullableString(d.Get("state").(string)),
		StreetAddress:     utils.NullableString(d.Get("street_address").(string)),
		Surname:           utils.NullableString(d.Get("surname").(string)),
		UsageLocation:     utils.NullableString(d.Get("usage_location").(string)),
		UserPrincipalName: utils.String(upn),

		PasswordProfile: &msgraph.UserPasswordProfile{
			ForceChangePasswordNextSignIn: utils.Bool(d.Get("force_password_change").(bool)),
			Password:                      utils.String(d.Get("password").(string)),
		},
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
		ID:             utils.String(d.Id()),
		AccountEnabled: utils.Bool(d.Get("account_enabled").(bool)),
		City:           utils.NullableString(d.Get("city").(string)),
		CompanyName:    utils.NullableString(d.Get("company_name").(string)),
		Country:        utils.NullableString(d.Get("country").(string)),
		Department:     utils.NullableString(d.Get("department").(string)),
		DisplayName:    utils.String(d.Get("display_name").(string)),
		GivenName:      utils.NullableString(d.Get("given_name").(string)),
		JobTitle:       utils.NullableString(d.Get("job_title").(string)),
		MailNickname:   utils.String(d.Get("mail_nickname").(string)),
		MobilePhone:    utils.NullableString(d.Get("mobile_phone").(string)),
		OfficeLocation: utils.NullableString(d.Get("office_location").(string)),
		PostalCode:     utils.NullableString(d.Get("postal_code").(string)),
		State:          utils.NullableString(d.Get("state").(string)),
		StreetAddress:  utils.NullableString(d.Get("street_address").(string)),
		Surname:        utils.NullableString(d.Get("surname").(string)),
		UsageLocation:  utils.NullableString(d.Get("usage_location").(string)),
	}

	if d.HasChange("password") {
		properties.PasswordProfile = &msgraph.UserPasswordProfile{
			ForceChangePasswordNextSignIn: utils.Bool(d.Get("force_password_change").(bool)),
			Password:                      utils.String(d.Get("password").(string)),
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

	tf.Set(d, "account_enabled", user.AccountEnabled)
	tf.Set(d, "city", user.City)
	tf.Set(d, "company_name", user.CompanyName)
	tf.Set(d, "country", user.Country)
	tf.Set(d, "department", user.Department)
	tf.Set(d, "display_name", user.DisplayName)
	tf.Set(d, "given_name", user.GivenName)
	tf.Set(d, "job_title", user.JobTitle)
	tf.Set(d, "mail", user.Mail)
	tf.Set(d, "mail_nickname", user.MailNickname)
	tf.Set(d, "mobile_phone", user.MobilePhone)
	tf.Set(d, "object_id", user.ID)
	tf.Set(d, "office_location", user.OfficeLocation)
	tf.Set(d, "onpremises_immutable_id", user.OnPremisesImmutableId)
	tf.Set(d, "onpremises_sam_account_name", user.OnPremisesSamAccountName)
	tf.Set(d, "onpremises_user_principal_name", user.OnPremisesUserPrincipalName)
	tf.Set(d, "postal_code", user.PostalCode)
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
