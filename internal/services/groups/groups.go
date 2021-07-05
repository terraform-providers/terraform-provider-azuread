package groups

import (
	"context"
	"fmt"

	"github.com/manicminer/hamilton/msgraph"
)

func groupFindByName(ctx context.Context, client *msgraph.GroupsClient, displayName string) (*[]msgraph.Group, error) {
	filter := fmt.Sprintf("displayName eq '%s'", displayName)
	groups, _, err := client.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("unable to list Groups with filter %q: %+v", filter, err)
	}

	result := make([]msgraph.Group, 0)
	if groups != nil {
		for _, group := range *groups {
			if group.DisplayName != nil && *group.DisplayName == displayName {
				result = append(result, group)
			}
		}
	}

	return &result, nil
}
