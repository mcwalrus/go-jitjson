package main

import (
	"encoding/json"
	"fmt"
)

// TODO need to handle
type JitJSON[T JSON] struct {
	data []byte
	val  *T
}

// TODO rename.
type JSON interface {
	json.Marshaler
	json.Unmarshaler
}

func NewJitJSON[T JSON](data []byte) (*JitJSON[T], error) {
	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid json")
	}

	jit := JitJSON[T]{
		data: data,
	}

	return &jit, nil
}

// TODO: pretty sure
func (jit *JitJSON[JSON]) Validate() error {
	if jit.data == nil && jit.val == nil {
		return fmt.Errorf("jit-json missing value")
	}
	return nil
}

func (jit *JitJSON[JSON]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}

	var err error
	jit.data, err = json.Marshal(&jit.val)
	if err != nil {
		return nil, err
	}

	return jit.data, nil
}

func (jit *JitJSON[JSON]) Unmarshal() (*JSON, error) {
	if jit.val != nil {
		return jit.val, nil
	}

	var err error
	err = json.Unmarshal(jit.data, &jit.val)
	if err != nil {
		return nil, err
	}

	return jit.val, nil
}
