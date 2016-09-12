// Package jsonapi provides structures and methods to implement JSON API
// compatible APIs. Most methods are tailored to be used together with the
// echo framework, yet all of them also have a native counterpart in the compat
// sub package that allows implementing APIs using the standard HTTP library.
package jsonapi

import (
	"io"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
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

// StructToMap will assign the fields of the source struct to a new map and
// additionally filter the map to only include the fields specified.
//
// Note: The "json" tag will be respected to write proper field names.
//
// Note: If no fields to filter are specified, not filtering will be applied.
func StructToMap(source interface{}, fields []string) Map {
	// prepare structs helper
	s := structs.New(source)
	s.TagName = "json"

	// create map
	m := Map(s.Map())

	// return map directly of no fields are specified
	if len(fields) == 0 {
		return m
	}

	// filter map
	for key := range m {
		ok := false

		// check if field is present
		for _, field := range fields {
			if field == key {
				ok = true
			}
		}

		// otherwise remove field
		if !ok {
			delete(m, key)
		}
	}

	return m
}

// Assign will assign the values in the map to the target struct.
//
// Note: The "json" tag will be respected to match field names.
func (m Map) Assign(target interface{}) error {
	// prepare decoder config
	cfg := &mapstructure.DecoderConfig{
		ZeroFields: true,
		Result:     target,
		TagName:    "json",
	}

	// create decoder
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}

	// run decoder
	err = decoder.Decode(m)
	if err != nil {
		return BadRequest(err.Error())
	}

	return nil
}
