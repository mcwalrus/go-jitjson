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
			want: []*AnyJitJSON{
				nil,
				{NewJitJSON(true)},
				{NewJitJSON(json.Number("123"))},
				{NewJitJSON("hello")},
			},
			wantErr: false,
		},
		{
			name:  "object value",
			input: `{"key1": null, "key2": true, "key3": 123, "key4": "hello"}`,
			want: map[string]*AnyJitJSON{
				"key1": nil,
				"key2": {NewJitJSON(true)},
				"key3": {NewJitJSON(json.Number("123"))},
				"key4": {NewJitJSON("hello")},
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
			if !tt.wantErr && !compareValues(t, a.v, tt.want) {
				t.Errorf("UnmarshalJSON() = %v, want %v", a.v, tt.want)
			}
		})
	}
}

func TestAnyJit_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   *AnyJitJSON
		want    string
		wantErr bool
	}{
		{
			name:    "null value",
			input:   &AnyJitJSON{nil},
			want:    "null",
			wantErr: false,
		},
		{
			name:    "boolean true",
			input:   &AnyJitJSON{NewJitJSON(true)},
			want:    "true",
			wantErr: false,
		},
		{
			name:    "boolean false",
			input:   &AnyJitJSON{NewJitJSON(false)},
			want:    "false",
			wantErr: false,
		},
		{
			name:    "number value",
			input:   &AnyJitJSON{NewJitJSON(json.Number("123.45"))},
			want:    "123.45",
			wantErr: false,
		},
		{
			name:    "string value",
			input:   &AnyJitJSON{NewJitJSON("hello")},
			want:    `"hello"`,
			wantErr: false,
		},
		{
			name: "array value",
			input: &AnyJitJSON{[]*AnyJitJSON{
				nil,
				{NewJitJSON(true)},
				{NewJitJSON(json.Number("123"))},
				{NewJitJSON("hello")},
			}},
			want:    `[null,true,123,"hello"]`,
			wantErr: false,
		},
		{
			name: "object value",
			input: &AnyJitJSON{map[string]*AnyJitJSON{
				"key1": nil,
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
			got, err := json.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.input.Type() != TypeObject {
				if string(got) != tt.want {
					t.Errorf("MarshalJSON() = %s, want %s", got, tt.want)
				}
			}
		})
	}
}

