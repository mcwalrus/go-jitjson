package jitjson

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestAnyJit_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		wantErr bool
	}{
		{
			name:    "null value",
			input:   "null",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "boolean true",
			input:   "true",
			want:    NewJitJSON(true),
			wantErr: false,
		},
		{
			name:    "boolean false",
			input:   "false",
			want:    NewJitJSON(false),
			wantErr: false,
		},
		{
			name:    "number value",
			input:   "123.45",
			want:    NewJitJSON(json.Number("123.45")),
			wantErr: false,
		},
		{
			name:    "string value",
			input:   `"hello"`,
			want:    NewJitJSON("hello"),
			wantErr: false,
		},
		{
			name:  "array value",
			input: `[null, true, 123, "hello"]`,
			want: []AnyJitJSON{
				{nil},
				{true},
				{json.Number("123")},
				{"hello"},
			},
			wantErr: false,
		},
		{
			name:  "object value",
			input: `{"key1": null, "key2": true, "key3": 123, "key4": "hello"}`,
			want: map[string]AnyJitJSON{
				"key1": {nil},
				"key2": {true},
				"key3": {json.Number("123")},
				"key4": {"hello"},
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `invalid`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			var a AnyJitJSON
			err := json.Unmarshal([]byte(tt.input), &a)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareValues(t, a.Value(), tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", a.Value(), tt.want)
			}
		})
	}
}

func TestAnyJit_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   AnyJitJSON
		want    string
		wantErr bool
	}{
		{
			name:    "null value",
			input:   AnyJitJSON{nil},
			want:    "null",
			wantErr: false,
		},
		{
			name:    "boolean true",
			input:   AnyJitJSON{NewJitJSON(true)},
			want:    "true",
			wantErr: false,
		},
		{
			name:    "boolean false",
			input:   AnyJitJSON{NewJitJSON(false)},
			want:    "false",
			wantErr: false,
		},
		{
			name:    "number value",
			input:   AnyJitJSON{NewJitJSON(json.Number("123.45"))},
			want:    "123.45",
			wantErr: false,
		},
		{
			name:    "string value",
			input:   AnyJitJSON{NewJitJSON("hello")},
			want:    `"hello"`,
			wantErr: false,
		},
		{
			name: "array value",
			input: AnyJitJSON{[]AnyJitJSON{
				{nil},
				{NewJitJSON(true)},
				{NewJitJSON(json.Number("123"))},
				{NewJitJSON("hello")},
			}},
			want:    `[null,true,123,"hello"]`,
			wantErr: false,
		},
		{
			name: "object value",
			input: AnyJitJSON{map[string]AnyJitJSON{
				"key1": {nil},
				"key2": {NewJitJSON(true)},
				"key3": {NewJitJSON(json.Number("123"))},
				"key4": {NewJitJSON("hello")},
			}},
			want:    `{"key1":null,"key2":true,"key3":123,"key4":"hello"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			// var anyJit AnyJitJSON
			got, err := json.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

// compareValues is a helper function to compare the values of AnyJitJSON.
func compareValues(t *testing.T, got, want interface{}) bool {
	t.Helper()

	switch got := got.(type) {
	case *JitJSON[bool]:
		return compareJitJSON(t, got, want)
	case *JitJSON[json.Number]:
		return compareJitJSON(t, got, want)
	case *JitJSON[string]:
		return compareJitJSON(t, got, want)
	case *JitJSON[AnyJitJSON]:
		return compareJitJSON(t, got, want)
	case []AnyJitJSON:
		// want, ok := want.([]*JitJSON[AnyJitJSON])
		// if !ok || len(got) != len(want) {
		// 	t.Log("slice does not match")
		// 	return false
		// }
		// for i := range got {
		// 	jitAny, err := got[i].Unmarshal()
		// 	if err != nil {
		// 		return false
		// 	}
		// 	wantJitAny, err := want[i].Unmarshal()
		// 	if err != nil {
		// 		return false
		// 	}
		// 	if !compareValues(t, jitAny.v, wantJitAny.v) {
		// 		return false
		// 	}
		// }
		return true
	case map[string]AnyJitJSON:
		// want, ok := want.(map[string]*JitJSON[AnyJitJSON])
		// if !ok || len(got) != len(want) {
		// 	t.Log("map does not match")
		// 	return false
		// }
		// for k, v := range got {
		// 	wantV, ok := want[k]
		// 	if !ok {
		// 		return false
		// 	}
		// 	jitAny, err := v.Unmarshal()
		// 	if err != nil {
		// 		return false
		// 	}
		// 	wantJitAny, err := wantV.Unmarshal()
		// 	if err != nil {
		// 		return false
		// 	}
		// 	if !compareValues(t, jitAny.v, wantJitAny.v) {
		// 		return false
		// 	}
		// }
		return true
	default:
		return got == want
	}
}

// compareJitJSON is a helper function to compare the values of JitJSON. Input 'w' (want)
// be type asserted to a JitJSON type, and the values will be compared.
// The argument 'got' should be a JitJSON type, without the value set by deferred parsing.
// Once checked, we can compare the data and value between the two JitJSON types.
func compareJitJSON[T comparable](t *testing.T, got *JitJSON[T], w interface{}) bool {
	t.Helper()
	var err error

	want, ok := w.(*JitJSON[T])
	if !ok || want == nil {
		t.Error("unexpected type")
		return false
	}

	if got.val != nil {
		t.Error("val should be nil")
		return false
	}

	_, err = want.Unmarshal()
	if err != nil {
		t.Error(err)
		return false
	}

	_, err = got.Unmarshal()
	if err != nil {
		t.Error(err)
		return false
	}

	if got.val != got.val {
		t.Error("values do not match")
		return false
	}

	_, err = got.Marshal()
	if err != nil {
		t.Error(err)
		return false
	}

	_, err = want.Marshal()
	if err != nil {
		t.Error(err)
		return false
	}

	if !bytes.Equal(got.data, want.data) {
		t.Error("data do not match")
		return false
	}

	return true
}
