package jsonapi

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalMinimumSinglePayload(t *testing.T) {
	writer := bytes.NewBuffer(nil)

	err := MarshalPayload(writer, &Payload{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
			},
		},
	})
	assert.NoError(t, err)

	assert.JSONEq(t, `{
  		"data": {
    		"type": "foo"
		}
	}`, writer.String())
}

func TestMarshalSinglePayload(t *testing.T) {
	writer := bytes.NewBuffer(nil)

	err := MarshalPayload(writer, &Payload{
		Data: &HybridResource{
			One: &Resource{
				Type: "foo",
				ID:   "1",
			},
		},
	})
	assert.NoError(t, err)

	assert.JSONEq(t, `{
  		"data": {
    		"type": "foo",
    		"id": "1"
		}
	}`, writer.String())
}
