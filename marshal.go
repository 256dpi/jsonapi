package jsonapi

import (
	"encoding/json"
	"io"
)

func (c *HybridResource) MarshalJSON() ([]byte, error) {
	if c.Many != nil {
		return json.Marshal(c.Many)
	}

	return json.Marshal(c.One)
}

func (c *HybridPayload) MarshalJSON() ([]byte, error) {
	if c.Many != nil {
		return json.Marshal(c.Many)
	}

	return json.Marshal(c.One)
}

func MarshalPayload(writer io.Writer, payload *Payload) error {
	err := json.NewEncoder(writer).Encode(payload)
	if err != nil {
		return err
	}

	return nil
}
