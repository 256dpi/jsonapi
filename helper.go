package jsonapi

import "net/http"

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
