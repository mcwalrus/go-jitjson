package jitjson_test

import (
	"encoding/json"
	"testing"

	"github.com/mcwalrus/go-jitjson"
)

type customType struct {
	Word   string
	Number int
}

var (
	myStruct = customType{
		Word:   word,
		Number: number,
	}
)

const (
	word   = "ghfi;dsafldskafhdsafdslafhiadsfds;ajhfidlasafdnverfaieldadjfaisdalenaciiafjds"
	number = 901901199010191
)

// BenchmarkMarshalJSON marshals by encoding/json.
func BenchmarkMarshalJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {
			_, err := json.Marshal(&myStruct)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkMarshalJitJSON1 marshals all by jit-json.
func BenchmarkMarshalJitJSON1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {

			jit, err := jitjson.NewJitJSON[customType](myStruct)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}

			_, err = jit.Marshal()
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkMarshalJitJSON2 marshals half by jit-json.
func BenchmarkMarshalJitJSON2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {

			jit, err := jitjson.NewJitJSON[customType](myStruct)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}

			if i%2 == 0 {
				continue
			}

			_, err = jit.Marshal()
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkMarshalJitJSON3 marshals none by jit-json.
func BenchmarkMarshalJitJSON3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {
			_, err := jitjson.NewJitJSON[customType](myStruct)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkUnmarshalJSON un-marshals by encoding/json.
func BenchmarkUnmarshalJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {
			_, err := json.Marshal(&myStruct)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkUnmarshalJitJSON1 un-marshals all by jit-json.
func BenchmarkUnmarshalJitJSON1(b *testing.B) {
	data, err := json.Marshal(myStruct)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {

			jit, err := jitjson.NewJitJSON[customType](data)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}

			_, err = jit.Unmarshal()
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkUnmarshalJitJSON2 un-marshals half by jit-json.
func BenchmarkUnmarshalJitJSON2(b *testing.B) {
	data, err := json.Marshal(myStruct)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {

			jit, err := jitjson.NewJitJSON[customType](data)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}

			if i%2 == 0 {
				continue
			}

			_, err = jit.Unmarshal()
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}

// BenchmarkUnmarshalJitJSON3 un-marshals none by jit-json.
func BenchmarkUnmarshalJitJSON3(b *testing.B) {
	data, err := json.Marshal(myStruct)
	if err != nil {
		b.Error(err)
		b.FailNow()
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {

			_, err := jitjson.NewJitJSON[customType](data)
			if err != nil {
				b.Error(err)
				b.FailNow()
			}
		}
	}
}