func TestAnyJitJSON_IsNull(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{
			name: "null value",
			data: []byte("null"),
			want: true,
		},
		{
			name: "non-null value",
			data: []byte("true"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AnyJitJSON
			if err := json.Unmarshal(tt.data, &a); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			if got := a.IsNull(); got != tt.want {
				t.Errorf("IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyJitJSON_AsBool(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		want   bool
		wantOk bool
	}{
		{
			name:   "boolean true",
			data:   []byte("true"),
			want:   true,
			wantOk: true,
		},
		{
			name:   "boolean false",
			data:   []byte("false"),
			want:   false,
			wantOk: true,
		},
		{
			name:   "non-boolean value",
			data:   []byte(`"string"`),
			want:   false,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AnyJitJSON
			if err := json.Unmarshal(tt.data, &a); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			got, ok := a.AsBool()
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("AsBool() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestAnyJitJSON_AsNumber(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		want   json.Number
		wantOk bool
	}{
		{
			name:   "number value",
			data:   []byte("123.45"),
			want:   json.Number("123.45"),
			wantOk: true,
		},
		{
			name:   "non-number value",
			data:   []byte(`"string"`),
			want:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a AnyJitJSON
			if err := json.Unmarshal(tt.data, &a); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			got, ok := a.AsNumber()
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("AsNumber() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestAnyJitJSON_AsString(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		want   string
		wantOk bool
	}{
		{
			name:   "null value",
			data:   []byte(`null`),
			want:   "",
			wantOk: false,
		},
		{
			name:   "string value",
			data:   []byte(`"hello"`),
			want:   "hello",
			wantOk: true,
		},
		{
			name:   "boolean is also string value",
			data:   []byte(`"true"`),
			want:   "true",
			wantOk: true,
		},
		{
			name:   "number is also string value",
			data:   []byte(`"123.45"`),
			want:   `"123.45"`,
			wantOk: true,
		},
		{
			name:   "array text is also string value",
			data:   []byte(`"[1,2,3]"`),
			want:   "[1,2,3]",
			wantOk: true,
		},
		{
			name:   "object text is also string value",
			data:   []byte(`{"key": "value"}`),
			want:   `{"key": "value"}`,
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			var a AnyJitJSON
			if err := json.Unmarshal(tt.data, &a); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}
			got, ok := a.AsString()
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("AsString() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestAnyJitJSON_AsArray(t *testing.T) {
	tests := []struct {
		name   string
		a      *AnyJitJSON
		want   []*AnyJitJSON
		wantOk bool
	}{
		{
			name:   "array value",
			a:      &AnyJitJSON{[]*AnyJitJSON{NewAnyJitJSON([]byte("null")), NewAnyJitJSON([]byte("true"))}},
			want:   []*AnyJitJSON{NewAnyJitJSON([]byte("null")), NewAnyJitJSON([]byte("true"))},
			wantOk: true,
		},
		{
			name:   "non-array value",
			a:      &AnyJitJSON{NewJitJSON("string")},
			want:   nil,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.a.AsArray()
			if !compareAnyJitJSONArray(t, got, tt.want) || ok != tt.wantOk {
				t.Errorf("AsArray() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestAnyJitJSON_AsObject(t *testing.T) {
	tests := []struct {
		name   string
		a      *AnyJitJSON
		want   map[string]*AnyJitJSON
		wantOk bool
	}{
		{
			name:   "object value",
			a:      &AnyJitJSON{map[string]*AnyJitJSON{"key": NewAnyJitJSON([]byte("true"))}},
			want:   map[string]*AnyJitJSON{"key": NewAnyJitJSON([]byte("true"))},
			wantOk: true,
		},
		{
			name:   "non-object value",
			a:      &AnyJitJSON{NewJitJSON("string")},
			want:   nil,
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.a.AsObject()
			if !compareAnyJitJSONObject(t, got, tt.want) || ok != tt.wantOk {
				t.Errorf("AsObject() = %v, %v, want %v, %v", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func compareAnyJitJSONArray(t *testing.T, got, want []*AnyJitJSON) bool {
	t.Helper()
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if !compareAnyJitJSON(t, got[i], want[i]) {
			return false
		}
	}
	return true
}

func compareAnyJitJSONObject(t *testing.T, got, want map[string]*AnyJitJSON) bool {
	t.Helper()
	if len(got) != len(want) {
		return false
	}
	for k, v := range got {
		if !compareAnyJitJSON(t, v, want[k]) {
			return false
		}
	}
	return true
}

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
	case []*AnyJitJSON:
		want, ok := want.([]*AnyJitJSON)
		if !ok || len(got) != len(want) {
			t.Log("slice does not match")
			return false
		}
		for i := range got {
			if !compareAnyJitJSON(t, got[i], want[i]) {
				return false
			}
		}
		return true
	case map[string]*AnyJitJSON:
		want, ok := want.(map[string]*AnyJitJSON)
		if !ok || len(got) != len(want) {
			t.Log("map does not match")
			return false
		}
		for k, v := range got {
			wantV, ok := want[k]
			if !ok {
				t.Errorf("unexpected key: %s", k)
				return false
			}
			if !compareAnyJitJSON(t, v, wantV) {
				return false
			}
		}
		return true
	default:
		return got == want
	}
}

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

func compareAnyJitJSON(t *testing.T, got *AnyJitJSON, want *AnyJitJSON) bool {
	t.Helper()

	if (got == nil) != (want == nil) {
		t.Error("unexpected value")
		return false
	}
	if got == nil {
		return true
	}

	if (got.v == nil) != (want.v == nil) {
		t.Error("unexpected value")
		return false
	}
	if got.v == nil {
		return true
	}

	return compareValues(t, got.v, want.v)
}
