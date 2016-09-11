package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequestError(t *testing.T) {
	invalidAccept := newTestRequester("GET", "", nil)
	invalidAccept.header["Accept"] = "foo"

	invalidContentType := newTestRequester("GET", "", nil)
	invalidContentType.header["Content-Type"] = "foo"

	missingContentType := newTestRequester("POST", "foo", nil)

	list := []Requester{
		invalidAccept,
		invalidContentType,
		missingContentType,
		newTestRequester("PUT", "", nil),
		newTestRequester("GET", "", nil),
		newTestRequester("POST", "", nil),
		newTestRequester("GET", "/", nil),
		newTestRequester("GET", "foo/bar/baz/qux", nil),
		newTestRequester("GET", "foo/bar/baz/qux/quux", nil),
		newTestRequester("GET", "foo", map[string][]string{
			"page[number]": []string{"bar"},
		}),
		newTestRequester("GET", "foo", map[string][]string{
			"page[size]": []string{"bar"},
		}),
		newTestRequester("GET", "foo", map[string][]string{
			"page[number]": []string{"1"},
		}),
		newTestRequester("GET", "foo", map[string][]string{
			"page[size]": []string{"1"},
		}),
		newTestRequester("GET", "foo", map[string][]string{
			"page[number]": []string{"bar", "baz"},
		}),
		newTestRequester("GET", "foo", map[string][]string{
			"page[size]": []string{"bar", "baz"},
		}),
		newTestRequester("PATCH", "foo", nil),
	}

	for _, r := range list {
		req, err := ParseRequest(r, "")
		assert.Error(t, err)
		assert.Nil(t, req)
	}
}

func TestParseRequestPrefix(t *testing.T) {
	list := map[string]string{
		"bar":           "",
		"/bar":          "",
		"foo/bar":       "foo",
		"/foo/bar":      "foo",
		"foo/bar/":      "/foo",
		"/foo/bar/":     "foo/",
		"baz/foo/bar/":  "/baz/foo",
		"/baz/foo/bar/": "baz/foo/",
	}

	for url, prefix := range list {
		r := newTestRequester("GET", url, nil)

		req, err := ParseRequest(r, prefix)
		assert.NoError(t, err)
		assert.Equal(t, "bar", req.ResourceType)
	}
}

func TestParseRequestResource(t *testing.T) {
	r := newTestRequester("GET", "foo", nil)

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
	}, req)
}

func TestParseRequestResourceID(t *testing.T) {
	r := newTestRequester("GET", "foo/1", nil)

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       FindResource,
		ResourceType: "foo",
		ResourceID:   "1",
	}, req)
}

func TestParseRequestRelatedResource(t *testing.T) {
	r := newTestRequester("GET", "foo/1/bar", nil)

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:          GetRelatedResources,
		ResourceType:    "foo",
		ResourceID:      "1",
		RelatedResource: "bar",
	}, req)
}

func TestParseRequestRelationship(t *testing.T) {
	r := newTestRequester("GET", "foo/1/relationships/bar", nil)

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       GetRelationship,
		ResourceType: "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, req)
}

func TestParseRequestIntent(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		doc    bool
	}{
		{"GET", "/posts", ListResources, false},
		{"GET", "/posts/1", FindResource, false},
		{"POST", "/posts", CreateResource, true},
		{"PATCH", "/posts/1", UpdateResource, true},
		{"DELETE", "/posts/1", DeleteResource, false},
		{"GET", "/posts/1/author", GetRelatedResources, false},
		{"GET", "/posts/1/relationships/author", GetRelationship, false},
		{"PATCH", "/posts/1/relationships/author", SetRelationship, true},
		{"POST", "/posts/1/relationships/comments", AppendToRelationship, true},
		{"DELETE", "/posts/1/relationships/comments", RemoveFromRelationship, true},
	}

	for _, entry := range list {
		r := newTestRequester(entry.method, entry.url, nil)
		r.header["Content-Type"] = MediaType

		req, err := ParseRequest(r, "")
		assert.NoError(t, err)
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.doc, req.Intent.DocumentExpected())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.method, req.Intent.RequestMethod())
	}
}

func TestParseRequestInclude(t *testing.T) {
	r1 := newTestRequester("GET", "foo", map[string][]string{
		"include": []string{"bar,baz"},
	})

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz"},
	}, req)

	r2 := newTestRequester("GET", "foo", map[string][]string{
		"include": []string{"bar", "baz,qux"},
	})

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestSorting(t *testing.T) {
	r1 := newTestRequester("GET", "foo", map[string][]string{
		"sort": []string{"bar,baz"},
	})

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz"},
	}, req)

	r2 := newTestRequester("GET", "foo", map[string][]string{
		"sort": []string{"bar", "baz,qux"},
	})

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestPage(t *testing.T) {
	r := newTestRequester("GET", "foo", map[string][]string{
		"page[number]": []string{"1"},
		"page[size]":   []string{"2"},
	})

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		PageNumber:   1,
		PageSize:     2,
	}, req)
}

func TestParseRequestFields(t *testing.T) {
	r := newTestRequester("GET", "foo", map[string][]string{
		"fields[foo]": []string{"bar,baz"},
	})

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Fields: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestParseRequestFilters(t *testing.T) {
	r := newTestRequester("GET", "foo", map[string][]string{
		"filter[foo]": []string{"bar,baz"},
	})

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Filters: map[string][]string{
			"foo": {"bar", "baz"},
		},
	}, req)
}

func TestZeroIntentRequestMethod(t *testing.T) {
	assert.Empty(t, Intent(0).RequestMethod())
}

func BenchmarkParseRequest(b *testing.B) {
	r := newTestRequester("GET", "foo/1", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}

func BenchmarkParseRequestFilterAndSort(b *testing.B) {
	r := newTestRequester("GET", "foo/1?sort=bar&filter[baz]=qux", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRequest(r, "")
	}
}
