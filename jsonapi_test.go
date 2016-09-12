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

func TestMapAssignInvalidTarget(t *testing.T) {
	m := Map{"foo": "foo"}

	err := m.Assign(nil)
	assert.Error(t, err)
}

func TestMapAssign(t *testing.T) {
	var test struct {
		Foo string
	}

	m := Map{"foo": "foo"}

	err := m.Assign(&test)
	assert.NoError(t, err)
	assert.Equal(t, "foo", test.Foo)
}

func TestMapAssignWithTag(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	m := Map{"bar": "foo"}

	err := m.Assign(&test)
	assert.NoError(t, err)
	assert.Equal(t, "foo", test.Foo)
}

func TestMapAssignUnmatchedFields(t *testing.T) {
	var test struct {
		Foo string
	}

	m := Map{"bar": "foo"}

	err := m.Assign(&test)
	assert.NoError(t, err)
	assert.Equal(t, "", test.Foo)
}

func TestMapAssignInvalidType(t *testing.T) {
	var test struct {
		Foo string
	}

	m := Map{"foo": 1}

	err := m.Assign(&test)
	assert.Error(t, err)
	assert.Equal(t, "", test.Foo)
}
