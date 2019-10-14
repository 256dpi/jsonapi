package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var objectPrefix = []byte("{")
var arrayPrefix = []byte("[")

// A Resource is carried by a document and provides the basic structure for
// JSON API resource objects and resource identifier objects.
//
// See: http://jsonapi.org/format/#document-resource-objects and
// http://jsonapi.org/format/#document-resource-identifier-objects.
type Resource struct {
	// The mandatory type of the resource.
	Type string `json:"type"`

	// The mandatory id of the resource.
	//
	// Exception: The id is not required when the resource object
	// originates at the client and represents a new resource to be created on
	// the server.
	ID string `json:"id,omitempty"`

	// An attributes map representing some of the resource's data.
	//
	// Note: Numbers are left as strings to avoid issues with mismatching types
	// when they are later assigned to a struct.
	Attributes Map `json:"attributes,omitempty"`

	// A relationships object describing relationships between the resource and
	// other JSON API resources.
	Relationships map[string]*Document `json:"relationships,omitempty"`

	// Non-standard meta-information about the resource.
	//
	// Note: Numbers are left as strings to avoid issues with mismatching types
	// when they are later assigned to a struct.
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
func (r *HybridResource) MarshalJSON() ([]byte, error) {
	if r.Many != nil {
		return json.Marshal(r.Many)
	}

	return json.Marshal(r.One)
}

// UnmarshalJSON detects if the passed JSON is a single object or a list.
func (r *HybridResource) UnmarshalJSON(doc []byte) error {
	// prepare decoder
	dec := json.NewDecoder(bytes.NewReader(doc))
	dec.UseNumber()

	if bytes.HasPrefix(doc, objectPrefix) {
		return dec.Decode(&r.One)
	}

	if bytes.HasPrefix(doc, arrayPrefix) {
		return dec.Decode(&r.Many)
	}

	return errors.New("expected data to be an object or array")
}

// WriteResource will wrap the passed resource, links and included resources in
// a document and write it to the passed response writer.
func WriteResource(w http.ResponseWriter, status int, resource *Resource, links *DocumentLinks, included ...*Resource) error {
	return WriteResponse(w, status, &Document{
		Data: &HybridResource{
			One: resource,
		},
		Links:    links,
		Included: included,
	})
}

// WriteResources will wrap the passed resources, links and included resources
// in a document and write it to the passed response writer.
func WriteResources(w http.ResponseWriter, status int, resources []*Resource, links *DocumentLinks, included ...*Resource) error {
	return WriteResponse(w, status, &Document{
		Data: &HybridResource{
			Many: resources,
		},
		Links:    links,
		Included: included,
	})
}
