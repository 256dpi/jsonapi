// Package jsonapi provides structure and methods to implement JSON API
// compatible APIs. Most methods are tailored to be used together with the
// echo framework, yet all of them also have a native counterpart that allows
// implementing APIs using the standard library.
package jsonapi

import "github.com/mitchellh/mapstructure"

// MediaType is the official JSON API media type that should be used by
// all requests and responses.
const MediaType = "application/vnd.api+json"

// MapToStruct wraps the mapstructure package to provide a simple way to assign
// map values to a struct. The optional parameter tag name defaults to "json".
func MapToStruct(m interface{}, s interface{}, tagName ...string) error {
	// read optional tag name
	_tagName := "json"
	if len(tagName) > 0 {
		_tagName = tagName[0]
	}

	// prepare decoder config
	cfg := &mapstructure.DecoderConfig{
		ZeroFields: true,
		Result:     s,
		TagName:    _tagName,
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
