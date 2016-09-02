package jsonapi

import (
	"fmt"
	"net/http"
)

// Error objects provide additional information about problems encountered while
// performing an operation.
//
// See: http://jsonapi.org/format/#errors.
type Error struct {
	// A unique identifier for this particular occurrence of the problem.
	ID string `json:"id,omitempty"`

	// Continuing links to other resources.
	Links *ErrorLinks `json:"links,omitempty"`

	// The HTTP status code applicable to this problem.
	Status int `json:"status,string,omitempty"`

	// An application-specific error code.
	Code string `json:"code,omitempty"`

	// A short, human-readable summary of the problem.
	Title string `json:"title,omitempty"`

	// A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail,omitempty"`

	// A parameter or pointer reference to the source of the error.
	Source *ErrorSource `json:"source,omitempty"`

	// Non-standard meta-information about the error.
	Meta Map `json:"meta,omitempty"`
}

// ErrorLinks contains continuing links to other resources.
//
// See: http://jsonapi.org/format/#errors.
type ErrorLinks struct {
	// A link that leads to further details about this particular occurrence of
	// the problem.
	About string `json:"about"`
}

// ErrorSource contains a parameter or pointer reference to the source of the
// error.
//
// See: http://jsonapi.org/format/#errors.
type ErrorSource struct {
	// A string indicating which URI query parameter caused the error.
	Parameter string `json:"parameter,omitempty"`

	// A JSON Pointer to the associated entity in the request document.
	Pointer string `json:"pointer,omitempty"`
}

// Error returns a string representation of the error for logging purposes.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

// WriteErrorFromStatus will write an error to the response writer that has
// been derived from the passed status code.
//
// Note: If the passed status code is not a valid HTTP status code, a 500 status
// code will be used instead.
func WriteErrorFromStatus(w http.ResponseWriter, status int) {
	// get text
	str := http.StatusText(status)

	// check text
	if str == "" {
		status = http.StatusInternalServerError
		str = http.StatusText(http.StatusInternalServerError)
	}

	WriteDocument(w, status, &Document{
		Errors: []*Error{
			{
				Status: status,
				Title:  str,
			},
		},
	})
}

// WriteError will write the passed error to the response writer.
//
// Note: If the supplied error is not an Error it will call WriteErrorFromStatus
// with StatusInternalServerError. Does the passed Error have an invalid or zero
// Status it will be corrected to 500 - Internal Server Error.
func WriteError(w http.ResponseWriter, err error) {
	anError, ok := err.(*Error)
	if !ok {
		WriteErrorFromStatus(w, http.StatusInternalServerError)
		return
	}

	// set status
	if str := http.StatusText(anError.Status); str == "" {
		anError.Status = http.StatusInternalServerError
	}

	WriteDocument(w, anError.Status, &Document{
		Errors: []*Error{anError},
	})
}

// TODO: Write a list of errors?
