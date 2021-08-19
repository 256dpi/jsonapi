package jsonapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// NullLink can be used as a link to encode it as a null value.
const NullLink = "NULL"

var nullBytes = []byte("null")

// Link is document link that also can be null by setting NullLink.
type Link string

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Link) UnmarshalJSON(data []byte) error {
	// handle null
	if bytes.Compare(data, nullBytes) == 0 {
		*l = NullLink
		return nil
	}

	// handle string
	var str string
	err := json.Unmarshal(data, &str)
	*l = Link(str)

	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (l Link) MarshalJSON() ([]byte, error) {
	// handle null
	if l == NullLink {
		return nullBytes, nil
	}

	return json.Marshal(string(l))
}

// DocumentLinks are a set of links related to a documents primary data.
//
// See: http://jsonapi.org/format/#document-links.
type DocumentLinks struct {
	Self     Link `json:"self,omitempty"`
	Related  Link `json:"related,omitempty"`
	First    Link `json:"first,omitempty"`
	Previous Link `json:"prev,omitempty"`
	Next     Link `json:"next,omitempty"`
	Last     Link `json:"last,omitempty"`
}

// A Document is the root structure of every JSON API response. It is also used
// to include relationships.
//
// See: http://jsonapi.org/format/#document-top-level.
type Document struct {
	// The documents primary data in the form of a single resource or a list
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
	//
	// Note: Numbers are left as strings to avoid issues with mismatching types
	// when they are later assigned to a struct.
	Meta Map `json:"meta,omitempty"`
}

// ParseDocument will decode a JSON API document from the passed reader.
//
// Note: If the read document contains errors the first Error will be returned
// as an error.
func ParseDocument(r io.Reader) (*Document, error) {
	// TODO: Check document validity more in depth?

	// prepare document
	var doc Document

	// prepare decoder
	dec := json.NewDecoder(r)
	dec.UseNumber()

	// decode body
	err := dec.Decode(&doc)
	if err != nil {
		return nil, BadRequest(err.Error())
	}

	// check for errors
	if len(doc.Errors) > 0 {
		return nil, doc.Errors[0]
	}

	return &doc, nil
}

// WriteResponse will write the status and supplied document to the passed
// response writer.
func WriteResponse(w http.ResponseWriter, status int, doc *Document) error {
	// set content type
	w.Header().Set("Content-Type", MediaType)

	// write status
	w.WriteHeader(status)

	// write document
	return json.NewEncoder(w).Encode(doc)
}
