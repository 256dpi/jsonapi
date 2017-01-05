package jsonapi

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	var test struct {
		Foo string
	}

	test.Foo = "foo"

	m, err := StructToMap(&test, nil)
	assert.NoError(t, err)
	assert.Equal(t, Map{
		"Foo": "foo",
	}, m)
}

func TestStructToMapWithTag(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	test.Foo = "foo"

	m, err := StructToMap(&test, nil)
	assert.NoError(t, err)
	assert.Equal(t, Map{
		"bar": "foo",
	}, m)
}

func TestStructToMapOmitEmpty(t *testing.T) {
	var test struct {
		Foo string `json:",omitempty"`
	}

	m, err := StructToMap(&test, nil)
	assert.NoError(t, err)
	assert.Equal(t, Map{}, m)
}

func TestStructToMapInvalidStruct(t *testing.T) {
	var test struct {
		Foo func()
	}

	m, err := StructToMap(&test, nil)
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestStructToMapInvalidStruct2(t *testing.T) {
	test := "foo"

	m, err := StructToMap(&test, nil)
	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestStructToMapFiltering(t *testing.T) {
	var test struct {
		Foo string `json:"bar"`
	}

	test.Foo = "foo"

	m, err := StructToMap(&test, []string{"baz"})
	assert.NoError(t, err)
	assert.Equal(t, Map{}, m)

	m, err = StructToMap(&test, []string{"bar"})
	assert.NoError(t, err)
	assert.Equal(t, Map{
		"bar": "foo",
	}, m)
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

func TestMapAssignInvalidMap(t *testing.T) {
	var test struct {
		Foo string
	}

	m := Map{"foo": func() {}}

	err := m.Assign(&test)
	assert.Error(t, err)
	assert.Equal(t, "", test.Foo)
}

func TestBigNumberConversion(t *testing.T) {
	type test struct {
		Int8    int8
		Int64   int64
		UInt8   uint8
		UInt64  uint64
		Float32 float32
		Float64 float64
	}

	i := &test{
		Int8:    math.MaxInt8,
		Int64:   math.MaxInt64,
		UInt8:   math.MaxUint8,
		UInt64:  math.MaxUint64,
		Float32: math.MaxFloat32,
		Float64: math.MaxFloat64,
	}

	m, err := StructToMap(i, nil)
	assert.NoError(t, err)
	assert.NotNil(t, m)

	ii := &test{}

	err = m.Assign(ii)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(i, ii))
}

func BenchmarkStructToMap(b *testing.B) {
	var test struct {
		Foo string
	}

	test.Foo = "foo"

	for i := 0; i < b.N; i++ {
		_, err := StructToMap(&test, nil)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkMapAssign(b *testing.B) {
	var test struct {
		Foo string
	}

	m := Map{"foo": "foo"}

	for i := 0; i < b.N; i++ {
		err := m.Assign(&test)
		if err != nil {
			panic(err)
		}
	}
}
