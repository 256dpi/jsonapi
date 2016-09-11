package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/engine/standard"
	"github.com/stretchr/testify/assert"
)

func TestResourceAssignAttributesInvalidTarget(t *testing.T) {
	r := &Resource{
		Attributes: Map{"foo": "foo"},
	}

	err := r.AssignAttributes(nil)
	assert.Error(t, err)
}

func TestResourceAssignAttributes(t *testing.T) {
	var test struct {
		Foo string
	}

	r := &Resource{
		Attributes: Map{"foo": "foo"},
	}

	err := r.AssignAttributes(&test)
	assert.NoError(t, err)
	assert.Equal(t, "foo", test.Foo)
}

func TestResourceAssignAttributesWithTag(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	r := &Resource{
		Attributes: Map{"bar": "foo"},
	}

	err := r.AssignAttributes(&test)
	assert.NoError(t, err)
	assert.Equal(t, "foo", test.Foo)
}

func TestResourceAssignAttributesUnmatchedFields(t *testing.T) {
	var test struct {
		Foo string
	}

	r := &Resource{
		Attributes: Map{"bar": "foo"},
	}

	err := r.AssignAttributes(&test)
	assert.NoError(t, err)
	assert.Equal(t, "", test.Foo)
}

func TestResourceAssignAttributesInvalidType(t *testing.T) {
	var test struct {
		Foo string
	}

	r := &Resource{
		Attributes: Map{"foo": 1},
	}

	err := r.AssignAttributes(&test)
	assert.Error(t, err)
	assert.Equal(t, "", test.Foo)
}

func TestWriteResourceEmpty(t *testing.T) {
	// TODO: Should this raise an error?

	res, rec := constructResponseAndRecorder()

	err := WriteResource(res, http.StatusOK, &Resource{}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": ""
		}
	}`, rec.Body.String())
}

func TestWriteResource(t *testing.T) {
	res, rec := constructResponseAndRecorder()

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
	}`, rec.Body.String())
}

func TestWriteResourcesEmpty(t *testing.T) {
	// TODO: Should this raise an error?

	res, rec := constructResponseAndRecorder()

	err := WriteResources(res, http.StatusOK, []*Resource{{}}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": [
			{
				"type": ""
			}
		]
	}`, rec.Body.String())
}

func TestWriteResources(t *testing.T) {
	res, rec := constructResponseAndRecorder()

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
	}`, rec.Body.String())
}

func TestWriteResourceRelationship(t *testing.T) {
	res, rec := constructResponseAndRecorder()

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
	}`, rec.Body.String())
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
		res := standard.NewResponse(httptest.NewRecorder(), nil)

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
		res := standard.NewResponse(httptest.NewRecorder(), nil)

		err := WriteResources(res, http.StatusOK, resources, nil)
		if err != nil {
			panic(err)
		}
	}
}
