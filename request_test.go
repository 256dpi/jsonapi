package jsonapi

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestRequest(method, path string) *http.Request {
	r, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}

	return r
}

func TestParseRequestError(t *testing.T) {
	invalidAccept := newTestRequest("GET", "posts")
	invalidAccept.Header.Set("Accept", "foo")

	invalidContentType := newTestRequest("GET", "posts")
	invalidContentType.Header.Set("Content-Type", "foo")

	list := []struct {
		r *http.Request
		e string
	}{
		{
			r: newTestRequest("HEAD", ""),
			e: "bad request: unsupported method",
		},
		{
			r: newTestRequest("GET", "a/b/c/d/e"),
			e: "bad request: invalid URL segment count",
		},
		{
			r: newTestRequest("GET", "a///d"),
			e: "bad request: found empty URL segments",
		},
		{
			r: newTestRequest("GET", ""),
			e: "bad request: found empty URL segments",
		},
		{
			r: newTestRequest("GET", "a/b/c/d"),
			e: "bad request: invalid URL relationship format",
		},
		{
			r: newTestRequest("POST", "a/b"),
			e: "bad request: the URL and method combination is invalid",
		},
		{
			r: invalidContentType,
			e: "bad request: invalid content type header",
		},
		{
			r: invalidAccept,
			e: "not acceptable: invalid accept header",
		},
		{
			r: newTestRequest("POST", "foo"),
			e: "bad request: missing content type header",
		},
		{
			r: newTestRequest("GET", "foo?page[number]=bar"),
			e: "bad request: invalid page number",
		},
		{
			r: newTestRequest("GET", "foo?page[size]=bar"),
			e: "bad request: invalid page size",
		},
		{
			r: newTestRequest("GET", "foo?page[number]=bar&page[number]=baz"),
			e: "bad request: more than one page number",
		},
		{
			r: newTestRequest("GET", "foo?page[size]=bar&page[size]=baz"),
			e: "bad request: more than one page size",
		},
		{
			r: newTestRequest("GET", "foo?page[offset]=bar"),
			e: "bad request: invalid page offset",
		},
		{
			r: newTestRequest("GET", "foo?page[limit]=bar"),
			e: "bad request: invalid page limit",
		},
		{
			r: newTestRequest("GET", "foo?page[offset]=bar&page[offset]=baz"),
			e: "bad request: more than one page offset",
		},
		{
			r: newTestRequest("GET", "foo?page[limit]=bar&page[limit]=baz"),
			e: "bad request: more than one page limit",
		},
		{
			r: newTestRequest("GET", "foo?page[number]=1"),
			e: "bad request: missing page size",
		},
		{
			r: newTestRequest("GET", "foo?page[offset]=1"),
			e: "bad request: missing page limit",
		},
	}

	for _, i := range list {
		req, err := ParseRequest(i.r, "")
		assert.Error(t, err)
		assert.Equal(t, i.e, err.Error())
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

	for path, prefix := range list {
		r := newTestRequest("GET", path)

		req, err := ParseRequest(r, prefix)
		assert.NoError(t, err)
		assert.Equal(t, "bar", req.ResourceType)
	}
}

func TestParseRequestResource(t *testing.T) {
	r := newTestRequest("GET", "foo")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
	}, req)
}

func TestParseRequestResourceID(t *testing.T) {
	r := newTestRequest("GET", "foo/1")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       FindResource,
		ResourceType: "foo",
		ResourceID:   "1",
	}, req)
}

func TestParseRequestRelatedResource(t *testing.T) {
	r := newTestRequest("GET", "foo/1/bar")

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
	r := newTestRequest("GET", "foo/1/relationships/bar")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       GetRelationship,
		ResourceType: "foo",
		ResourceID:   "1",
		Relationship: "bar",
	}, req)
}

