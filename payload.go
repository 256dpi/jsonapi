package jsonapi

// TODO: Add Composer service that helps constructing payload and resources.

// Map is a general purpose map of string keys and arbitrary values.
type Map map[string]interface{}

// See: http://jsonapi.org/format/#document-links.
type PayloadLinks struct {
	Self     string `json:"self,omitempty"`
	Related  string `json:"related,omitempty"`
	First    string `json:"first,omitempty"`
	Previous string `json:"prev,omitempty"`
	Next     string `json:"next,omitempty"`
	Last     string `json:"last,omitempty"`
}

// Payload is the root structure of every JSON API document and used for
// providing embedded relationships.
//
// See: http://jsonapi.org/format/#document-top-level.
type Payload struct {
	// The payload's primary data in the form of a single resource or a list
	// of resources.
	Data *HybridResource `json:"data,omitempty"`

	// A list of resources that are related to the primary data and/or other
	// included resources.
	Included []Resource `json:"included,omitempty"`

	// A set of links related to the primary data.
	Links *PayloadLinks `json:"links,omitempty"`

	// A list of errors that occurred during the request.
	Errors []*Error `json:"errors,omitempty"`

	// Non-standard meta-information about the payload.
	Meta Map `json:"meta,omitempty"`
}

// HybridPayload is a transparent type that enables concrete marshalling and
// unmarshalling of a single payload value or a list of payloads.
type HybridPayload struct {
	// A single Payload.
	One *Payload

	// A list of Payloads.
	Many []*Payload
}

// A Resource is carried by a payload and provides the basic structure for
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
	Relationships map[string]HybridPayload `json:"relationships,omitempty"`

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

// Rationale
//
// Parse JSON to a plain map
// Read map into generalized a structure
// --> do crazy ORM stuff
//
// ORM systems must unmarshal and marshal to a go data structure rather than a
// bytes array
//
// --> do crazy ORM stuff
// Write generalized structure into map
// Generate JSON from a plain map
