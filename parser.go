package jitjson

// JSONParser is an interface that allows injection of custom JSON parsers into the jitjson library.
// This enables applications to use alternative JSON implementations beyond the standard library
// versions of encoding/json/v1 and encoding/json/v2, or other high-performance JSON libraries.
//
// Example usage:
//
//	// Define custom parser
//	type customParser struct{}
//
//	func (p *customParser) Marshal(v interface{}) ([]byte, error) { /* implementation */ }
//	func (p *customParser) Unmarshal(data []byte, v interface{}) error { /* implementation */ }
//
//	// Use custom parser
//	jit := jitjson.NewCustom(value, &customParser{})
//	jsonEncoding, err := jit.Marshal()
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(string(jsonEncoding))
type JSONParser interface {
	// Marshal encodes the given value v into JSON bytes.
	// The behavior should be equivalent to encoding/json.Marshal.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal decodes JSON data into the value pointed to by v.
	// The behavior should be equivalent to encoding/json.Unmarshal.
	Unmarshal(data []byte, v interface{}) error
}