func TestParseRequest(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		doc    bool
		base   string
	}{
		{"GET", "/posts", ListResources, false, "/posts"},
		{"GET", "/posts/1", FindResource, false, "/posts/1"},
		{"POST", "/posts", CreateResource, true, "/posts"},
		{"PATCH", "/posts/1", UpdateResource, true, "/posts/1"},
		{"DELETE", "/posts/1", DeleteResource, false, "/posts/1"},
		{"GET", "/posts/1/author", GetRelatedResources, false, "/posts/1"},
		{"GET", "/posts/1/relationships/author", GetRelationship, false, "/posts/1"},
		{"PATCH", "/posts/1/relationships/author", SetRelationship, true, "/posts/1"},
		{"POST", "/posts/1/relationships/comments", AppendToRelationship, true, "/posts/1"},
		{"DELETE", "/posts/1/relationships/comments", RemoveFromRelationship, true, "/posts/1"},
	}

	for _, entry := range list {
		r := newTestRequest(entry.method, entry.url)
		r.Header.Set("Content-Type", MediaType)

		req, err := ParseRequest(r, "")
		assert.NoError(t, err)
		assert.Equal(t, entry.method, req.Intent.RequestMethod())
		assert.Equal(t, entry.url, req.Path())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.doc, req.Intent.DocumentExpected())
		assert.Equal(t, entry.base, req.Base())
	}
}

func TestParseRequestWithPrefix(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		doc    bool
		base   string
	}{
		{"GET", "/api/posts", ListResources, false, "/api/posts"},
		{"GET", "/api/posts/1", FindResource, false, "/api/posts/1"},
		{"POST", "/api/posts", CreateResource, true, "/api/posts"},
		{"PATCH", "/api/posts/1", UpdateResource, true, "/api/posts/1"},
		{"DELETE", "/api/posts/1", DeleteResource, false, "/api/posts/1"},
		{"GET", "/api/posts/1/author", GetRelatedResources, false, "/api/posts/1"},
		{"GET", "/api/posts/1/relationships/author", GetRelationship, false, "/api/posts/1"},
		{"PATCH", "/api/posts/1/relationships/author", SetRelationship, true, "/api/posts/1"},
		{"POST", "/api/posts/1/relationships/comments", AppendToRelationship, true, "/api/posts/1"},
		{"DELETE", "/api/posts/1/relationships/comments", RemoveFromRelationship, true, "/api/posts/1"},
	}

	for _, entry := range list {
		r := newTestRequest(entry.method, entry.url)
		r.Header.Set("Content-Type", MediaType)

		req, err := ParseRequest(r, "api")
		assert.NoError(t, err)
		assert.Equal(t, entry.method, req.Intent.RequestMethod())
		assert.Equal(t, entry.url, req.Path())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.doc, req.Intent.DocumentExpected())
		assert.Equal(t, entry.base, req.Base())
	}
}

func TestParseRequestActions(t *testing.T) {
	list := []struct {
		method string
		url    string
		intent Intent
		base   string
	}{
		{"GET", "/posts/foo", CollectionAction, "/posts"},
		{"PATCH", "/posts/foo", CollectionAction, "/posts"},
		{"POST", "/posts/1/foo", ResourceAction, "/posts/1"},
		{"DELETE", "/posts/1/foo", ResourceAction, "/posts/1"},
	}

	parser := &Parser{
		CollectionActions: map[string][]string{
			"foo": {"GET", "PATCH"},
		},
		ResourceActions: map[string][]string{
			"foo": {"POST", "DELETE"},
		},
	}

	for _, entry := range list {
		r := newTestRequest(entry.method, entry.url)
		r.Header.Set("Content-Type", MediaType)

		req, err := parser.ParseRequest(r)
		assert.NoError(t, err)
		assert.Empty(t, req.Intent.RequestMethod())
		assert.Equal(t, entry.url, req.Path())
		assert.Equal(t, entry.url, req.Self())
		assert.Equal(t, entry.intent, req.Intent)
		assert.Equal(t, entry.base, req.Base())
	}
}

func TestParseRequestInclude(t *testing.T) {
	r1 := newTestRequest("GET", "foo?include=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz"},
	}, req)

	r2 := newTestRequest("GET", "foo?include=bar&include=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Include:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestSorting(t *testing.T) {
	r1 := newTestRequest("GET", "foo?sort=bar,baz")

	req, err := ParseRequest(r1, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz"},
	}, req)

	r2 := newTestRequest("GET", "foo?sort=bar&sort=baz,qux")

	req, err = ParseRequest(r2, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		Sorting:      []string{"bar", "baz", "qux"},
	}, req)
}

func TestParseRequestPagedPagination(t *testing.T) {
	r := newTestRequest("GET", "foo?page[number]=1&page[size]=5")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		PageNumber:   1,
		PageSize:     5,
	}, req)
}

