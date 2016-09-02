package jsonapi

import (
	"encoding/json"
	"io"
)

// UnmarshalDocument reads data from a reader and tries to decode a new document.
// TODO: Remove.
func UnmarshalDocument(r io.Reader) (*Document, error) {
	var doc Document
	err := json.NewDecoder(r).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
