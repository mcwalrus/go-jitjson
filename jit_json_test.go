package jitjson_test

import (
	"encoding/json"
	"testing"

	"github.com/MaxCollier/go-jitjson"
	"github.com/stretchr/testify/require"
)

type person interface {
	Name() string
	json.Marshaler
	json.Unmarshaler
}

type Person struct {
	Name_ string `json:"name"`
}

func (p Person) Name() string {
	return p.Name_
}

func TestJitJSON(t *testing.T) {

	// Since the underlying JSON type could be a pointer to some type, the representation
	// is ignored when the value is set as it could become out of sync with the value.
	t.Run("check interface", func(t *testing.T) {
		pJIT, err := jitjson.NewJitJSON[person]([]byte{})
		require.NoError(t, err)

		type Person struct {
			Name_ string `json:"name"`
		}

	})
}
