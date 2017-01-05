package jsonapi

import (
	"fmt"
	"net/http"
	"sync"
)

var singleErrorDocumentPool = sync.Pool{
	New: func() interface{} {
		return &Document{
			Errors: make([]*Error, 1),
		}
	},
}

var multiErrorDocumentPool = sync.Pool{
	New: func() interface{} {
		return &Document{}
	},
}

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
	//
	// Note: Numbers are left as strings to avoid issues with mismatching types
	// when they are later assigned to a struct.
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

// WriteError will write the passed error to the response writer.
//
// Note: If the supplied error is not an Error a new InternalServerError is used
// instead. Does the passed Error have an invalid or zero status code it will be
// corrected to the Internal Server Error status code.
func WriteError(w http.ResponseWriter, err error) error {
	anError, ok := err.(*Error)
	if !ok {
		anError = InternalServerError("")
	}

	// set status
	if str := http.StatusText(anError.Status); str == "" {
		anError.Status = http.StatusInternalServerError
	}

	// get document from pool
	doc := singleErrorDocumentPool.Get().(*Document)

	// put document back when finished
	defer func() {
		doc.Errors[0] = nil
		singleErrorDocumentPool.Put(doc)
	}()

	// set error
	doc.Errors[0] = anError

	return WriteResponse(w, anError.Status, doc)
}

// WriteErrorList will write the passed errors to the the response writer.
// The function will calculate a common status code for all the errors.
//
// Does a passed Error have an invalid or zero status code it will be corrected
// to the Internal Server Error status code.
func WriteErrorList(w http.ResponseWriter, errors ...*Error) error {
	// write internal server error if no errors are passed
	if len(errors) == 0 {
		return WriteError(w, nil)
	}

	// prepare common status
	commonStatus := 0

	for i, err := range errors {
		// check for zero and invalid status
		if str := http.StatusText(err.Status); str == "" {
			err.Status = http.StatusInternalServerError
		}

		// take the first status directly
		if i == 0 {
			commonStatus = err.Status
			continue
		}

		// check if the same or already 500
		if commonStatus == err.Status || commonStatus == 500 {
			continue
		}

		// settle on 500 if already in 500er range
		if err.Status >= 500 {
			commonStatus = 500
			continue
		}

		// settle on 400 if in 400er range
		commonStatus = 400
	}

	// get document from pool
	doc := multiErrorDocumentPool.Get().(*Document)

	// put document back when finished
	defer func() {
		doc.Errors = nil
		multiErrorDocumentPool.Put(doc)
	}()

	// set errors
	doc.Errors = errors

	return WriteResponse(w, commonStatus, doc)
}

// ErrorFromStatus will return an error that has been derived from the passed
// status code.
//
// Note: If the passed status code is not a valid HTTP status code, an Internal
// Server Error status code will be used instead.
func ErrorFromStatus(status int, detail string) *Error {
	// get text
	str := http.StatusText(status)

	// check text
	if str == "" {
		status = http.StatusInternalServerError
		str = http.StatusText(http.StatusInternalServerError)
	}

	return &Error{
		Status: status,
		Title:  str,
		Detail: detail,
	}
}

// NotFound returns a new not found error.
func NotFound(detail string) *Error {
	return ErrorFromStatus(http.StatusNotFound, detail)
}

// BadRequest returns a new bad request error.
func BadRequest(detail string) *Error {
	return ErrorFromStatus(http.StatusBadRequest, detail)
}

// BadRequestParam returns a new bad request error with a parameter source.
func BadRequestParam(detail, param string) *Error {
	err := ErrorFromStatus(http.StatusBadRequest, detail)
	err.Source = &ErrorSource{
		Parameter: param,
	}

	return err
}

// BadRequestPointer returns a new bad request error with a pointer source.
func BadRequestPointer(detail, pointer string) *Error {
	err := ErrorFromStatus(http.StatusBadRequest, detail)
	err.Source = &ErrorSource{
		Pointer: pointer,
	}

	return err
}

// InternalServerError returns na new internal server error.
func InternalServerError(detail string) *Error {
	return ErrorFromStatus(http.StatusInternalServerError, detail)
}
