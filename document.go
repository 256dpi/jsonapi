package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
)

// TODO: Add Composer service that helps constructing documents and resources.

// DocumentLinks are a set of links related to a documents primary data.
//
// See: http://jsonapi.org/format/#document-links.
type DocumentLinks struct {
	Self     string `json:"self,omitempty"`
	Related  string `json:"related,omitempty"`
	First    string `json:"first,omitempty"`
	Previous string `json:"prev,omitempty"`
	Next     string `json:"next,omitempty"`
	Last     string `json:"last,omitempty"`
}

// A Document is the root structure of every JSON API response. It is also used
// to include relationships.
//
// See: http://jsonapi.org/format/#document-top-level.
type Document struct {
	// The documents's primary data in the form of a single resource or a list
	// of resources.
	Data *HybridResource `json:"data,omitempty"`

	// A list of resources that are related to the primary data and/or other
	// included resources.
	Included []Resource `json:"included,omitempty"`

	// A set of links related to the primary data.
	Links *DocumentLinks `json:"links,omitempty"`

	// A list of errors that occurred during the request.
	Errors []*Error `json:"errors,omitempty"`

	// Non-standard meta-information about the document.
	Meta Map `json:"meta,omitempty"`
}

// HybridDocument is a transparent type that enables concrete marshalling and
// unmarshalling of a single document value or a list of documents.
type HybridDocument struct {
	// A single document.
	One *Document

	// A list of documents.
	Many []*Document
}

// MarshalJSON will either encode a list or a single object.
func (c *HybridDocument) MarshalJSON() ([]byte, error) {
	if c.Many != nil {
		return json.Marshal(c.Many)
	}

	return json.Marshal(c.One)
}

// UnmarshalJSON detects if the passed JSON is a single object or a list.
func (c *HybridDocument) UnmarshalJSON(doc []byte) error {
	if bytes.HasPrefix(doc, objectSuffix) {
		return json.Unmarshal(doc, &c.One)
	}

	if bytes.HasPrefix(doc, arraySuffix) {
		return json.Unmarshal(doc, &c.Many)
	}

	return errors.New("invalid")
}
