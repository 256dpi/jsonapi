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
	}, StructToMap(&test, nil))
}

func TestStructToMapWithTag(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	test.Foo = "foo"

	assert.Equal(t, Map{
		"bar": "foo",
	}, StructToMap(&test, nil))
}

func TestStructToMapOmitEmpty(t *testing.T) {
	var test struct {
		Foo string `json:",omitempty"`
	}

	assert.Equal(t, Map{}, StructToMap(&test, nil))
}

func TestStructToMapFiltering(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	test.Foo = "foo"

	assert.Equal(t, Map{}, StructToMap(&test, []string{"baz"}))
	assert.Equal(t, Map{
		"bar": "foo",
	}, StructToMap(&test, []string{"bar"}))
}
