package parse

import "fmt"

type ApplicationPreAuthorizedId struct {
	ObjectId string
	AppId    string
}

func NewApplicationPreAuthorizedID(objectId, appId string) ApplicationPreAuthorizedId {
	return ApplicationPreAuthorizedId{
		ObjectId: objectId,
		AppId:    appId,
	}
}

func (id ApplicationPreAuthorizedId) String() string {
	return id.ObjectId + "/preAuthorizedApplication/" + id.AppId
}

func ApplicationPreAuthorizedID(idString string) (*ApplicationPreAuthorizedId, error) {
	id, err := ObjectSubResourceID(idString, "preAuthorizedApplication")
	if err != nil {
		return nil, fmt.Errorf("unable to parse Pre-Authorized Application ID: %v", err)
	}

	return &ApplicationPreAuthorizedId{
		ObjectId: id.objectId,
		AppId:    id.subId,
	}, nil
}
