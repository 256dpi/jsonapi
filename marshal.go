package jsonapi

import (
	"encoding/json"
	"io"
)

// MarshalJSON will either encode a list or a single object.
func (c *HybridResource) MarshalJSON() ([]byte, error) {
	if c.Many != nil {
		return json.Marshal(c.Many)
	}

	return json.Marshal(c.One)
}

// MarshalJSON will either encode a list or a single object.
func (c *HybridDocument) MarshalJSON() ([]byte, error) {
	if c.Many != nil {
		return json.Marshal(c.Many)
	}

	return json.Marshal(c.One)
}

// MarshalDocument will encode the document to the writer.
func MarshalDocument(w io.Writer, doc *Document) error {
	err := json.NewEncoder(w).Encode(doc)
	if err != nil {
		return err
	}

	return nil
}
