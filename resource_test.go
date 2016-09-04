package jsonapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteResourceEmpty(t *testing.T) {
	// TODO: Should this raise an error?

	rec := httptest.NewRecorder()
	err := WriteResource(rec, http.StatusOK, &Resource{}, nil)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"data": {
			"type": ""
		}
	}`, rec.Body.String())
}

func TestWriteResource(t *testing.T) {
	rec := httptest.NewRecorder()
	err := WriteResource(rec, http.StatusOK, &Resource{
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

	rec := httptest.NewRecorder()
	err := WriteResources(rec, http.StatusOK, []*Resource{{}}, nil)
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
	rec := httptest.NewRecorder()
	err := WriteResources(rec, http.StatusOK, []*Resource{
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
	rec := httptest.NewRecorder()
	err := WriteResource(rec, http.StatusOK, &Resource{
		Type: "foo",
		ID:   "1",
		Relationships: Relationships{
			"bar": {
				One: &Document{
					Data: &HybridResource{
						One: &Resource{
							Type: "bar",
							ID:   "2",
						},
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

func TestWriteResourceRelationships(t *testing.T) {
	rec := httptest.NewRecorder()
	err := WriteResource(rec, http.StatusOK, &Resource{
		Type: "foo",
		ID:   "1",
		Relationships: Relationships{
			"bar": {
				Many: []*Document{
					{
						Data: &HybridResource{
							One: &Resource{
								Type: "bar",
								ID:   "2",
							},
						},
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
				"bar": [
					{
						"data": {
							"type": "bar",
							"id": "2"
						}
					}
				]
			}
		}
	}`, rec.Body.String())
}

func BenchmarkWriteResource(b *testing.B) {
	res := &Resource{
		Type: "foo",
		ID:   "1",
		Attributes: Map{
			"foo": "bar",
			"bar": "baz",
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := WriteResource(httptest.NewRecorder(), http.StatusOK, res, nil)
		if err != nil {
			panic(err)
		}
	}
}
