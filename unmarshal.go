package jsonapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

var objectSuffix = []byte("{")
var arraySuffix = []byte("[")

func (c *HybridResource) UnmarshalJSON(payload []byte) error {
	if bytes.HasPrefix(payload, objectSuffix) {
		return json.Unmarshal(payload, &c.One)
	}

	if bytes.HasPrefix(payload, arraySuffix) {
		return json.Unmarshal(payload, &c.Many)
	}

	return errors.New("invalid")
}

func (c *HybridPayload) UnmarshalJSON(payload []byte) error {
	if bytes.HasPrefix(payload, objectSuffix) {
		return json.Unmarshal(payload, &c.One)
	}

	if bytes.HasPrefix(payload, arraySuffix) {
		return json.Unmarshal(payload, &c.Many)
	}

	return errors.New("invalid")
}

func UnmarshalPayload(reader io.Reader) (*Payload, error) {
	var payload Payload
	err := json.NewDecoder(reader).Decode(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
