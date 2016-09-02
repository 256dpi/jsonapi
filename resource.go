package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Map is a general purpose map of string keys and arbitrary values.
type Map map[string]interface{}

// A Resource is carried by a document and provides the basic structure for
// JSON API resource objects and resource identifier objects.
//
// See: http://jsonapi.org/format/#document-resource-objects.
// See: http://jsonapi.org/format/#document-resource-identifier-objects.
type Resource struct {
	// The mandatory type of the resource.
	Type string `json:"type"`

	// The mandatory id of the resource.
	//
	// Exception: The id is not required when the resource object
	// originates at the client and represents a new resource to be created on
	// the server.
	ID string `json:"id,omitempty"`

	// An attributes object representing some of the resource's data.
	Attributes Map `json:"attributes,omitempty"`

	// A relationships object describing relationships between the resource and
	// other JSON API resources.
	Relationships map[string]HybridDocument `json:"relationships,omitempty"`

	// Non-standard meta-information about the resource.
	Meta Map `json:"meta,omitempty"`
}

// HybridResource is a transparent type that enables concrete marshalling and
// unmarshalling of a single resource value or a list of resources.
type HybridResource struct {
	// A single Resource.
	One *Resource

	// A list of Resources.
	Many []*Resource
}

// MarshalJSON will either encode a list or a single object.
func (c *HybridResource) MarshalJSON() ([]byte, error) {
	if c.Many != nil {
		return json.Marshal(c.Many)
	}

	return json.Marshal(c.One)
}

// UnmarshalJSON detects if the passed JSON is a single object or a list.
func (c *HybridResource) UnmarshalJSON(doc []byte) error {
	if bytes.HasPrefix(doc, objectSuffix) {
		return json.Unmarshal(doc, &c.One)
	}

	if bytes.HasPrefix(doc, arraySuffix) {
		return json.Unmarshal(doc, &c.Many)
	}

	return errors.New("invalid")
}
