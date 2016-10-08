package jsonapi

import (
	"encoding/json"
	"io"
	"net/http"
)

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
	Included []*Resource `json:"included,omitempty"`

	// A set of links related to the primary data.
	Links *DocumentLinks `json:"links,omitempty"`

	// A list of errors that occurred during the request.
	Errors []*Error `json:"errors,omitempty"`

	// Non-standard meta-information about the document.
	Meta Map `json:"meta,omitempty"`
}

// ParseDocument will decode a JSON API document from the passed reader.
//
// Note: If the read document contains errors the first Error will be returned
// as an error.
func ParseDocument(r io.Reader) (*Document, error) {
	// prepare document
	var doc Document

	// decode body
	err := json.NewDecoder(r).Decode(&doc)
	if err != nil {
		return nil, BadRequest(err.Error())
	}

	// check for errors
	if len(doc.Errors) > 0 {
		return nil, doc.Errors[0]
	}

	return &doc, nil
}

// WriteResponse will write the the status and supplied document to the passed
// response writer.
func WriteResponse(res http.ResponseWriter, status int, doc *Document) error {
	// set content type
	res.Header().Set("Content-Type", MediaType)

	// write status
	res.WriteHeader(status)

	// write document
	return json.NewEncoder(res).Encode(doc)
}
