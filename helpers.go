package jsonapi

import "net/http"

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
