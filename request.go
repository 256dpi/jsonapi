package jsonapi

import (
	"net/http"
	"net/url"
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

	// CollectionAction is a variation of the following requests:
	// GET /posts/top-titles
	// POST /posts/lock
	// PATCH /posts/settings
	// DELETE /posts/cache
	CollectionAction

	// ResourceAction is a variation of the following requests:
	// GET /posts/1/meta-data
	// POST /posts/1/publish
	// PATCH /posts/1/settings
	// DELETE /posts/1/history
	ResourceAction
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

	// The prefix of the endpoint e.g. "api". It should not contain any prefix
	// or suffix slashes.
	Prefix string

	// The fragments parsed from the URL of the request. The fragments should not
	// contain any prefix or suffix slashes.
	ResourceType     string
	ResourceID       string
	RelatedResource  string
	Relationship     string
	CollectionAction string
	ResourceAction   string

	// The requested resources to be included in the response. This is read
	// from the "include" query parameter.
	Include []string

	// The pagination details of the request. Zero values mean no pagination
	// details have been provided. These values are read from the "page[number]",
	// "page[size]", "page[offset]" and "page[limit]" query parameters. These
	// parameters do not belong to the standard, but are recommended.
	PageNumber int64
	PageSize   int64
	PageOffset int64
	PageLimit  int64

	// The sorting that has been requested. This is read from the "sort" query
	// parameter.
	Sorting []string

	// The sparse fields that have been requested. This is read from the "fields"
	// query parameter.
	Fields map[string][]string

	// The filtering that has been requested. This is read from the "filter"
	// query parameter. This parameter does not belong to the standard, but is
	// recommended.
	Filters map[string][]string
}

// ParseRequest is a short-hand for Parser.ParseRequest and will be removed in
// future releases.
func ParseRequest(r *http.Request, prefix string) (*Request, error) {
	return (&Parser{Prefix: prefix}).ParseRequest(r)
}

// A Parser is used to parse incoming requests.
type Parser struct {
	// Prefix is the expected prefix of the endpoint.
	Prefix string

	// A list of valid collection actions and the allowed methods.
	//
	// Note: Make sure the actions do not conflict with the resource id format.
	CollectionActions map[string][]string

	// A list of valid resource actions and the allowed methods.
	//
	// Note: Make sure the actions do not contain "relationships" or use
	// related resource types.
	ResourceActions map[string][]string
}

// ParseRequest will parse the passed request and return a new Request with the
// parsed data. It will return an error if the content type, request method or
// url is invalid. Any returned error can directly be written using WriteError.
func (p *Parser) ParseRequest(r *http.Request) (*Request, error) {
	// get method
	method := r.Method

	// map method to action
	if method != "GET" && method != "POST" && method != "PATCH" && method != "DELETE" {
		return nil, BadRequest("unsupported method")
	}

	// allocate new request
	req := &Request{
		Prefix: strings.Trim(p.Prefix, "/"),
	}

	// de-prefix and trim path
	location := strings.TrimPrefix(strings.Trim(r.URL.Path, "/"), req.Prefix+"/")

	// split path
	segments := strings.Split(location, "/")
	if len(segments) == 0 || len(segments) > 4 {
		return nil, BadRequest("invalid URL segment count")
	}

	// check for invalid segments
	for _, s := range segments {
		if s == "" {
			return nil, BadRequest("found empty URL segments")
		}
	}

	// set resource
	req.ResourceType = segments[0]
	level := 1

	// return early if a collection action is provided
	if len(segments) == 2 {
		if action, ok := p.CollectionActions[segments[1]]; ok {
			for _, m := range action {
				if method == m {
					req.Intent = CollectionAction
					req.CollectionAction = segments[1]
					return req, nil
				}
			}
		}
	}

	// set resource id
	if len(segments) > 1 {
		req.ResourceID = segments[1]
		level = 2
	}

	// return early if a resource action is provided
	if len(segments) == 3 {
		if action, ok := p.ResourceActions[segments[2]]; ok {
			for _, m := range action {
				if method == m {
					req.Intent = ResourceAction
					req.ResourceAction = segments[2]
					return req, nil
				}
			}
		}
	}

	// set related resource
	if len(segments) == 3 && segments[2] != "relationships" {
		req.RelatedResource = segments[2]
		level = 3
	}

	// set relationship
	if len(segments) == 4 && segments[2] == "relationships" {
		req.Relationship = segments[3]
		level = 4
	}

	// final check
	if len(segments) > 2 && (req.RelatedResource == "" && req.Relationship == "") {
		return nil, BadRequest("invalid URL relationship format")
	}

	// calculate intent
	switch method {
	case "GET":
		switch level {
		case 1:
			req.Intent = ListResources
		case 2:
			req.Intent = FindResource
		case 3:
			req.Intent = GetRelatedResources
		case 4:
			req.Intent = GetRelationship
		}
	case "POST":
		switch level {
		case 1:
			req.Intent = CreateResource
		case 4:
			req.Intent = AppendToRelationship
		}
	case "PATCH":
		switch level {
		case 2:
			req.Intent = UpdateResource
		case 4:
			req.Intent = SetRelationship
		}
	case "DELETE":
		switch level {
		case 2:
			req.Intent = DeleteResource
		case 4:
			req.Intent = RemoveFromRelationship
		}
	}

	// check intent
	if req.Intent == 0 {
		return nil, BadRequest("the URL and method combination is invalid")
	}

	// check headers for standard requests
	if req.Intent != CollectionAction && req.Intent != ResourceAction {
		// check content type header
		contentType := r.Header.Get("Content-Type")
		if contentType != "" && contentType != MediaType {
			return nil, BadRequest("invalid content type header")
		}

		// check accept header
		accept := r.Header.Get("Accept")
		if accept != "" && accept != "*/*" && accept != "application/*" && accept != "application/json" && accept != MediaType {
			return nil, ErrorFromStatus(http.StatusNotAcceptable, "invalid accept header")
		}
	}

	// check if request should come with a document and has no content type set
	if req.Intent.DocumentExpected() && r.Header.Get("Content-Type") == "" {
		return nil, BadRequest("missing content type header")
	}

	for key, values := range r.URL.Query() {
		// set included resources
		if key == "include" {
			for _, v := range values {
				req.Include = append(req.Include, strings.Split(v, ",")...)
			}

			continue
		}

		// set sorting
		if key == "sort" {
			for _, v := range values {
				req.Sorting = append(req.Sorting, strings.Split(v, ",")...)
			}

			continue
		}

		// set page number
		if key == "page[number]" {
			if len(values) != 1 {
				return nil, BadRequestParam("more than one page number", "page[number]")
			}

			n, err := strconv.ParseInt(values[0], 10, 0)
			if err != nil {
				return nil, BadRequestParam("invalid page number", "page[number]")
			}

			req.PageNumber = n
			continue
		}

		// set page size
		if key == "page[size]" {
			if len(values) != 1 {
				return nil, BadRequestParam("more than one page size", "page[size]")
			}

			n, err := strconv.ParseInt(values[0], 10, 0)
			if err != nil {
				return nil, BadRequestParam("invalid page size", "page[size]")
			}

			req.PageSize = n
			continue
		}

		// set page offset
		if key == "page[offset]" {
			if len(values) != 1 {
				return nil, BadRequestParam("more than one page offset", "page[offset]")
			}

			n, err := strconv.ParseInt(values[0], 10, 0)
			if err != nil {
				return nil, BadRequestParam("invalid page offset", "page[offset]")
			}

			req.PageOffset = n
			continue
		}

		// set page limit
		if key == "page[limit]" {
			if len(values) != 1 {
				return nil, BadRequestParam("more than one page limit", "page[limit]")
			}

			n, err := strconv.ParseInt(values[0], 10, 0)
			if err != nil {
				return nil, BadRequestParam("invalid page limit", "page[limit]")
			}

			req.PageLimit = n
			continue
		}

		// set sparse fields
		if strings.HasPrefix(key, "fields[") && strings.HasSuffix(key, "]") {
			if req.Fields == nil {
				req.Fields = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				req.Fields[typ] = append(req.Fields[typ], strings.Split(v, ",")...)
			}
		}

		// set filters
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			if req.Filters == nil {
				req.Filters = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				req.Filters[typ] = append(req.Filters[typ], strings.Split(v, ",")...)
			}
		}
	}

	// check that page number is set if page size is set
	if req.PageNumber > 0 && req.PageSize <= 0 {
		return nil, BadRequestParam("missing page size", "page[number]")
	}

	// check that page size is set if page number is set
	if req.PageSize > 0 && req.PageNumber <= 0 {
		return nil, BadRequestParam("missing page number", "page[size]")
	}

	// check that page limit is set if page offset is set
	if req.PageOffset > 0 && req.PageLimit <= 0 {
		return nil, BadRequestParam("missing page limit", "page[limit]")
	}

	return req, nil
}

