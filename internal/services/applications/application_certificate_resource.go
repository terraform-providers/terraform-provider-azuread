package applications

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/helpers"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/applications/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func applicationCertificateResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: applicationCertificateResourceCreate,
		ReadContext:   applicationCertificateResourceRead,
		DeleteContext: applicationCertificateResourceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CertificateID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"application_object_id": {
				Description:      "The object ID of the application for which this certificate should be created",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"encoding": {
				Description: "Specifies the encoding used for the supplied certificate data",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "pem",
				ValidateFunc: validation.StringInSlice([]string{
					"base64",
					"hex",
					"pem",
				}, false),
			},

			"key_id": {
				Description:      "A UUID used to uniquely identify this certificate. If omitted, a random UUID will be automatically generated",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"start_date": {
				Description:  "The start date from which the certificate is valid, formatted as an RFC3339 date string (e.g. `2018-01-01T01:02:03Z`). If this isn't specified, the current date and time are use",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"end_date": {
				Description:   "The end date until which the certificate is valid, formatted as an RFC3339 date string (e.g. `2018-01-01T01:02:03Z`). If omitted, the API will decide a suitable expiry date, which is typically around 2 years from the start date",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"end_date_relative"},
				ValidateFunc:  validation.IsRFC3339Time,
			},

			"end_date_relative": {
				Description:      "A relative duration for which the certificate is valid until, for example `240h` (10 days) or `2400h30m`",
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ConflictsWith:    []string{"end_date"},
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"type": {
				Description: "The type of key/certificate",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"AsymmetricX509Cert",
					"Symmetric",
				}, false),
			},

			"value": {
				Description: "The certificate data, which can be PEM encoded, base64 encoded DER or hexadecimal encoded DER. See also the `encoding` argumen",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
		},
	}
}

func applicationCertificateResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient
	objectId := d.Get("application_object_id").(string)

	credential, err := helpers.KeyCredentialForResource(d)
	if err != nil {
		attr := ""
		if kerr, ok := err.(helpers.CredentialError); ok {
			attr = kerr.Attr()
		}
		return tf.ErrorDiagPathF(err, attr, "Generating certificate credentials for application with object ID %q", objectId)
	}

	if credential.KeyId == nil {
		return tf.ErrorDiagF(errors.New("keyId for certificate credential is nil"), "Creating certificate credential")
	}
	id := parse.NewCredentialID(objectId, "certificate", *credential.KeyId)

	tf.LockByName(applicationResourceName, id.ObjectId)
	defer tf.UnlockByName(applicationResourceName, id.ObjectId)

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(nil, "application_object_id", "Application with object ID %q was not found", id.ObjectId)
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving application with object ID %q", id.ObjectId)
	}

	newCredentials := make([]msgraph.KeyCredential, 0)
	if app.KeyCredentials != nil {
		for _, cred := range *app.KeyCredentials {
			if cred.KeyId != nil && strings.EqualFold(*cred.KeyId, *credential.KeyId) {
				return tf.ImportAsExistsDiag("azuread_application_certificate", id.String())
			}
			newCredentials = append(newCredentials, cred)
		}
	}

	newCredentials = append(newCredentials, *credential)

	properties := msgraph.Application{
		ID:             &id.ObjectId,
		KeyCredentials: &newCredentials,
	}
	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Adding certificate for application with object ID %q", id.ObjectId)
	}

	d.SetId(id.String())

	return applicationCertificateResourceRead(ctx, d, meta)
}

func applicationCertificateResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing certificate credential with ID %q", d.Id())
	}

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Application with ID %q for %s credential %q was not found - removing from state!", id.ObjectId, id.KeyType, id.KeyId)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving Application with object ID %q", id.ObjectId)
	}

	var credential *msgraph.KeyCredential
	if app.KeyCredentials != nil {
		for _, cred := range *app.KeyCredentials {
			if cred.KeyId != nil && strings.EqualFold(*cred.KeyId, id.KeyId) {
				credential = &cred
				break
			}
		}
	}

	if credential == nil {
		log.Printf("[DEBUG] Certificate credential %q (ID %q) was not found - removing from state!", id.KeyId, id.ObjectId)
		d.SetId("")
		return nil
	}

	tf.Set(d, "application_object_id", id.ObjectId)
	tf.Set(d, "key_id", id.KeyId)
	tf.Set(d, "type", string(credential.Type))

	startDate := ""
	if v := credential.StartDateTime; v != nil {
		startDate = v.Format(time.RFC3339)
	}
	tf.Set(d, "start_date", startDate)

	endDate := ""
	if v := credential.EndDateTime; v != nil {
		endDate = v.Format(time.RFC3339)
	}
	tf.Set(d, "end_date", endDate)

	return nil
}

func applicationCertificateResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing certificate credential with ID %q", d.Id())
	}

	tf.LockByName(applicationResourceName, id.ObjectId)
	defer tf.UnlockByName(applicationResourceName, id.ObjectId)

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(fmt.Errorf("Application was not found"), "application_object_id", "Retrieving Application with ID %q", id.ObjectId)
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving application with object ID %q", id.ObjectId)
	}

	newCredentials := make([]msgraph.KeyCredential, 0)
	if app.KeyCredentials != nil {
		for _, cred := range *app.KeyCredentials {
			if cred.KeyId != nil && !strings.EqualFold(*cred.KeyId, id.KeyId) {
				newCredentials = append(newCredentials, cred)
			}
		}
	}

	properties := msgraph.Application{
		ID:             &id.ObjectId,
		KeyCredentials: &newCredentials,
	}
	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Removing certificate credential %q from application with object ID %q", id.KeyId, id.ObjectId)
	}

	return nil
}
