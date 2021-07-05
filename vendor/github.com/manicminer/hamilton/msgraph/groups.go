package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/manicminer/hamilton/odata"
)

// GroupsClient performs operations on Groups.
type GroupsClient struct {
	BaseClient Client
}

// NewGroupsClient returns a new GroupsClient.
func NewGroupsClient(tenantId string) *GroupsClient {
	return &GroupsClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// List returns a list of Groups, optionally filtered using OData.
func (c *GroupsClient) List(ctx context.Context, filter string) (*[]Group, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/groups",
			Params:      params,
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		Groups []Group `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.Groups, status, nil
}

// Create creates a new Group.
func (c *GroupsClient) Create(ctx context.Context, group Group) (*Group, int, error) {
	var status int
	body, err := json.Marshal(group)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}
	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/groups",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Post(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var newGroup Group
	if err := json.Unmarshal(respBody, &newGroup); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &newGroup, status, nil
}

// Get retrieves a Group.
func (c *GroupsClient) Get(ctx context.Context, id string) (*Group, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var group Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &group, status, nil
}

// GetDeleted retrieves a deleted O365 Group.
// TODO: add test coverage once API supports creating O365 groups.
func (c *GroupsClient) GetDeleted(ctx context.Context, id string) (*Group, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/directory/deletedItems/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var group Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &group, status, nil
}

// Update amends an existing Group.
func (c *GroupsClient) Update(ctx context.Context, group Group) (int, error) {
	var status int
	body, err := json.Marshal(group)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}
	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s", *group.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("GroupsClient.BaseClient.Patch(): %v", err)
	}
	return status, nil
}

// Delete removes a Group.
func (c *GroupsClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("GroupsClient.BaseClient.Delete(): %v", err)
	}
	return status, nil
}

// DeletePermanently removes a deleted O365 Group permanently.
// TODO: add test coverage once API supports creating O365 groups.
func (c *GroupsClient) DeletePermanently(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/directory/deletedItems/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("GroupsClient.BaseClient.Delete(): %v", err)
	}
	return status, nil
}

// ListDeleted retrieves a list of recently deleted O365 groups, optionally filtered using OData.
// TODO: add test coverage once API supports creating O365 groups
func (c *GroupsClient) ListDeleted(ctx context.Context, filter string) (*[]Group, int, error) {
	params := url.Values{}
	if filter != "" {
		params.Add("$filter", filter)
	}
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/directory/deleteditems/microsoft.graph.group",
			Params:      params,
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	var data struct {
		DeletedGroups []Group `json:"value"`
	}
	if err = json.Unmarshal(respBody, &data); err != nil {
		return nil, status, err
	}
	return &data.DeletedGroups, status, nil
}

// RestoreDeleted restores a recently deleted O365 Group.
// TODO: add test coverage once API supports creating O365 groups.
func (c *GroupsClient) RestoreDeleted(ctx context.Context, id string) (*Group, int, error) {
	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/directory/deletedItems/%s/restore", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Post(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var restoredGroup Group
	if err = json.Unmarshal(respBody, &restoredGroup); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &restoredGroup, status, nil
}

// ListMembers retrieves the members of the specified Group.
// id is the object ID of the group.
func (c *GroupsClient) ListMembers(ctx context.Context, id string) (*[]string, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s/members", id),
			Params:      url.Values{"$select": []string{"id"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		Members []struct {
			Type string `json:"@odata.type"`
			Id   string `json:"id"`
		} `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	ret := make([]string, len(data.Members))
	for i, v := range data.Members {
		ret[i] = v.Id
	}
	return &ret, status, nil
}

// GetMember retrieves a single member of the specified Group.
// groupId is the object ID of the group.
// memberId is the object ID of the member object.
func (c *GroupsClient) GetMember(ctx context.Context, groupId, memberId string) (*string, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s/members/%s/$ref", groupId, memberId),
			Params:      url.Values{"$select": []string{"id,url"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		Context string `json:"@odata.context"`
		Type    string `json:"@odata.type"`
		Id      string `json:"id"`
		Url     string `json:"url"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.Id, status, nil
}

// AddMembers adds a new member to a Group.
// First populate the Members field of the Group using the AppendMember method of the model, then call this method.
func (c *GroupsClient) AddMembers(ctx context.Context, group *Group) (int, error) {
	var status int
	// Patching group members support up to 20 members per request
	var memberChunks [][]string
	if group.Members == nil || len(*group.Members) == 0 {
		return status, fmt.Errorf("no members specified")
	}
	members := *group.Members
	max := len(members)
	// Chunk into slices of 20 for batching
	for i := 0; i < max; i += 20 {
		end := i + 20
		if end > max {
			end = max
		}
		memberChunks = append(memberChunks, members[i:end])
	}
	for _, members := range memberChunks {
		// don't fail if a member already exists
		checkMemberAlreadyExists := func(resp *http.Response, o *odata.OData) bool {
			if resp.StatusCode == http.StatusBadRequest && o.Error != nil {
				return o.Error.Match(odata.ErrorAddedObjectReferencesAlreadyExist)
			}
			return false
		}

		data := Group{
			Members: &members,
		}
		body, err := json.Marshal(data)
		if err != nil {
			return status, fmt.Errorf("json.Marshal(): %v", err)
		}
		_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
			Body:                   body,
			ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
			ValidStatusCodes:       []int{http.StatusNoContent},
			ValidStatusFunc:        checkMemberAlreadyExists,
			Uri: Uri{
				Entity:      fmt.Sprintf("/groups/%s", *group.ID),
				HasTenantId: true,
			},
		})
		if err != nil {
			return status, fmt.Errorf("GroupsClient.BaseClient.Patch(): %v", err)
		}
	}
	return status, nil
}

// RemoveMembers removes members from a Group.
// groupId is the object ID of the group.
// memberIds is a *[]string containing object IDs of members to remove.
func (c *GroupsClient) RemoveMembers(ctx context.Context, id string, memberIds *[]string) (int, error) {
	var status int
	if memberIds == nil || len(*memberIds) == 0 {
		return status, fmt.Errorf("no members specified")
	}
	for _, memberId := range *memberIds {
		// check for membership before attempting deletion
		if _, status, err := c.GetMember(ctx, id, memberId); err != nil {
			if status == http.StatusNotFound {
				continue
			}
			return status, err
		}

		// despite the above check, sometimes members are just gone
		checkMemberGone := func(resp *http.Response, o *odata.OData) bool {
			if resp.StatusCode == http.StatusBadRequest && o.Error != nil {
				return o.Error.Match(odata.ErrorRemovedObjectReferencesDoNotExist)
			}
			return false
		}

		var err error
		_, status, _, err = c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
			ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
			ValidStatusCodes:       []int{http.StatusNoContent},
			ValidStatusFunc:        checkMemberGone,
			Uri: Uri{
				Entity:      fmt.Sprintf("/groups/%s/members/%s/$ref", id, memberId),
				HasTenantId: true,
			},
		})
		if err != nil {
			return status, fmt.Errorf("GroupsClient.BaseClient.Delete(): %v", err)
		}
	}
	return status, nil
}

// ListOwners retrieves the owners of the specified Group.
// id is the object ID of the group.
func (c *GroupsClient) ListOwners(ctx context.Context, id string) (*[]string, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s/owners", id),
			Params:      url.Values{"$select": []string{"id"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		Owners []struct {
			Type string `json:"@odata.type"`
			Id   string `json:"id"`
		} `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	ret := make([]string, len(data.Owners))
	for i, v := range data.Owners {
		ret[i] = v.Id
	}
	return &ret, status, nil
}

// GetOwner retrieves a single owner for the specified Group.
// groupId is the object ID of the group.
// ownerId is the object ID of the owning object.
func (c *GroupsClient) GetOwner(ctx context.Context, groupId, ownerId string) (*string, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/groups/%s/owners/%s/$ref", groupId, ownerId),
			Params:      url.Values{"$select": []string{"id,url"}},
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("GroupsClient.BaseClient.Get(): %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("ioutil.ReadAll(): %v", err)
	}
	var data struct {
		Context string `json:"@odata.context"`
		Type    string `json:"@odata.type"`
		Id      string `json:"id"`
		Url     string `json:"url"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}
	return &data.Id, status, nil
}

// AddOwners adds a new owner to a Group.
// First populate the Owners field of the Group using the AppendOwner method of the model, then call this method.
func (c *GroupsClient) AddOwners(ctx context.Context, group *Group) (int, error) {
	var status int
	if group.Owners == nil || len(*group.Owners) == 0 {
		return status, fmt.Errorf("no owners specified")
	}
	for _, owner := range *group.Owners {
		// don't fail if an owner already exists
		checkOwnerAlreadyExists := func(resp *http.Response, o *odata.OData) bool {
			if resp.StatusCode == http.StatusBadRequest && o.Error != nil {
				return o.Error.Match(odata.ErrorAddedObjectReferencesAlreadyExist)
			}
			return false
		}

		data := struct {
			Owner string `json:"@odata.id"`
		}{
			Owner: owner,
		}
		body, err := json.Marshal(data)
		if err != nil {
			return status, fmt.Errorf("json.Marshal(): %v", err)
		}
		_, status, _, err = c.BaseClient.Post(ctx, PostHttpRequestInput{
			Body:                   body,
			ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
			ValidStatusCodes:       []int{http.StatusNoContent},
			ValidStatusFunc:        checkOwnerAlreadyExists,
			Uri: Uri{
				Entity:      fmt.Sprintf("/groups/%s/owners/$ref", *group.ID),
				HasTenantId: true,
			},
		})
		if err != nil {
			return status, fmt.Errorf("GroupsClient.BaseClient.Post(): %v", err)
		}
	}
	return status, nil
}

// RemoveOwners removes owners from a Group.
// groupId is the object ID of the group.
// ownerIds is a *[]string containing object IDs of owners to remove.
func (c *GroupsClient) RemoveOwners(ctx context.Context, id string, ownerIds *[]string) (int, error) {
	var status int
	if ownerIds == nil || len(*ownerIds) == 0 {
		return status, fmt.Errorf("no owners specified")
	}
	for _, ownerId := range *ownerIds {
		// check for ownership before attempting deletion
		if _, status, err := c.GetOwner(ctx, id, ownerId); err != nil {
			if status == http.StatusNotFound {
				continue
			}
			return status, err
		}

		// despite the above check, sometimes owners are just gone
		checkOwnerGone := func(resp *http.Response, o *odata.OData) bool {
			if resp.StatusCode == http.StatusBadRequest && o.Error != nil {
				return o.Error.Match(odata.ErrorRemovedObjectReferencesDoNotExist)
			}
			return false
		}

		var err error
		_, status, _, err = c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
			ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
			ValidStatusCodes:       []int{http.StatusNoContent},
			ValidStatusFunc:        checkOwnerGone,
			Uri: Uri{
				Entity:      fmt.Sprintf("/groups/%s/owners/%s/$ref", id, ownerId),
				HasTenantId: true,
			},
		})
		if err != nil {
			return status, fmt.Errorf("GroupsClient.BaseClient.Delete(): %v", err)
		}
	}
	return status, nil
}
