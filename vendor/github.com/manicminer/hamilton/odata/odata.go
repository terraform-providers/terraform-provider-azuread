package odata

import (
	"encoding/json"
	"fmt"
	"regexp"
)

const (
	ErrorAddedObjectReferencesAlreadyExist   = "One or more added object references already exist"
	ErrorConflictingObjectPresentInDirectory = "A conflicting object with one or more of the specified property values is present in the directory"
	ErrorRemovedObjectReferencesDoNotExist   = "One or more removed object references do not exist"
	ErrorServicePrincipalInvalidAppId        = "The appId '.+' of the service principal does not reference a valid application object."
)

// OData is used to unmarshall OData metadata from an API response.
type OData struct {
	Context      *string `json:"@odata.context"`
	MetadataEtag *string `json:"@odata.metadataEtag"`
	Type         *string `json:"@odata.type"`
	Count        *string `json:"@odata.count"`
	NextLink     *string `json:"@odata.nextLink"`
	Delta        *string `json:"@odata.delta"`
	DeltaLink    *string `json:"@odata.deltaLink"`
	Id           *string `json:"@odata.id"`
	Etag         *string `json:"@odata.etag"`

	Error *Error `json:"-"`

	Value *[]json.RawMessage `json:"value"`
}

func (o *OData) UnmarshalJSON(data []byte) error {
	// Perform unmarshalling using a local type
	type odata OData
	var o2 odata
	if err := json.Unmarshal(data, &o2); err != nil {
		return err
	}
	*o = OData(o2)

	// Look for errors in the "error" and "odata.error" fields
	var e map[string]json.RawMessage
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}
	for _, k := range []string{"error", "odata.error"} {
		if v, ok := e[k]; ok {
			var e2 Error
			if err := json.Unmarshal(v, &e2); err != nil {
				return err
			}
			o.Error = &e2
			break
		}
	}
	return nil
}

// Error is used to unmarshal an API error message.
type Error struct {
	Code            *string          `json:"code"`
	Date            *string          `json:"date"`
	Message         *string          `json:"-"`
	RawMessage      *json.RawMessage `json:"message"` // sometimes a string, sometimes an object :/
	ClientRequestId *string          `json:"client-request-id"`
	RequestId       *string          `json:"request-id"`

	InnerError *Error `json:"innerError"` // nested errors

	Details *[]struct {
		Code   *string `json:"code"`
		Target *string `json:"target"`
	} `json:"details"`

	Values *[]struct {
		Item  string `json:"item"`
		Value string `json:"value"`
	} `json:"values"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	// Perform unmarshalling using a local type
	type error Error
	var e2 error
	if err := json.Unmarshal(data, &e2); err != nil {
		return err
	}
	*e = Error(e2)

	// Unmarshal the message, which can be a plain string or an object wrapping a message
	if raw := e.RawMessage; raw != nil && len(*raw) > 0 {
		switch string((*raw)[0]) {
		case "\"":
			var s string
			if err := json.Unmarshal(*raw, &s); err != nil {
				return err
			}
			e.Message = &s
		case "{":
			var m map[string]interface{}
			if err := json.Unmarshal(*raw, &m); err != nil {
				return err
			}
			if v, ok := m["value"]; ok {
				if s, ok := v.(string); ok {
					e.Message = &s
				}
			}
		default:
			return fmt.Errorf("unrecognised error message: %#v", string(*raw))
		}
	}
	return nil
}

func (e Error) String() (s string) {
	if e.Code != nil {
		s = *e.Code
	}
	if e.Message != nil {
		s = fmt.Sprintf("%s: %s", s, *e.Message)
	}
	return
}

func (e Error) Match(errorText string) bool {
	re := regexp.MustCompile(errorText)
	return re.MatchString(e.String())
}
