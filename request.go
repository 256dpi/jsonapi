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

// DocumentExpected returns whether a request using this intent is expected to
// include a JSON API document.
//
// Note: A response from an API may always include a document that at least
// contains one ore more errors.
func (i Intent) DocumentExpected() bool {
	switch i {
	case CreateResource, UpdateResource, SetRelationship,
		AppendToRelationship, RemoveFromRelationship:
		return true
	}

	return false
}

// RequestMethod returns the matching HTTP request method for an Intent.
func (i Intent) RequestMethod() string {
	switch i {
	case ListResources, FindResource, GetRelatedResources, GetRelationship:
		return "GET"
	case CreateResource, AppendToRelationship:
		return "POST"
	case UpdateResource, SetRelationship:
		return "PATCH"
	case DeleteResource, RemoveFromRelationship:
		return "DELETE"
	}

	return ""
}

// A Request contains all JSON API related information parsed from a low level
// request.
type Request struct {
	// The parsed JSON API intent of the request.
	Intent Intent

	// The fragments parsed from the URL of the request.
	Prefix          string
	ResourceType    string
	ResourceID      string
	RelatedResource string
	Relationship    string

	// The requested resources to be included in the response.
	Include []string

	// The pagination details of the request. Zero values mean no pagination
	// details have been provided.
	PageNumber int
	PageSize   int

	// The sorting that has been requested.
	Sorting []string

	// The sparse fieldsets that have been requested.
	Fields map[string][]string

	// The filtering that has been requested.
	Filters map[string][]string
}

// ParseRequest will parse the passed request and return a new Request with the
// parsed data. It will return an error if the content type, request method or
// url is invalid.
//
// Note: The returned error can directly be written using WriteError.
func ParseRequest(r *http.Request, prefix string) (*Request, error) {
	// get method
	method := r.Method

	// map method to action
	if method != "GET" && method != "POST" && method != "PATCH" && method != "DELETE" {
		return nil, BadRequest("Unsupported method")
	}

	// allocate new request
	jr := &Request{
		Prefix: strings.Trim(prefix, "/"),
	}

	// check content type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && contentType != MediaType {
		return nil, BadRequest("Invalid content type header")
	}

	// check accept header
	accept := r.Header.Get("Accept")
	if accept != "" && accept != "*/*" && accept != "application/*" && accept != MediaType {
		return nil, ErrorFromStatus(http.StatusNotAcceptable, "Invalid accept header")
	}

	// de-prefix and trim path
	url := strings.TrimPrefix(strings.Trim(r.URL.Path, "/"), jr.Prefix+"/")

	// split path
	segments := strings.Split(url, "/")
	if len(segments) == 0 || len(segments) > 4 {
		return nil, BadRequest("Invalid URL segment count")
	}

	// check for invalid segments
	for _, s := range segments {
		if s == "" {
			return nil, BadRequest("Found empty URL segments")
		}
	}

	// set resource
	jr.ResourceType = segments[0]
	level := 1

	// set resource id
	if len(segments) > 1 {
		jr.ResourceID = segments[1]
		level = 2
	}

	// set related resource
	if len(segments) == 3 && segments[2] != "relationships" {
		jr.RelatedResource = segments[2]
		level = 3
	}

	// set relationship
	if len(segments) == 4 && segments[2] == "relationships" {
		jr.Relationship = segments[3]
		level = 4
	}

	// final check
	if len(segments) > 2 && (jr.RelatedResource == "" && jr.Relationship == "") {
		return nil, BadRequest("Invalid URL relationship format")
	}

	// calculate intent
	switch method {
	case "GET":
		switch level {
		case 1:
			jr.Intent = ListResources
		case 2:
			jr.Intent = FindResource
		case 3:
			jr.Intent = GetRelatedResources
		case 4:
			jr.Intent = GetRelationship
		}
	case "POST":
		switch level {
		case 1:
			jr.Intent = CreateResource
		case 4:
			jr.Intent = AppendToRelationship
		}
	case "PATCH":
		switch level {
		case 2:
			jr.Intent = UpdateResource
		case 4:
			jr.Intent = SetRelationship
		}
	case "DELETE":
		switch level {
		case 2:
			jr.Intent = DeleteResource
		case 4:
			jr.Intent = RemoveFromRelationship
		}
	}

	// check intent
	if jr.Intent == 0 {
		return nil, BadRequest("The URL and method combination is invalid")
	}

	// check if request should come with a document and has content type set
	if jr.Intent.DocumentExpected() && contentType == "" {
		return nil, BadRequest("Missing content type header")
	}

	for key, values := range r.URL.Query() {
		// set included resources
		if key == "include" {
			for _, v := range values {
				jr.Include = append(jr.Include, strings.Split(v, ",")...)
			}

			continue
		}

		// set sorting
		if key == "sort" {
			for _, v := range values {
				jr.Sorting = append(jr.Sorting, strings.Split(v, ",")...)
			}

			continue
		}

		// set page number
		if key == "page[number]" {
			if len(values) != 1 {
				return nil, BadRequestParam("More than one value", "page[number]")
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, BadRequestParam("Not a number", "page[number]")
			}

			jr.PageNumber = n
			continue
		}

		// set page size
		if key == "page[size]" {
			if len(values) != 1 {
				return nil, BadRequestParam("More than one value", "page[size]")
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, BadRequestParam("Not a number", "page[size]")
			}

			jr.PageSize = n
			continue
		}

		// set sparse fields
		if strings.HasPrefix(key, "fields[") && strings.HasSuffix(key, "]") {
			if jr.Fields == nil {
				jr.Fields = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				jr.Fields[typ] = append(jr.Fields[typ], strings.Split(v, ",")...)
			}
		}

		// set filters
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			if jr.Filters == nil {
				jr.Filters = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				jr.Filters[typ] = append(jr.Filters[typ], strings.Split(v, ",")...)
			}
		}
	}

	// check page size
	if jr.PageNumber > 0 && jr.PageSize <= 0 {
		return nil, BadRequestParam("Missing page size", "page[number]")
	}

	// check page number
	if jr.PageSize > 0 && jr.PageNumber <= 0 {
		return nil, BadRequestParam("Missing page number", "page[size]")
	}

	return jr, nil
}

// Base will generate the base URL for this request, which includes the type and
// id if present.
func (r *Request) Base() string {
	segments := []string{r.Prefix, r.ResourceType}

	// add id if available
	if r.ResourceID != "" {
		segments = append(segments, r.ResourceID)
	}

	return strings.Join(segments, "/")
}

// Self will generate the "self" URL for this request, which includes all path
// elements if available.
func (r *Request) Self() string {
	segments := []string{r.Prefix, r.ResourceType}

	// add id if available
	if r.ResourceID != "" {
		segments = append(segments, r.ResourceID)

		// add related resource or relationship
		if r.RelatedResource != "" {
			segments = append(segments, r.RelatedResource)
		} else if r.Relationship != "" {
			segments = append(segments, "relationships", r.Relationship)
		}
	}

	return strings.Join(segments, "/")
}