func TestParseRequestOffsetPagination(t *testing.T) {
	r := newTestRequest("GET", "foo?page[offset]=10&page[limit]=5")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		PageOffset:   10,
		PageLimit:    5,
	}, req)
}

func TestParseRequestOffsetPaginationWithZeroOffset(t *testing.T) {
	r := newTestRequest("GET", "foo?page[offset]=0&page[limit]=5")

	req, err := ParseRequest(r, "")
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:       ListResources,
		ResourceType: "foo",
		PageOffset:   0,
		PageLimit:    5,
	}, req)
}

func TestParseRequestFields(t *testing.T) {
	r := newTestRequest("GET", "foo?fields[foo]=bar,baz")

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
	r := newTestRequest("GET", "foo?filter[foo]=bar,baz")

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

func TestCollectionActionsAcceptHeader(t *testing.T) {
	r := newTestRequest("POST", "posts/foo")
	r.Header.Set("Content-Type", "application/octet-stream")
	r.Header.Set("Accept", "application/octet-stream")

	parser := Parser{
		CollectionActions: map[string][]string{
			"foo": {"POST"},
		},
	}
	req, err := parser.ParseRequest(r)
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:           CollectionAction,
		ResourceType:     "posts",
		CollectionAction: "foo",
	}, req)
}

func TestResourceActionsAcceptHeader(t *testing.T) {
	r := newTestRequest("POST", "posts/1/foo")
	r.Header.Set("Content-Type", "application/octet-stream")
	r.Header.Set("Accept", "application/octet-stream")

	parser := Parser{
		ResourceActions: map[string][]string{
			"foo": {"POST"},
		},
	}
	req, err := parser.ParseRequest(r)
	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Intent:         ResourceAction,
		ResourceType:   "posts",
		ResourceID:     "1",
		ResourceAction: "foo",
	}, req)
}

func TestRequestSelf(t *testing.T) {
	req := Request{
		Intent:       ListResources,
		ResourceType: "posts",
		Include:      []string{"foo", "bar"},
		PageNumber:   1,
		PageSize:     2,
		PageOffset:   3,
		PageLimit:    4,
		Sorting:      []string{"foo", "-bar"},
		Fields: map[string][]string{
			"foo": {"f1", "f2"},
			"bar": {"b1", "b2"},
		},
		Filters: map[string][]string{
			"foo": {"f1", "f2"},
			"bar": {"b1", "b2"},
		},
	}
	assert.Equal(t, strings.Join([]string{
		"/posts",
		"?fields[bar]=b1,b2",
		"&fields[foo]=f1,f2",
		"&filter[bar]=b1,b2",
		"&filter[foo]=f1,f2",
		"&include=foo,bar",
		"&page[limit]=4",
		"&page[number]=1",
		"&page[offset]=3",
		"&page[size]=2",
		"&sort=foo,-bar",
	}, ""), req.Self())
}

func TestRequestQuery(t *testing.T) {
	req := Request{
		Include:    []string{"foo", "bar"},
		PageNumber: 1,
		PageSize:   2,
		PageOffset: 3,
		PageLimit:  4,
		Sorting:    []string{"foo", "-bar"},
		Fields: map[string][]string{
			"foo": {"f5", "f6"},
			"bar": {"b7", "b8"},
		},
		Filters: map[string][]string{
			"foo": {"f9", "f10"},
			"bar": {"b11", "b12"},
		},
	}

	assert.Equal(t, url.Values{
		"include":      []string{"foo,bar"},
		"page[number]": []string{"1"},
		"page[size]":   []string{"2"},
		"page[offset]": []string{"3"},
		"page[limit]":  []string{"4"},
		"sort":         []string{"foo,-bar"},
		"fields[foo]":  []string{"f5,f6"},
		"fields[bar]":  []string{"b7,b8"},
		"filter[foo]":  []string{"f9,f10"},
		"filter[bar]":  []string{"b11,b12"},
	}, req.Query())
}

func BenchmarkParseRequest(b *testing.B) {
	r := newTestRequest("GET", "foo/1")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseRequest(r, "")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkParseRequestFilterAndSort(b *testing.B) {
	r := newTestRequest("GET", "foo/1?sort=bar&filter[baz]=qux")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseRequest(r, "")
		if err != nil {
			panic(err)
		}
	}
}
