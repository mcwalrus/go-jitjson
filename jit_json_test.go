package jitjson_test

import (
	"testing"

	"github.com/MaxCollier/go-jitjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: remove the assert / require framework to make bundle smaller.
func TestNilJitJSON(t *testing.T) {
	data := []byte(`null`)

	jit, err := jitjson.NewJitJSON[struct{}](data)
	require.NoError(t, err)

	value, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, struct{}{}, value)
}

func TestIntJitJSON(t *testing.T) {
	data := []byte(`1`)

	jit, err := jitjson.NewJitJSON[int](data)
	require.NoError(t, err)

	value, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, 1, value)
}

func TestFloatJitJSON(t *testing.T) {
	data := []byte(`1.00000001`)

	jit, err := jitjson.NewJitJSON[float64](data)
	require.NoError(t, err)

	value, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, 1.00000001, value)
}

func TestStringJitJSON(t *testing.T) {
	data := []byte(`"some text"`)

	jit, err := jitjson.NewJitJSON[string](data)
	require.NoError(t, err)

	value, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, "some text", value)
}

func TestArrayJitJSON(t *testing.T) {
	data := []byte(`[1, 2, 3, 4, 5]`)

	jit, err := jitjson.NewJitJSON[[]int](data)
	require.NoError(t, err)

	value, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2, 3, 4, 5}, value)
}

func TestMapJitJSON(t *testing.T) {
	data := []byte(`{
		"1": "two",
		"three": 4.0,
		"5": [6, 7, "eight", 9, "ten"]
	}`)

	jit, err := jitjson.NewJitJSON[map[string]interface{}](data)
	require.NoError(t, err)

	m, err := jit.Unmarshal()
	require.NoError(t, err)
	require.NotNil(t, m)

	value, ok := m["1"]
	assert.True(t, ok)
	assert.Equal(t, "two", value)

	value, ok = m["three"]
	assert.True(t, ok)
	assert.Equal(t, 4.0, value)

	value, ok = m["5"]
	assert.True(t, ok)
	assert.Equal(t, []interface{}{6.0, 7.0, "eight", 9.0, "ten"}, value)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestExample1(t *testing.T) {
	data := []byte(`
        {
            "name": "Willy Wonka",
            "age":  42
        }
    `)

	expected := Person{
		Name: "Willy Wonka",
		Age:  42,
	}

	jit, err := jitjson.NewJitJSON[Person](data)
	require.NoError(t, err)

	person, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, expected, person)
}

func TestExample2(t *testing.T) {
	person := Person{
		Name: "Charlie Bucket",
		Age:  12,
	}

	expected := []byte(`{"name":"Charlie Bucket","age":12}`)

	jit, err := jitjson.NewJitJSON[Person](nil)
	require.NoError(t, err)

	jit.Set(person)
	data, err := jit.Marshal()
	require.NoError(t, err)

	assert.Equal(t, expected, data)
}

func TestJitJSON1(t *testing.T) {
	data := []byte(`{
		"name": "Chris",
		"age": 42
	}`)

	jit, err := jitjson.NewJitJSON[Person](data)
	require.NoError(t, err)

	person, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, "Chris", person.Name)
	assert.Equal(t, 42, person.Age)
}

type Boater struct {
	Person Person  `json:"person"`
	Boat   *string `json:"boat"`
}

func TestJitJSON2(t *testing.T) {
	data := []byte(`{
		"person": {
			"name": "Chris",
			"age": 42
		},
		"boat": "Stout"
	}`)

	jit, err := jitjson.NewJitJSON[Boater](data)
	require.NoError(t, err)

	boater, err := jit.Unmarshal()
	require.NoError(t, err)

	assert.Equal(t, "Chris", boater.Person.Name)
	assert.Equal(t, 42, boater.Person.Age)
	assert.Nil(t, boater.Boat)
}

func TestJitJSON3(t *testing.T) {
	data := []byte(`{
		"person": {
			"name": "Chris",
			"age": 42
		},
		"boat": "Stout"
	}`)

	jit, err := jitjson.NewJitJSON[Boater](data)
	require.NoError(t, err)

	boater, err := jit.Unmarshal()
	require.NoError(t, err)

	if assert.NotNil(t, boater.Boat) {
		assert.Equal(t, "Stout", *boater.Boat)
	}

	if t.Failed() {
		return
	}

	newBoat := "Steeze"
	boater.Boat = &newBoat
	jit.Set(boater)

	boater, err = jit.Unmarshal()
	require.NoError(t, err)

	if assert.NotNil(t, boater.Boat) {
		assert.Equal(t, "Steeze", *boater.Boat)
	}
}

type testInterface interface {
	method() int
}

func TestErrorsJitJSON(t *testing.T) {

	// triggers pointer type error.
	t.Run("pointer", func(t *testing.T) {
		_, err := jitjson.NewJitJSON[*int]([]byte("null"))
		require.Error(t, err)
	})

	// triggers invalid type error.
	t.Run("interface", func(t *testing.T) {
		_, err := jitjson.NewJitJSON[interface{}]([]byte("null"))
		assert.Error(t, err)
	})

	// triggers invalid type error.
	t.Run("invalid type", func(t *testing.T) {
		_, err := jitjson.NewJitJSON[testInterface]([]byte("null"))
		assert.Error(t, err)
	})
}
