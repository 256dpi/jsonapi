package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

var objectSuffix = []byte("{")
var arraySuffix = []byte("[")

var responseDocumentPool = sync.Pool{
	New: func() interface{} {
		return &Document{
			Data: &HybridResource{},
		}
	},
}

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
	Relationships map[string]*Document `json:"relationships,omitempty"`

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
func (r *HybridResource) MarshalJSON() ([]byte, error) {
	if r.Many != nil {
		return json.Marshal(r.Many)
	}

	return json.Marshal(r.One)
}

// UnmarshalJSON detects if the passed JSON is a single object or a list.
func (r *HybridResource) UnmarshalJSON(doc []byte) error {
	if bytes.HasPrefix(doc, objectSuffix) {
		return json.Unmarshal(doc, &r.One)
	}

	if bytes.HasPrefix(doc, arraySuffix) {
		return json.Unmarshal(doc, &r.Many)
	}

	return errors.New("Expected data to be an object or array")
}

// WriteResource will wrap the passed resource, links and included resources in
// a document and write it to the passed response writer.
func WriteResource(w http.ResponseWriter, status int, res *Resource, links *DocumentLinks, included ...*Resource) error {
	// TODO: Validate resource?

	// get document from pool
	doc := getResponseDocumentFromPool()

	// put document back when finished
	defer responseDocumentPool.Put(doc)

	// set data
	doc.Data.One = res
	doc.Links = links
	doc.Included = included

	return WriteResponse(w, status, doc)
}

// WriteResources will wrap the passed resources, links and included resources
// in a document and write it to the passed response writer.
func WriteResources(w http.ResponseWriter, status int, res []*Resource, links *DocumentLinks, included ...*Resource) error {
	// get document from pool
	doc := getResponseDocumentFromPool()

	// put document back when finished
	defer responseDocumentPool.Put(doc)

	// set data
	doc.Data.Many = res
	doc.Links = links
	doc.Included = included

	return WriteResponse(w, status, doc)
}

func getResponseDocumentFromPool() *Document {
	// get document from pool
	doc := responseDocumentPool.Get().(*Document)

	// reset document
	doc.Data.One = nil
	doc.Data.Many = nil
	doc.Included = nil
	doc.Links = nil
	doc.Errors = nil
	doc.Meta = nil

	return doc
}
