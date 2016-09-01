package jsonapi

import "fmt"

// See: http://jsonapi.org/format/#errors.
type ErrorPayload struct {
	// A list of errors that occurred during the request.
	List []*Error `json:"errors"`

	// Non-standard meta-information about the payload.
	Meta Map `json:"meta,omitempty"`
}

// See: http://jsonapi.org/format/#errors.
type ErrorLinks struct {
	// A link that leads to further details about this particular occurrence of the problem.
	About string `json:"about"`
}

// See: http://jsonapi.org/format/#errors.
type ErrorSource struct {
	// A string indicating which URI query parameter caused the error.
	Parameter string `json:"parameter,omitempty"`

	// A JSON Pointer to the associated entity in the request document.
	Pointer string `json:"pointer,omitempty"`
}

// TODO: Make Title and Detail mandatory?
// TODO: Use opaque int type for status?
// TODO: Use error type for detail or title?

// See: http://jsonapi.org/format/#errors.
type Error struct {
	// A unique identifier for this particular occurrence of the problem.
	ID string `json:"id,omitempty"`

	// Continuing links to other resources.
	Links *ErrorLinks `json:"links,omitempty"`

	// The HTTP status code applicable to this problem.
	Status string `json:"status,omitempty"`

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

// Error returns a string representation of the error for logging purposes.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}
