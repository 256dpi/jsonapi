package jsonapi

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	withServer(func(base string, server *Server) {
		c := NewClient(ClientConfig{
			BaseURI: base,
		})

		// list
		doc, err := c.List("foo")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)

		// find
		doc, err = c.Find("foo", "bar")
		assert.Error(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, &Error{
			Status: http.StatusNotFound,
			Title:  "not found",
			Detail: "unknown resource",
		}, err)

		// create
		doc, err = c.Create(&Resource{
			Type: "foo",
			ID:   "bar",
			Attributes: Map{
				"foo": "bar",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				One: &Resource{
					Type: "foo",
					ID:   "bar",
					Attributes: Map{
						"foo": "bar",
					},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo/bar",
			},
		}, doc)

		// list
		doc, err = c.List("foo")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{
					{
						Type: "foo",
						ID:   "bar",
						Attributes: Map{
							"foo": "bar",
						},
					},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)

		// find
		doc, err = c.Find("foo", "bar")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				One: &Resource{
					Type: "foo",
					ID:   "bar",
					Attributes: Map{
						"foo": "bar",
					},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo/bar",
			},
		}, doc)

		// update
		doc, err = c.Update(&Resource{
			Type: "foo",
			ID:   "bar",
			Attributes: Map{
				"foo": "baz",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				One: &Resource{
					Type: "foo",
					ID:   "bar",
					Attributes: Map{
						"foo": "baz",
					},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo/bar",
			},
		}, doc)

		// list
		doc, err = c.List("foo")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{
					{
						Type: "foo",
						ID:   "bar",
						Attributes: Map{
							"foo": "baz",
						},
					},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)

		// find
		doc, err = c.Find("foo", "bar")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				One: &Resource{
					Type: "foo",
					ID:   "bar",
					Attributes: Map{
						"foo": "baz",
					},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo/bar",
			},
		}, doc)

		// delete
		err = c.Delete("foo", "bar")
		assert.NoError(t, err)

		// list
		doc, err = c.List("foo")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)

		// find
		doc, err = c.Find("foo", "bar")
		assert.Error(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, &Error{
			Status: http.StatusNotFound,
			Title:  "not found",
			Detail: "unknown resource",
		}, err)
	})
}
