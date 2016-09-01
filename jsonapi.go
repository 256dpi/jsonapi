// Package jsonapi provides structure ans methods to implement JSON API
// compatible APIs.
package jsonapi

import "net/http"

// ContentType is the official JSON API content type that should be used by
// all requests and responses.
const ContentType = "application/vnd.api+json"

// SetContentType will set the proper JSON API content type on the passed
// response writer.
func SetContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentType)
}
