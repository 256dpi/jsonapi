package jsonapi

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerCRUD(t *testing.T) {
	withServer(func(client *Client, server *Server) {
		// list
		doc, err := client.List("foo")
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
		doc, err = client.Find("foo", "bar")
		assert.Error(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, &Error{
			Status: http.StatusNotFound,
			Title:  "not found",
			Detail: "unknown resource",
		}, err)

		// create
		doc, err = client.Create(&Resource{
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
		doc, err = client.List("foo")
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
		doc, err = client.Find("foo", "bar")
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
		doc, err = client.Update(&Resource{
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
		doc, err = client.List("foo")
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
		doc, err = client.Find("foo", "bar")
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
		err = client.Delete("foo", "bar")
		assert.NoError(t, err)

		// list
		doc, err = client.List("foo")
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
		doc, err = client.Find("foo", "bar")
		assert.Error(t, err)
		assert.NotNil(t, doc)
		assert.Equal(t, &Error{
			Status: http.StatusNotFound,
			Title:  "not found",
			Detail: "unknown resource",
		}, err)
	})
}

func TestServerPagination(t *testing.T) {
	withServer(func(client *Client, server *Server) {
		server.Data["foo"] = map[string]*Resource{}
		for i := 0; i < 5; i++ {
			id := strconv.Itoa(i)
			server.Data["foo"][id] = &Resource{
				Type: "foo",
				ID:   id,
			}
		}

		// all
		doc, err := client.List("foo")
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{
					{Type: "foo", ID: "0"},
					{Type: "foo", ID: "1"},
					{Type: "foo", ID: "2"},
					{Type: "foo", ID: "3"},
					{Type: "foo", ID: "4"},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)

		// number and size
		doc, err = client.List("foo", Request{
			PageNumber: 1,
			PageSize:   2,
		})
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{
					{Type: "foo", ID: "2"},
					{Type: "foo", ID: "3"},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)

		// offset and limit
		doc, err = client.List("foo", Request{
			PageOffset: 3,
			PageLimit:  5,
		})
		assert.NoError(t, err)
		assert.Equal(t, &Document{
			Data: &HybridResource{
				Many: []*Resource{
					{Type: "foo", ID: "3"},
					{Type: "foo", ID: "4"},
				},
			},
			Links: &DocumentLinks{
				Self: "/foo",
			},
		}, doc)
	})
}
