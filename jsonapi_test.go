package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	var test struct {
		Foo string
	}

	test.Foo = "foo"

	assert.Equal(t, Map{
		"Foo": "foo",
	}, StructToMap(&test))
}

func TestStructToMapWithTag(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	test.Foo = "foo"

	assert.Equal(t, Map{
		"bar": "foo",
	}, StructToMap(&test))
}

func TestStructToMapOmitEmpty(t *testing.T) {
	var test struct {
		Foo string `json:",omitempty"`
	}

	assert.Equal(t, Map{}, StructToMap(&test))
}
