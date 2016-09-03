// Package jsonapi provides structure ans methods to implement JSON API
// compatible APIs.
package jsonapi

var objectSuffix = []byte("{")
var arraySuffix = []byte("[")

// MediaType is the official JSON API media type that should be used by
// all requests and responses.
const MediaType = "application/vnd.api+json"
