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
	"github.com/manicminer/hamilton/msgraph"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/applications/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/utils"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func applicationPreAuthorizedResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: applicationPreAuthorizedResourceCreate,
		ReadContext:   applicationPreAuthorizedResourceRead,
		UpdateContext: applicationPreAuthorizedResourceUpdate,
		DeleteContext: applicationPreAuthorizedResourceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ApplicationPreAuthorizedID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"application_object_id": {
				Description:      "The object ID of the application to which this pre-authorized application should be added",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"authorized_app_id": {
				Description:      "The application ID of the pre-authorized application",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"permission_ids": {
				Description: "The IDs of the permission scopes required by the pre-authorized application",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validate.UUID,
				},
			},
		},
	}
}

func applicationPreAuthorizedResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient
	id := parse.NewApplicationPreAuthorizedID(d.Get("application_object_id").(string), d.Get("authorized_app_id").(string))

	tf.LockByName(applicationResourceName, id.ObjectId)
	defer tf.UnlockByName(applicationResourceName, id.ObjectId)

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(nil, "application_object_id", "Application with object ID %q was not found", id.ObjectId)
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving application with object ID %q", id.ObjectId)
	}
	if app == nil || app.ID == nil {
		return tf.ErrorDiagF(errors.New("nil application or application with nil ID was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}

	newPreAuthorizedApps := make([]msgraph.ApiPreAuthorizedApplication, 0)
	if app.Api != nil && app.Api.PreAuthorizedApplications != nil {
		for _, a := range *app.Api.PreAuthorizedApplications {
			if a.AppId != nil && strings.EqualFold(*a.AppId, id.AppId) {
				return tf.ImportAsExistsDiag("azuread_application_pre_authorized", id.String())
			}
			newPreAuthorizedApps = append(newPreAuthorizedApps, a)
		}
	}

	newPreAuthorizedApps = append(newPreAuthorizedApps, msgraph.ApiPreAuthorizedApplication{
		AppId:         utils.String(id.AppId),
		PermissionIds: tf.ExpandStringSlicePtr(d.Get("permission_ids").(*schema.Set).List()),
	})

	properties := msgraph.Application{
		ID: app.ID,
		Api: &msgraph.ApplicationApi{
			PreAuthorizedApplications: &newPreAuthorizedApps,
		},
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Adding pre-authorized application %q for application with object ID %q", id.AppId, id.ObjectId)
	}

	d.SetId(id.String())

	return applicationPreAuthorizedResourceRead(ctx, d, meta)
}

func applicationPreAuthorizedResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient
	id, err := parse.ApplicationPreAuthorizedID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing pre-authorized application ID %q", d.Id())
	}

	tf.LockByName(applicationResourceName, id.ObjectId)
	defer tf.UnlockByName(applicationResourceName, id.ObjectId)

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(nil, "application_object_id", "Application with object ID %q was not found", id.ObjectId)
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving application with object ID %q", id.ObjectId)
	}
	if app == nil || app.ID == nil {
		return tf.ErrorDiagF(errors.New("nil application or application with nil ID was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}
	if app.Api == nil || app.Api.PreAuthorizedApplications == nil {
		return tf.ErrorDiagF(errors.New("application with nil preAuthorizedApplications was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}

	found := false
	newPreAuthorizedApps := *app.Api.PreAuthorizedApplications
	for i, a := range newPreAuthorizedApps {
		if a.AppId != nil && strings.EqualFold(*a.AppId, id.AppId) {
			found = true
			newPreAuthorizedApps[i].PermissionIds = tf.ExpandStringSlicePtr(d.Get("permission_ids").(*schema.Set).List())
			break
		}
	}
	if !found {
		return tf.ErrorDiagF(fmt.Errorf("could not match an existing preAuthorizedApplication for %q", id.AppId), "retrieving application with object ID %q", id.ObjectId)
	}

	properties := msgraph.Application{
		ID: app.ID,
		Api: &msgraph.ApplicationApi{
			PreAuthorizedApplications: &newPreAuthorizedApps,
		},
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Updating pre-authorized application %q for application with object ID %q", id.AppId, id.ObjectId)
	}

	return applicationPreAuthorizedResourceRead(ctx, d, meta)
}

func applicationPreAuthorizedResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient
	id, err := parse.ApplicationPreAuthorizedID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing pre-authorized application ID %q", d.Id())
	}

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Application with ID %q for pre-authorized application %q was not found - removing from state!", id.ObjectId, id.AppId)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving Application with object ID %q", id.ObjectId)
	}
	if app == nil || app.ID == nil {
		return tf.ErrorDiagF(errors.New("nil application or application with nil ID was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}
	if app.Api == nil || app.Api.PreAuthorizedApplications == nil {
		return tf.ErrorDiagF(errors.New("application with nil preAuthorizedApplications was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}

	var preAuthorizedApp *msgraph.ApiPreAuthorizedApplication
	for _, a := range *app.Api.PreAuthorizedApplications {
		if a.AppId != nil && strings.EqualFold(*a.AppId, id.AppId) {
			preAuthorizedApp = &a
			break
		}
	}
	if preAuthorizedApp == nil {
		log.Printf("[DEBUG] No matching preAuthorizedApplication for ID %q - removing from state!", id)
		d.SetId("")
		return nil
	}

	d.Set("application_object_id", id.ObjectId)
	d.Set("authorized_app_id", id.AppId)
	d.Set("permission_ids", tf.FlattenStringSlicePtr(preAuthorizedApp.PermissionIds))

	return nil
}

func applicationPreAuthorizedResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Applications.ApplicationsClient
	id, err := parse.ApplicationPreAuthorizedID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing pre-authorized application ID %q", d.Id())
	}

	app, status, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Application with ID %q for pre-authorized application %q was not found - removing from state!", id.ObjectId, id.AppId)
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagPathF(err, "application_object_id", "Retrieving Application with object ID %q", id.ObjectId)
	}
	if app == nil || app.ID == nil {
		return tf.ErrorDiagF(errors.New("nil application or application with nil ID was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}
	if app.Api == nil || app.Api.PreAuthorizedApplications == nil {
		return tf.ErrorDiagF(errors.New("application with nil preAuthorizedApplications was returned"), "API error retrieving application with object ID %q", id.ObjectId)
	}

	newPreAuthorizedApps := make([]msgraph.ApiPreAuthorizedApplication, 0)
	for _, a := range *app.Api.PreAuthorizedApplications {
		if a.AppId != nil && !strings.EqualFold(*a.AppId, id.AppId) {
			newPreAuthorizedApps = append(newPreAuthorizedApps, a)
			break
		}
	}

	properties := msgraph.Application{
		ID: app.ID,
		Api: &msgraph.ApplicationApi{
			PreAuthorizedApplications: &newPreAuthorizedApps,
		},
	}

	if _, err := client.Update(ctx, properties); err != nil {
		return tf.ErrorDiagF(err, "Removing pre-authorized application %q from application with object ID %q", id.AppId, id.ObjectId)
	}

	return nil
}
