package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var objectSuffix = []byte("{")
var arraySuffix = []byte("[")

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

// UnmarshalDocument reads data from a reader and tries to decode a new document.
func UnmarshalDocument(r io.Reader) (*Document, error) {
	var doc Document
	err := json.NewDecoder(r).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
