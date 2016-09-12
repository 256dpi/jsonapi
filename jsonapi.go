// Package jsonapi provides structures and methods to implement JSON API
// compatible APIs. Most methods are tailored to be used together with the
// echo framework, yet all of them also have a native counterpart in the compat
// sub package that allows implementing APIs using the standard HTTP library.
package jsonapi

import (
	"io"

	"github.com/fatih/structs"
)

// MediaType is the official JSON API media type that should be used by
// all requests and responses.
const MediaType = "application/vnd.api+json"

// The Requester interface must be implemented by adapters to make framework
// specific request objects compatible with jsonapi.
type Requester interface {
	Method() string
	Get(key string) string
	Path() string
	QueryParams() map[string][]string
}

// The Responder interface must be implemented by adapters to make framework
// specific response objects compatible with jsonapi.
type Responder interface {
	io.Writer

	Set(key, value string)
	WriteHeader(status int)
}

// Map is a general purpose map of string keys and arbitrary values.
type Map map[string]interface{}

// StructToMap  will assign the fields of the source struct to a new map.
//
// Note: The "json" tag will be respected to write proper field names.
func StructToMap(source interface{}) Map {
	s := structs.New(source)
	s.TagName = "json"
	return Map(s.Map())
}
