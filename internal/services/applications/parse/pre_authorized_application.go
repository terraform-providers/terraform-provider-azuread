package parse

import "fmt"

type PreAuthorizedApplicationId struct {
	ObjectId string
	AppId    string
}

func NewPreAuthorizedApplicationID(objectId, appId string) PreAuthorizedApplicationId {
	return PreAuthorizedApplicationId{
		ObjectId: objectId,
		AppId:    appId,
	}
}

func (id PreAuthorizedApplicationId) String() string {
	return id.ObjectId + "/preAuthorizedApplication/" + id.AppId
}

func PreAuthorizedApplicationID(idString string) (*PreAuthorizedApplicationId, error) {
	id, err := ObjectSubResourceID(idString, "preAuthorizedApplication")
	if err != nil {
		return nil, fmt.Errorf("unable to parse Pre-Authorized Application ID: %v", err)
	}

	return &PreAuthorizedApplicationId{
		ObjectId: id.objectId,
		AppId:    id.subId,
	}, nil
}
