// Package jsonapi provides structure ans methods to implement JSON API
// compatible APIs.
package jsonapi

var objectSuffix = []byte("{")
var arraySuffix = []byte("[")

// ContentType is the official JSON API content type that should be used by
// all requests and responses.
const ContentType = "application/vnd.api+json"
