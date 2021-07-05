package groups

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-azuread/internal/clients"
	"github.com/hashicorp/terraform-provider-azuread/internal/services/groups/parse"
	"github.com/hashicorp/terraform-provider-azuread/internal/tf"
	"github.com/hashicorp/terraform-provider-azuread/internal/validate"
)

func groupMemberResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: groupMemberResourceCreate,
		ReadContext:   groupMemberResourceRead,
		DeleteContext: groupMemberResourceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: tf.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.GroupMemberID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"group_object_id": {
				Description:      "The object ID of the group you want to add the member to",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},

			"member_object_id": {
				Description:      "The object ID of the principal you want to add as a member to the group. Supported object types are Users, Groups or Service Principals",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validate.UUID,
			},
		},
	}
}

func groupMemberResourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient
	groupId := d.Get("group_object_id").(string)
	memberId := d.Get("member_object_id").(string)

	id := parse.NewGroupMemberID(groupId, memberId)

	tf.LockByName(groupResourceName, id.GroupId)
	defer tf.UnlockByName(groupResourceName, id.GroupId)

	group, status, err := client.Get(ctx, groupId)
	if err != nil {
		if status == http.StatusNotFound {
			return tf.ErrorDiagPathF(nil, "object_id", "Group with object ID %q was not found", groupId)
		}
		return tf.ErrorDiagPathF(err, "object_id", "Retrieving group with object ID: %q", groupId)
	}

	existingMembers, _, err := client.ListMembers(ctx, id.GroupId)
	if err != nil {
		return tf.ErrorDiagF(err, "Listing existing members for group with object ID: %q", id.GroupId)
	}
	if existingMembers != nil {
		for _, v := range *existingMembers {
			if strings.EqualFold(v, memberId) {
				return tf.ImportAsExistsDiag("azuread_group_member", id.String())
			}
		}
	}

	group.AppendMember(client.BaseClient.Endpoint, client.BaseClient.ApiVersion, memberId)

	if _, err := client.AddMembers(ctx, group); err != nil {
		return tf.ErrorDiagF(err, "Adding group member %q to group %q", memberId, groupId)
	}

	d.SetId(id.String())
	return groupMemberResourceRead(ctx, d, meta)
}

func groupMemberResourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient

	id, err := parse.GroupMemberID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing Group Member ID %q", d.Id())
	}

	members, _, err := client.ListMembers(ctx, id.GroupId)
	if err != nil {
		return tf.ErrorDiagF(err, "Retrieving members for group with object ID: %q", id.GroupId)
	}

	var memberObjectId string
	if members != nil {
		for _, objectId := range *members {
			if strings.EqualFold(objectId, id.MemberId) {
				memberObjectId = objectId
				break
			}
		}
	}

	if memberObjectId == "" {
		log.Printf("[DEBUG] Member with ID %q was not found in Group %q - removing from state", id.MemberId, id.GroupId)
		d.SetId("")
		return nil
	}

	tf.Set(d, "group_object_id", id.GroupId)
	tf.Set(d, "member_object_id", memberObjectId)

	return nil
}

func groupMemberResourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client).Groups.GroupsClient

	id, err := parse.GroupMemberID(d.Id())
	if err != nil {
		return tf.ErrorDiagPathF(err, "id", "Parsing Group Member ID %q", d.Id())
	}

	tf.LockByName(groupResourceName, id.GroupId)
	defer tf.UnlockByName(groupResourceName, id.GroupId)

	if _, err := client.RemoveMembers(ctx, id.GroupId, &[]string{id.MemberId}); err != nil {
		return tf.ErrorDiagF(err, "Removing member %q from group with object ID: %q", id.MemberId, id.GroupId)
	}

	return nil
}
