package jsonapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const ContentType = "application/vnd.api+json"

var ErrInvalidRequest = errors.New("invalid request")

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

// ParseRequest will parse the passed request and return a neq Request with the
// parsed data. It will return an error if the content type or url is invalid.
func ParseRequest(r *http.Request, prefix string) (*Request, error) {
	// check content type
	if r.Header.Get("Content-Type") != ContentType {
		return nil, errors.Wrap(ErrInvalidRequest, "invalid content type")
	}

	// de-prefix and trim path
	url := strings.Trim(strings.TrimPrefix(r.URL.Path, prefix), "/")

	// split path
	segments := strings.Split(url, "/")
	if len(segments) == 0 || len(segments) > 4 {
		return nil, errors.Wrap(ErrInvalidRequest, "invalid url segment count")
	}

	// check for invalid segments
	for _, s := range segments {
		if s == "" {
			return nil, errors.Wrap(ErrInvalidRequest, "found empty segments")
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
		return nil, errors.Wrap(ErrInvalidRequest, "invalid relationships")
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
				return nil, errors.Wrap(ErrInvalidRequest, "more than one value")
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, errors.Wrap(ErrInvalidRequest, "not a number")
			}

			req.PageNumber = n
			continue
		}

		// set page size
		if key == "page[size]" {
			if len(values) != 1 {
				return nil, errors.Wrap(ErrInvalidRequest, "more than one value")
			}

			n, err := strconv.Atoi(values[0])
			if err != nil {
				return nil, errors.Wrap(ErrInvalidRequest, "not a number")
			}

			req.PageSize = n
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

	// check page number and page size
	if (req.PageNumber > 0 && req.PageSize == 0) || (req.PageNumber == 0 && req.PageSize > 0) {
		return nil, errors.Wrap(ErrInvalidRequest, "pagination requires both parameters")
	}

	// Parse Body

	return req, nil
}