// Base will generate the base URL for this request, which includes the type and
// id if present.
func (r *Request) Base() string {
	// prepare segments
	var segments []string

	// add prefix if set
	if r.Prefix != "" {
		segments = append(segments, r.Prefix)
	}

	// add resource type
	segments = append(segments, r.ResourceType)

	// add id if available
	if r.ResourceID != "" {
		segments = append(segments, r.ResourceID)
	}

	return "/" + strings.Join(segments, "/")
}

// Self will generate the "self" URL for this request, which includes all path
// elements if available.
func (r *Request) Self() string {
	// prepare segments
	var segments []string

	// add prefix if set
	if r.Prefix != "" {
		segments = append(segments, r.Prefix)
	}

	// add resource type
	segments = append(segments, r.ResourceType)

	// add id if available
	if r.ResourceID != "" {
		segments = append(segments, r.ResourceID)

		// add related resource or relationship
		if r.RelatedResource != "" {
			segments = append(segments, r.RelatedResource)
		} else if r.Relationship != "" {
			segments = append(segments, "relationships", r.Relationship)
		} else if r.ResourceAction != "" {
			segments = append(segments, r.ResourceAction)
		}
	}

	// add collection action if available
	if r.CollectionAction != "" {
		segments = append(segments, r.CollectionAction)
	}

	return "/" + strings.Join(segments, "/")
}

// Query will collect and return all query parameters from the request.
func (r *Request) Query() url.Values {
	// prepare values
	values := url.Values{}

	// add include
	if len(r.Include) > 0 {
		values.Set("include", strings.Join(r.Include, ","))
	}

	// add page number and size
	if r.PageNumber > 0 && r.PageSize > 0 {
		values.Set("page[number]", strconv.FormatInt(r.PageNumber, 10))
		values.Set("page[size]", strconv.FormatInt(r.PageSize, 10))
	}

	// add page offset and limit
	if r.PageOffset > 0 || r.PageLimit > 0 {
		values.Set("page[offset]", strconv.FormatInt(r.PageOffset, 10))
		values.Set("page[limit]", strconv.FormatInt(r.PageLimit, 10))
	}

	// add include
	if len(r.Include) > 0 {
		values.Set("sort", strings.Join(r.Sorting, ","))
	}

	// add fields
	for name, fields := range r.Fields {
		values.Set("fields["+name+"]", strings.Join(fields, ","))
	}

	// add filters
	for name, filter := range r.Filters {
		values.Set("filer["+name+"]", strings.Join(filter, ","))
	}

	return values
}
