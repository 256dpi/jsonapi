// Package jsonapi provides structure ans methods to implement JSON API
// compatible APIs.
package jsonapi

import "net/http"

// MediaType is the official JSON API media type that should be used by
// all requests and responses.
const MediaType = "application/vnd.api+json"

var objectSuffix = []byte("{")
var arraySuffix = []byte("[")

// returns a bad request error
func badRequest(detail string) *Error {
	return &Error{
		Status: http.StatusBadRequest,
		Title:  http.StatusText(http.StatusBadRequest),
		Detail: detail,
	}
}

// returns a bad request error with a parameter source
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
