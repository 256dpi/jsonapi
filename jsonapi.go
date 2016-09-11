// Package jsonapi provides structures and methods to implement JSON API
// compatible APIs. Most methods are tailored to be used together with the
// echo framework, yet all of them also have a native counterpart in the compat
// sub package that allows implementing APIs using the standard HTTP library.
package jsonapi

// MediaType is the official JSON API media type that should be used by
// all requests and responses.
const MediaType = "application/vnd.api+json"

// Map is a general purpose map of string keys and arbitrary values.
type Map map[string]interface{}
