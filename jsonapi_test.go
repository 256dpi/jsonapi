package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToStructInvalidTarget(t *testing.T) {
	err := MapToStruct(Map{
		"foo": "foo",
	}, nil)

	assert.Error(t, err)
}

func TestMapToStruct(t *testing.T) {
	var test1 struct {
		Foo string
	}

	err := MapToStruct(Map{
		"foo": "foo",
	}, &test1)

	assert.NoError(t, err)
	assert.Equal(t, "foo", test1.Foo)
}

func TestMapToStructUnmatchedFields(t *testing.T) {
	var test1 struct {
		Foo string
	}

	err := MapToStruct(Map{
		"bar": "foo",
	}, &test1)

	assert.NoError(t, err)
	assert.Equal(t, "", test1.Foo)
}

func TestMapToStructInvalidType(t *testing.T) {
	var test1 struct {
		Foo string
	}

	err := MapToStruct(Map{
		"foo": 1,
	}, &test1)

	assert.Error(t, err)
	assert.Equal(t, "", test1.Foo)
}

func TestMapToStructTagName(t *testing.T) {
	var test1 struct {
		Foo string `foo:"bar"`
	}

	err := MapToStruct(Map{
		"bar": "baz",
	}, &test1, "foo")

	assert.NoError(t, err)
	assert.Equal(t, "baz", test1.Foo)
}
