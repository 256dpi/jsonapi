package jsonapi

import (
	"net/http"
	"strconv"
	"strings"
)

// A Request contains all JSON API related information parsed from a low level
// request.
type Request struct {
	// Location
	Resource        string
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
	// check content type
	if req.Header.Get("Content-Type") != ContentType {
		return nil, badRequest("Invalid content type")
	}

	// TODO: Content-Type header is only required when sending data.
	// TODO: Check "Accept" header.

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

	// allocate new request
	r := &Request{}

	// set resource
	r.Resource = segments[0]

	// set resource id
	if len(segments) > 1 {
		r.ResourceID = segments[1]
	}

	// set related resource
	if len(segments) == 3 && segments[2] != "relationships" {
		r.RelatedResource = segments[2]
	}

	// set relationship
	if len(segments) == 4 && segments[2] == "relationships" {
		r.Relationship = segments[3]
	}

	// final check
	if len(segments) > 2 && (r.RelatedResource == "" && r.Relationship == "") {
		return nil, badRequest("Invalid URL relationship format")
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

	// TODO: Parse Body.

	return r, nil
}
