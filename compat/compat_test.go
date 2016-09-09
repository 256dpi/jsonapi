package compat

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gonfire/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestParseHTTPRequest(t *testing.T) {
	u, err := url.Parse("foo")
	if err != nil {
		panic(err)
	}

	r := &http.Request{
		Method: "GET",
		URL:    u,
		Header: make(http.Header),
	}

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, jsonapi.ListResources, req.Intent)
	assert.Equal(t, "foo", req.ResourceType)
}

func TestWriteHTTPResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	err := WriteResponse(rec, http.StatusOK, &jsonapi.Document{
		Data: &jsonapi.HybridResource{
			One: &jsonapi.Resource{
				Type: "foo",
			},
		},
	})
	assert.NoError(t, err)
	assert.JSONEq(t, `{
  		"data": {
    		"type": "foo"
		}
	}`, rec.Body.String())
}

func TestWriteResourceHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	err := WriteResource(rec, http.StatusOK, &jsonapi.Resource{
		Type: "foo",
		ID:   "1",
		Attributes: jsonapi.Map{
			"foo": "bar",
		},
	}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": "foo",
			"id": "1",
			"attributes": {
				"foo": "bar"
			}
		}
	}`, rec.Body.String())
}

func TestWriteResourcesHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	err := WriteResources(rec, http.StatusOK, []*jsonapi.Resource{
		{
			Type: "foo",
			ID:   "1",
			Attributes: jsonapi.Map{
				"foo": "bar",
			},
		},
		{
			Type: "foo",
			ID:   "2",
			Attributes: jsonapi.Map{
				"foo": "bar",
			},
		},
	}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": [
			{
				"type": "foo",
				"id": "1",
				"attributes": {
					"foo": "bar"
				}
			},
			{
				"type": "foo",
				"id": "2",
				"attributes": {
					"foo": "bar"
				}
			}
		]
	}`, rec.Body.String())
}

func TestWriteErrorHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, &jsonapi.Error{
		Status: http.StatusNotFound,
		Title:  "Resource Not Found",
		Detail: "The requested resource cannot be found",
	})

	result := rec.Result()

	assert.Equal(t, http.StatusNotFound, result.StatusCode)
	assert.Equal(t, jsonapi.MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "404",
			"title": "Resource Not Found",
			"detail": "The requested resource cannot be found"
		}]
	}`, rec.Body.String())
}

func TestWriteErrorListHTTP(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteErrorList(rec, &jsonapi.Error{
		Status: http.StatusMethodNotAllowed,
	}, &jsonapi.Error{
		Status: http.StatusMethodNotAllowed,
	})

	result := rec.Result()

	assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	assert.Equal(t, jsonapi.MediaType, result.Header.Get("Content-Type"))
	assert.JSONEq(t, `{
		"errors": [{
			"status": "405"
		}, {
			"status": "405"
		}]
	}`, rec.Body.String())
}
