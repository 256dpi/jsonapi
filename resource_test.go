package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResourceEmpty(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResource(res, http.StatusOK, &Resource{}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": ""
		}
	}`, res.Body.String())
}

func TestWriteResource(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResource(res, http.StatusOK, &Resource{
		Type: "foo",
		ID:   "1",
		Attributes: Map{
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
	}`, res.Body.String())
}

func TestWriteResourcesEmpty(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResources(res, http.StatusOK, []*Resource{{}}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": [
			{
				"type": ""
			}
		]
	}`, res.Body.String())
}

func TestWriteResources(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResources(res, http.StatusOK, []*Resource{
		{
			Type: "foo",
			ID:   "1",
			Attributes: Map{
				"foo": "bar",
			},
		},
		{
			Type: "foo",
			ID:   "2",
			Attributes: Map{
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
	}`, res.Body.String())
}

func TestWriteResourceRelationship(t *testing.T) {
	res := httptest.NewRecorder()

	err := WriteResource(res, http.StatusOK, &Resource{
		Type: "foo",
		ID:   "1",
		Relationships: map[string]*Document{
			"bar": {
				Data: &HybridResource{
					One: &Resource{
						Type: "bar",
						ID:   "2",
					},
				},
			},
		},
	}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": "foo",
			"id": "1",
			"relationships": {
				"bar": {
					"data": {
						"type": "bar",
						"id": "2"
					}
				}
			}
		}
	}`, res.Body.String())
}

func BenchmarkWriteResource(b *testing.B) {
	resource := &Resource{
		Type: "foo",
		ID:   "1",
		Attributes: Map{
			"foo": "bar",
			"bar": "baz",
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()

		err := WriteResource(res, http.StatusOK, resource, nil)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkWriteResources(b *testing.B) {
	resources := []*Resource{
		{
			Type: "foo",
			ID:   "1",
			Attributes: Map{
				"foo": "bar",
				"bar": "baz",
			},
		},
		{
			Type: "foo",
			ID:   "1",
			Attributes: Map{
				"foo": "bar",
				"bar": "baz",
			},
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()

		err := WriteResources(res, http.StatusOK, resources, nil)
		if err != nil {
			panic(err)
		}
	}
}
