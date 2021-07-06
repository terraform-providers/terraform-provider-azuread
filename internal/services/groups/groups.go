package groups

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"
)

func groupDefaultMailNickname() string {
	charSet := "0123456789abcdef"
	result := make([]byte, 9)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 9; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	resultString := string(result)
	return resultString[:8] + "-" + resultString[8:]
}

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
