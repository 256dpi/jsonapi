package compat

import (
	"github.com/gonfire/jsonapi"
	"github.com/labstack/echo/engine/standard"
	"net/http"
)

// ParseRequest is a convenience method to parse a standard http.Request.
//
// See: https://godoc.org/github.com/gonfire/jsonapi#ParseRequest.
func ParseRequest(r *http.Request, prefix string) (*jsonapi.Request, error) {
	return jsonapi.ParseRequest(standard.NewRequest(r, nil), prefix)
}

// WriteResponse is a convenience method to write to a standard
// http.ResponseWriter.
//
// See: https://godoc.org/github.com/gonfire/jsonapi#WriteResponse.
func WriteResponse(res http.ResponseWriter, status int, doc *jsonapi.Document) error {
	return jsonapi.WriteResponse(standard.NewResponse(res, nil), status, doc)
}

// WriteResource is a convenience method to write to a standard
// http.ResponseWriter.
//
// See: https://godoc.org/github.com/gonfire/jsonapi#WriteResource.
func WriteResource(w http.ResponseWriter, status int, resource *jsonapi.Resource, links *jsonapi.DocumentLinks, included ...*jsonapi.Resource) error {
	return jsonapi.WriteResource(standard.NewResponse(w, nil), status, resource, links, included...)
}

// WriteResources is a convenience method to write to a standard
// http.ResponseWriter.
//
// See: https://godoc.org/github.com/gonfire/jsonapi#WriteResources.
func WriteResources(w http.ResponseWriter, status int, resources []*jsonapi.Resource, links *jsonapi.DocumentLinks, included ...*jsonapi.Resource) error {
	return jsonapi.WriteResources(standard.NewResponse(w, nil), status, resources, links, included...)
}

// WriteError is a convenience method to write to a standard
// http.ResponseWriter.
//
// See: https://godoc.org/github.com/gonfire/jsonapi#WriteError.
func WriteError(w http.ResponseWriter, err error) error {
	return jsonapi.WriteError(standard.NewResponse(w, nil), err)
}

// WriteErrorList is a convenience method to write to a standard
// http.ResponseWriter.
//
// See: https://godoc.org/github.com/gonfire/jsonapi#WriteErrorList.
func WriteErrorList(w http.ResponseWriter, errors ...*jsonapi.Error) error {
	return jsonapi.WriteErrorList(standard.NewResponse(w, nil), errors...)
}
