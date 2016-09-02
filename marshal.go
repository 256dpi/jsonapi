package jsonapi

import (
	"encoding/json"
	"io"
)

// MarshalDocument will encode the document to the writer.
// TODO: Remove.
func MarshalDocument(w io.Writer, doc *Document) error {
	err := json.NewEncoder(w).Encode(doc)
	if err != nil {
		return err
	}

	return nil
}
