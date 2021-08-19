// Package jsonapi provides structures and functions to implement JSON API
// compatible APIs. The library can be used with any framework and is built on
// top of the standard Go http library.
package jsonapi

import (
	"bytes"
	"encoding/json"
)

// MediaType is the official JSON API media type that should be used by
// all requests and responses.
const MediaType = "application/vnd.api+json"

// Map is a general purpose map of string keys and arbitrary values.
//
// Note: All methods in this package will leave numbers as strings to avoid
// issues with mismatching types when they are later assigned to a struct.
type Map map[string]interface{}

// StructToMap will assign the fields of the source struct to a new map and
// additionally filter the map to only include the fields specified.
//
// Note: The "json" tag will be respected to write proper field names. No
// filtering will be applied if fields is nil.
//
// Note: Numbers are left as strings to avoid issues with mismatching types
// when they are later assigned to a struct again.
//
// Warning: The function does actually convert the struct to json and then
// convert that json to a map. High performance applications might want to use
// a custom implementation that is much faster.
func StructToMap(source interface{}, fields []string) (Map, error) {
	// marshal struct as json
	buf, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	// prepare decoder
	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.UseNumber()

	// unmarshal json to map
	var m Map
	err = dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	// return map directly if no fields are specified
	if fields == nil {
		return m, nil
	}

	// filter map
	for key := range m {
		// check if field is present
		ok := false
		for _, field := range fields {
			if field == key {
				ok = true
			}
		}

		// otherwise, remove field
		if !ok {
			delete(m, key)
		}
	}

	return m, nil
}

// Assign will assign the values in the map to the target struct.
//
// Note: The "json" tag will be respected to match field names.
//
// Warning: The function does actually convert the map to json and then assign
// that json to the struct. High performance applications might want to use a
// custom implementation that is much faster.
func (m Map) Assign(target interface{}) error {
	// marshal map to json
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// unmarshal json to struct
	err = json.Unmarshal(buf, target)
	if err != nil {
		return err
	}

	return nil
}
