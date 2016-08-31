package jsonapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const ContentType = "application/vnd.api+json"

var ErrInvalidContentType = errors.New("invalid content type")

var ErrInvalidURL = errors.New("invalid url")

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

// TODO: Where to pass the URL prefix?
func ParseRequest(r *http.Request) (*Request, error) {
	// check content type
	if r.Header.Get("Content-Type") != ContentType {
		return nil, ErrInvalidContentType
	}

	// trim and split path
	segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(segments) == 0 || len(segments) > 4 {
		return nil, ErrInvalidURL
	}

	// check for invalid segments
	for _, s := range segments {
		if s == "" {
			return nil, ErrInvalidURL
		}
	}

	// allocate new request
	req := &Request{}

	// set resource
	req.Resource = segments[0]

	// set resource id
	if len(segments) > 1 {
		req.ResourceID = segments[1]
	}

	// set related resource
	if len(segments) == 3 && segments[2] != "relationships" {
		req.RelatedResource = segments[2]
	}

	// set relationship
	if len(segments) == 4 && segments[2] == "relationships" {
		req.Relationship = segments[3]
	}

	// final check
	if len(segments) > 2 && (req.RelatedResource == "" && req.Relationship == "") {
		return nil, ErrInvalidURL
	}

	// set included resources
	for key, values := range r.URL.Query() {
		if key == "include" {
			for _, v := range values {
				req.Include = append(req.Include, strings.Split(v, ",")...)
			}

			continue
		}

		if key == "sort" {
			for _, v := range values {
				req.Sorting = append(req.Sorting, strings.Split(v, ",")...)
			}

			continue
		}

		if key == "page[number]" {
			if len(values) != 1 {
				return nil, ErrInvalidURL
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, ErrInvalidURL
			}

			req.PageNumber = n
			continue
		}

		if key == "page[size]" {
			if len(values) != 1 {
				return nil, ErrInvalidURL
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, ErrInvalidURL
			}

			req.PageSize = n
			continue
		}

		if strings.HasPrefix(key, "fields[") && strings.HasSuffix(key, "]") {
			if req.Fields == nil {
				req.Fields = make(map[string][]string)
			}

			typ := key[7 : len(key)-1]

			for _, v := range values {
				req.Fields[typ] = append(req.Fields[typ], strings.Split(v, ",")...)
			}
		}

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

	// Parse Query Parameters
	// Parse Body

	return req, nil
}
