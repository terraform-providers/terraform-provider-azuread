package groups

import (
	"context"
	"fmt"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func groupFindByName(ctx context.Context, client *msgraph.GroupsClient, displayName string) (*[]msgraph.Group, error) {
	query := odata.Query{
		Filter: fmt.Sprintf("displayName eq '%s'", displayName),
	}
	groups, _, err := client.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("unable to list Groups with filter %q: %+v", query.Filter, err)
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
