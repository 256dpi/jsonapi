package jsonapi

import (
	"net/http"
	"strconv"
	"strings"
)

// An Intent represents a valid combination of a request method and a URL pattern.
type Intent int

const (
	_ Intent = iota

	// ListResources is a variation of the following request:
	// GET /posts
	ListResources

	// FindResource is a variation of the following request:
	// GET /posts/1
	FindResource

	// CreateResource is a variation of the following request:
	// POST /posts
	CreateResource

	// UpdateResource is a variation of the following request:
	// PATCH /posts/1
	UpdateResource

	// DeleteResource is a variation of the following request:
	// DELETE /posts/1
	DeleteResource

	// GetRelatedResources is a variation of the following requests:
	// GET /posts/1/author
	// GET /posts/1/comments
	GetRelatedResources

	// GetRelationship is a variation of the following requests:
	// GET /posts/1/relationships/author
	// GET /posts/1/relationships/comments
	GetRelationship

	// SetRelationship is a variation of the following requests:
	// PATCH /posts/1/relationships/author.
	// PATCH /posts/1/relationships/comments.
	SetRelationship

	// AppendToRelationship is a variation of the following request:
	// POST /posts/1/relationships/comments
	AppendToRelationship

	// RemoveFromRelationship is a variation of the following request:
	// DELETE /posts/1/relationships/comments
	RemoveFromRelationship
)

// A Request contains all JSON API related information parsed from a low level
// request.
type Request struct {
	// Intent
	Intent Intent

	// Location
	ResourceType    string
	ResourceID      string
	RelatedResource string
	Relationship    string

	// Inclusion of Related Resources
	Include []string

	// Pagination
	PageNumber int
	PageSize   int

	// Sorting
	Sorting []string

	// Sparse Fieldsets
	Fields map[string][]string

	// Filtering
	Filters map[string][]string
}

// ParseRequest will parse the passed request and return a new Request with the
// parsed data. It will return an error if the content type or url is invalid.
//
// Note: The returned error can directly be written using WriteError.
func ParseRequest(req *http.Request, prefix string) (*Request, error) {
	// get method
	method := req.Method

	// set overridden method if available
	if req.Header.Get("X-HTTP-Method-Override") != "" {
		method = req.Header.Get("X-HTTP-Method-Override")
	}

	// map method to action
	if method != "GET" && method != "POST" && method != "PATCH" && method != "DELETE" {
		return nil, badRequest("Unsupported method")
	}

	// allocate new request
	r := &Request{}

	// check content type header
	contentType := req.Header.Get("Content-Type")
	if contentType != "" && contentType != MediaType {
		return nil, badRequest("Invalid content type header")
	}

	// check accept header
	accept := req.Header.Get("Accept")
	if accept != MediaType {
		return nil, badRequest("Invalid accept header")
	}

	// de-prefix and trim path
	url := strings.Trim(strings.TrimPrefix(req.URL.Path, prefix), "/")

	// split path
	segments := strings.Split(url, "/")
	if len(segments) == 0 || len(segments) > 4 {
		return nil, badRequest("Invalid URL segment count")
	}

	// check for invalid segments
	for _, s := range segments {
		if s == "" {
			return nil, badRequest("Found empty URL segments")
		}
	}

	// set resource
	r.ResourceType = segments[0]
	level := 1

	// set resource id
	if len(segments) > 1 {
		r.ResourceID = segments[1]
		level = 2
	}

	// set related resource
	if len(segments) == 3 && segments[2] != "relationships" {
		r.RelatedResource = segments[2]
		level = 3
	}

	// set relationship
	if len(segments) == 4 && segments[2] == "relationships" {
		r.Relationship = segments[3]
		level = 4
	}

	// final check
	if len(segments) > 2 && (r.RelatedResource == "" && r.Relationship == "") {
		return nil, badRequest("Invalid URL relationship format")
	}

	// calculate intent
	switch method {
	case "GET":
		switch level {
		case 1:
			r.Intent = ListResources
		case 2:
			r.Intent = FindResource
		case 3:
			r.Intent = GetRelatedResources
		case 4:
			r.Intent = GetRelationship
		}
	case "POST":
		switch level {
		case 1:
			r.Intent = CreateResource
		case 4:
			r.Intent = AppendToRelationship
		}
	case "PATCH":
		switch level {
		case 2:
			r.Intent = UpdateResource
		case 4:
			r.Intent = SetRelationship
		}
	case "DELETE":
		switch level {
		case 2:
			r.Intent = DeleteResource
		case 4:
			r.Intent = RemoveFromRelationship
		}
	}

	// check intent
	if r.Intent == 0 {
		return nil, badRequest("The URL and method combination is invalid")
	}

	// check if request should come with a document and has content type set
	if r.DocumentExpected() && contentType == "" {
		return nil, badRequest("Missing content type header")
	}

	for key, values := range req.URL.Query() {
		// set included resources
		if key == "include" {
			for _, v := range values {
				r.Include = append(r.Include, strings.Split(v, ",")...)
			}

			continue
		}

		// set sorting
		if key == "sort" {
			for _, v := range values {
				r.Sorting = append(r.Sorting, strings.Split(v, ",")...)
			}

			continue
		}

		// set page number
		if key == "page[number]" {
			if len(values) != 1 {
				return nil, badRequestParam("More than one value", "page[number]")
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, badRequestParam("Not a number", "page[number]")
			}

			r.PageNumber = n
			continue
		}

		// set page size
		if key == "page[size]" {
			if len(values) != 1 {
				return nil, badRequestParam("More than one value", "page[size]")
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, badRequestParam("Not a number", "page[size]")
			}

			r.PageSize = n
			continue
		}

		// set sparse fields
		if strings.HasPrefix(key, "fields[") && strings.HasSuffix(key, "]") {
			if r.Fields == nil {
				r.Fields = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				r.Fields[typ] = append(r.Fields[typ], strings.Split(v, ",")...)
			}
		}

		// set filters
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			if r.Filters == nil {
				r.Filters = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				r.Filters[typ] = append(r.Filters[typ], strings.Split(v, ",")...)
			}
		}
	}

	// check page size
	if r.PageNumber > 0 && r.PageSize <= 0 {
		return nil, badRequestParam("Missing page size", "page[number]")
	}

	// check page number
	if r.PageSize > 0 && r.PageNumber <= 0 {
		return nil, badRequestParam("Missing page number", "page[size]")
	}

	return r, nil
}

// DocumentExpected returns whether the request is expected to come with a
// JSON API document.
func (r *Request) DocumentExpected() bool {
	switch r.Intent {
	case CreateResource, UpdateResource, SetRelationship,
		AppendToRelationship, RemoveFromRelationship:
		return true
	}

	return false
}

func badRequest(detail string) *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: detail,
	}
}

func badRequestParam(detail, param string) *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: detail,
		Source: &ErrorSource{
			Parameter: param,
		},
	}
}
