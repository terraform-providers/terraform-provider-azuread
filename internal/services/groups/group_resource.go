package groups

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

const groupResourceName = "azuread_group"

func groupResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: groupResourceCreate,
		ReadContext:   groupResourceRead,
		UpdateContext: groupResourceUpdate,
		DeleteContext: groupResourceDelete,

		CustomizeDiff: groupResourceCustomizeDiff,

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
			"display_name": {
				Description:      "The display name for the group",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validate.NoEmptyStrings,
			},

			"description": {
				Description: "The description for the group",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"mail_enabled": {
				Description:  "Whether the group is a mail enabled, with a shared group mailbox. At least one of `mail_enabled` or `security_enabled` must be specified. A group can be mail enabled _and_ security enabled",
				Type:         schema.TypeBool,
				Optional:     true,
				AtLeastOneOf: []string{"mail_enabled", "security_enabled"},
			},

			"members": {
				Description: "A set of members who should be present in this group. Supported object types are Users, Groups or Service Principals",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validate.UUID,
				},
			},

			"owners": {
				Description: "A set of owners who own this group. Supported object types are Users or Service Principals",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validate.UUID,
				},
			},

			"prevent_duplicate_names": {
				Description: "If `true`, will return an error if an existing group is found with the same name",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},

			"security_enabled": {
				Description:  "Whether the group is a security group for controlling access to in-app resources. At least one of `security_enabled` or `mail_enabled` must be specified. A group can be security enabled _and_ mail enabled",
				Type:         schema.TypeBool,
				Optional:     true,
				AtLeastOneOf: []string{"mail_enabled", "security_enabled"},
			},

			"types": {
				Description: "A set of group types to configure for the group. The only supported type is `Unified`, which specifies a Microsoft 365 group. Required when `mail_enabled` is true",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(msgraph.GroupTypeUnified),
					}, false),
				},
			},

			"object_id": {
				Description: "The object ID of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func groupResourceCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	client := meta.(*clients.Client).Groups.GroupsClient
	oldDisplayName, newDisplayName := diff.GetChange("display_name")
	mailEnabled := diff.Get("mail_enabled").(bool)
	groupTypes := make([]msgraph.GroupType, 0)
	for _, v := range diff.Get("types").(*schema.Set).List() {
		groupTypes = append(groupTypes, msgraph.GroupType(v.(string)))
	}
	hasGroupType := func(value msgraph.GroupType) bool {
		for _, v := range groupTypes {
			if value == v {
				return true
			}
		}
		return false
	}

	if mailEnabled && !hasGroupType(msgraph.GroupTypeUnified) {
		return fmt.Errorf("`types` must contain %q for mail-enabled groups", msgraph.GroupTypeUnified)
	}

	if !mailEnabled && hasGroupType(msgraph.GroupTypeUnified) {
		return fmt.Errorf("`mail_enabled` must be true for unified groups")
	}

	if diff.Get("prevent_duplicate_names").(bool) &&
		(oldDisplayName.(string) == "" || oldDisplayName.(string) != newDisplayName.(string)) {
		result, err := groupFindByName(ctx, client, newDisplayName.(string))
		if err != nil {
			return fmt.Errorf("could not check for existing application(s): %+v", err)
		}
		if result != nil && len(*result) > 0 {
			for _, existingGroup := range *result {
				if existingGroup.ID == nil {
					return fmt.Errorf("API error: group returned with nil object ID during duplicate name check")
				}
				if diff.Id() == "" || diff.Id() == *existingGroup.ID {
					return tf.ImportAsDuplicateError("azuread_group", *existingGroup.ID, newDisplayName.(string))
				}
			}
		}
	}

	return nil
}

func groupResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient
	callerId := meta.(*clients.Client).Claims.ObjectId
	displayName := d.Get("display_name").(string)

	// Perform this check at apply time to catch any duplicate names created during the same apply
	if d.Get("prevent_duplicate_names").(bool) {
		result, err := groupFindByName(ctx, client, displayName)
		if err != nil {
			return tf.ErrorDiagPathF(err, "name", "Could not check for existing groups(s)")
		}
		if result != nil && len(*result) > 0 {
			existingGroup := (*result)[0]
			if existingGroup.ID == nil {
				return tf.ErrorDiagF(errors.New("API returned group with nil object ID during duplicate name check"), "Bad API response")
			}
			return tf.ImportAsDuplicateDiag("azuread_group", *existingGroup.ID, displayName)
		}
	}

	mailNickname, err := uuid.GenerateUUID()
	if err != nil {
		return tf.ErrorDiagF(err, "Failed to generate mailNickname")
	}

	groupTypes := make([]msgraph.GroupType, 0)
	for _, v := range d.Get("types").(*schema.Set).List() {
		groupTypes = append(groupTypes, msgraph.GroupType(v.(string)))
	}

	properties := msgraph.Group{
		Description:     utils.NullableString(d.Get("description").(string)),
		DisplayName:     utils.String(displayName),
		GroupTypes:      groupTypes,
		MailEnabled:     utils.Bool(d.Get("mail_enabled").(bool)),
		MailNickname:    utils.String(mailNickname),
		SecurityEnabled: utils.Bool(d.Get("security_enabled").(bool)),
	}

	// Add the caller as the group owner to prevent lock-out after creation
	properties.AppendOwner(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, callerId)
	removeInitialOwner := true

	group, _, err := client.Create(ctx, properties)
	if err != nil {
		return tf.ErrorDiagF(err, "Creating group %q", displayName)
	}

	if group.ID == nil {
		return tf.ErrorDiagF(errors.New("API returned group with nil object ID"), "Bad API Response")
	}

	d.SetId(*group.ID)

	// Configure owners after the group is created, so they can be set one-by-one
	if v, ok := d.GetOk("owners"); ok {
		owners := v.(*schema.Set).List()
		for _, o := range owners {
			group.AppendOwner(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, o.(string))

			// If the authenticated principal is included in the owners list, make sure to not remove them after the fact
			if strings.EqualFold(callerId, o.(string)) {
				removeInitialOwner = false
			}
		}
		if _, err := client.AddOwners(ctx, group); err != nil {
			return tf.ErrorDiagF(err, "Could not add owners to group with ID: %q", d.Id())
		}
	}

	// Configure members after the group is created, so they can be reliably batched
	if v, ok := d.GetOk("members"); ok {
		members := v.(*schema.Set).List()
		for _, o := range members {
			group.AppendMember(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, o.(string))
		}
		if _, err := client.AddMembers(ctx, group); err != nil {
			return tf.ErrorDiagF(err, "Could not add members to group with ID: %q", d.Id())
		}
	}

	// Remove the initial owner
	if removeInitialOwner {
		ownersToRemove := []string{callerId}
		if _, err := client.RemoveOwners(ctx, *group.ID, &ownersToRemove); err != nil {
			return tf.ErrorDiagF(err, "Could not remove temporary owner of group with ID: %q", d.Id())
		}
	}

	return groupResourceRead(ctx, d, meta)
}

func groupResourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient
	groupId := d.Id()
	displayName := d.Get("display_name").(string)

	tf.LockByName(groupResourceName, groupId)
	defer tf.UnlockByName(groupResourceName, groupId)

	// Perform this check at apply time to catch any duplicate names created during the same apply
	if d.Get("prevent_duplicate_names").(bool) {
		result, err := groupFindByName(ctx, client, displayName)
		if err != nil {
			return tf.ErrorDiagPathF(err, "display_name", "Could not check for existing group(s)")
		}
		if result != nil && len(*result) > 0 {
			for _, existingGroup := range *result {
				if existingGroup.ID == nil {
					return tf.ErrorDiagF(errors.New("API returned group with nil object ID during duplicate name check"), "Bad API response")
				}

				if *existingGroup.ID != groupId {
					return tf.ImportAsDuplicateDiag("azuread_group", *existingGroup.ID, displayName)
				}
			}
		}
	}

	group := msgraph.Group{
		ID:              utils.String(groupId),
		Description:     utils.NullableString(d.Get("description").(string)),
		DisplayName:     utils.String(displayName),
		MailEnabled:     utils.Bool(d.Get("mail_enabled").(bool)),
		SecurityEnabled: utils.Bool(d.Get("security_enabled").(bool)),
	}

	if _, err := client.Update(ctx, group); err != nil {
		return tf.ErrorDiagF(err, "Updating group with ID: %q", d.Id())
	}

	if v, ok := d.GetOk("members"); ok && d.HasChange("members") {
		members, _, err := client.ListMembers(ctx, *group.ID)
		if err != nil {
			return tf.ErrorDiagF(err, "Could not retrieve members for group with ID: %q", d.Id())
		}

		existingMembers := *members
		desiredMembers := *tf.ExpandStringSlicePtr(v.(*schema.Set).List())
		membersForRemoval := utils.Difference(existingMembers, desiredMembers)
		membersToAdd := utils.Difference(desiredMembers, existingMembers)

		if membersForRemoval != nil {
			if _, err = client.RemoveMembers(ctx, d.Id(), &membersForRemoval); err != nil {
				return tf.ErrorDiagF(err, "Could not remove members from group with ID: %q", d.Id())
			}
		}

		if membersToAdd != nil {
			for _, m := range membersToAdd {
				group.AppendMember(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, m)
			}

			if _, err := client.AddMembers(ctx, &group); err != nil {
				return tf.ErrorDiagF(err, "Could not add members to group with ID: %q", d.Id())
			}
		}
	}

	if v, ok := d.GetOk("owners"); ok && d.HasChange("owners") {
		owners, _, err := client.ListOwners(ctx, *group.ID)
		if err != nil {
			return tf.ErrorDiagF(err, "Could not retrieve owners for group with ID: %q", d.Id())
		}

		existingOwners := *owners
		desiredOwners := *tf.ExpandStringSlicePtr(v.(*schema.Set).List())
		ownersForRemoval := utils.Difference(existingOwners, desiredOwners)
		ownersToAdd := utils.Difference(desiredOwners, existingOwners)

		if ownersToAdd != nil {
			for _, m := range ownersToAdd {
				group.AppendOwner(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, m)
			}

			if _, err := client.AddOwners(ctx, &group); err != nil {
				return tf.ErrorDiagF(err, "Could not add owners to group with ID: %q", d.Id())
			}
		}

		if ownersForRemoval != nil {
			if _, err = client.RemoveOwners(ctx, d.Id(), &ownersForRemoval); err != nil {
				return tf.ErrorDiagF(err, "Could not remove owners from group with ID: %q", d.Id())
			}
		}
	}

	return groupResourceRead(ctx, d, meta)
}

func groupResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient

	group, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			log.Printf("[DEBUG] Group with ID %q was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return tf.ErrorDiagF(err, "Retrieving group with object ID: %q", d.Id())
	}

	tf.Set(d, "description", group.Description)
	tf.Set(d, "display_name", group.DisplayName)
	tf.Set(d, "mail_enabled", group.MailEnabled)
	tf.Set(d, "object_id", group.ID)
	tf.Set(d, "security_enabled", group.SecurityEnabled)
	tf.Set(d, "types", group.GroupTypes)

	owners, _, err := client.ListOwners(ctx, *group.ID)
	if err != nil {
		return tf.ErrorDiagPathF(err, "owners", "Could not retrieve owners for group with object ID %q", d.Id())
	}
	tf.Set(d, "owners", owners)

	members, _, err := client.ListMembers(ctx, *group.ID)
	if err != nil {
		return tf.ErrorDiagPathF(err, "owners", "Could not retrieve members for group with object ID %q", d.Id())
	}
	tf.Set(d, "members", members)

	preventDuplicates := false
	if v := d.Get("prevent_duplicate_names").(bool); v {
		preventDuplicates = v
	}
	tf.Set(d, "prevent_duplicate_names", preventDuplicates)

	return nil
}

func groupResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient

	_, status, err := client.Get(ctx, d.Id())
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(fmt.Errorf("Group was not found"), "id", "Retrieving group with object ID %q", d.Id())
		}
		return tf.ErrorDiagPathF(err, "id", "Retrieving group with object ID: %q", d.Id())
	}

	if _, err := client.Delete(ctx, d.Id()); err != nil {
		return tf.ErrorDiagF(err, "Deleting group with object ID: %q", d.Id())
	}

	return nil
}
