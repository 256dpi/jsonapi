package jsonapi

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalError(t *testing.T) {
	writer := bytes.NewBuffer(nil)
	err := MarshalError(writer, &Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"status": "404",
		"title": "Resource Not Found",
		"detail": "The requested resource cannot be found"
	}`, writer.String())
}

func TestUnmarshalError(t *testing.T) {
	err, _err := UnmarshalError(stringReader(`{
		"status": "404",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
	}`))
	assert.NoError(t, _err)
	assert.Equal(t, &Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	}, err)
}
