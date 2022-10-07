package main

import (
	"encoding/json"
	"fmt"
)

// JSON struct requires marshal and unmarshal json methods or json tags.
type JSON interface {
	json.Marshaler
	json.Unmarshaler
}

// JitJSON provides 'just-in-time' complication for both marshalling and
// unmarshalling some JSON struct.
type JitJSON[T JSON] struct {
	data []byte
	val  *T
}

// NewJitJSON creates new JitJSON from data.
func NewJitJSON[T JSON](data []byte) (*JitJSON[T], error) {
	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid json")
	}

	jit := JitJSON[T]{
		data: data,
	}

	return &jit, nil
}

// Set new value to jitJSON.
// JitJSON byte representation is wiped to avoid data mismatch.
func (jit *JitJSON[JSON]) Set(val JSON) {
	jit.data = nil
	jit.val = &val
}

// Marshal provides the wrapper function to json.Marshal.
// JitJSON caches the byte representation once marshal has occurred.
func (jit *JitJSON[JSON]) Marshal() ([]byte, error) {
	if jit.data != nil {
		return jit.data, nil
	}

	var err error
	jit.data, err = json.Marshal(jit.val)
	if err != nil {
		return nil, err
	}

	return jit.data, nil
}

// Unmarshal provides the wrapper function to json.Unmarshal.
// JitJSON caches the returned value once unmarshal has occurred.
func (jit *JitJSON[JSON]) Unmarshal() (JSON, error) {
	if jit.val != nil {
		return *jit.val, nil
	}

	var val JSON
	err := json.Unmarshal(jit.data, jit.val)
	if err != nil {
		return val, err
	}

	return *jit.val, nil
}
