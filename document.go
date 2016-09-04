package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

var responseDocumentPool = sync.Pool{
	New: func() interface{} {
		return &Document{
			Data: &HybridResource{},
		}
	},
}

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

// ParseBody will decode a JSON API document from the passed request body reader.
func ParseBody(r io.Reader) (*Document, error) {
	// prepare document
	var doc Document

	// decode body
	err := json.NewDecoder(r).Decode(&doc)
	if err != nil {
		return nil, badRequest(err.Error())
	}

	// check for errors
	if len(doc.Errors) > 0 {
		return nil, badRequest("Request contains errors")
	}

	// check data
	if doc.Data == nil || (doc.Data.One == nil && len(doc.Data.Many) == 0) {
		return nil, badRequest("Mising data")
	}

	return &doc, nil
}

// WriteResponse will write the the status and supplied document to the passed
// response writer.
func WriteResponse(w http.ResponseWriter, status int, doc *Document) error {
	// set content type
	w.Header().Set("Content-Type", MediaType)

	// write status
	w.WriteHeader(status)

	// write document
	return json.NewEncoder(w).Encode(doc)
}

// WriteResource will wrap the passed resource, links and included resources in
// a document and write it to the passed response writer.
func WriteResource(w http.ResponseWriter, status int, res *Resource, links *DocumentLinks, included ...*Resource) error {
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
