package jitjson

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// JSONParser is an interface that allows injection of custom JSON parsers into the jitjson library.
// This enables applications to use alternative JSON implementations beyond the standard library,
// such as encoding/json/v2, or other high-performance JSON libraries.
//
// The default parser is "encoding/json" from the standard library. Custom parsers can be
// registered using RegisterParser() and set as the default using SetDefaultParser().
//
// Example usage:
//
//	type customParser struct{}
//
//	func (p *customParser) Name() string { return "custom-json" }
//	func (p *customParser) Marshal(v interface{}) ([]byte, error) { /* implementation */ }
//	func (p *customParser) Unmarshal(data []byte, v interface{}) error { /* implementation */ }
//
//	jitjson.MustRegisterParser(&customParser{})
//	jitjson.MustSetDefaultParser("custom-json")
type JSONParser interface {
	// Name returns a unique identifier for this parser implementation.
	// This name is used for registration and selection of parsers.
	Name() string

	// Marshal encodes the given value v into JSON bytes.
	// The behavior should be equivalent to encoding/json.Marshal.
	Marshal(v interface{}) ([]byte, error)

	// Unmarshal decodes JSON data into the value pointed to by v.
	// The behavior should be equivalent to encoding/json.Unmarshal.
	Unmarshal(data []byte, v interface{}) error
}

// jsonParser is the default JSONParser using the encoding/json package.
type jsonParser struct{}

var _ JSONParser = (*jsonParser)(nil)

func (j *jsonParser) Name() string {
	return "encoding/json"
}

func (j *jsonParser) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonParser) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// parsers registry with the current default parser pointers
var (
	parsers           map[string]JSONParser
	defaultParser     unsafe.Pointer
	defaultParserName atomic.Value
	registryOnce      sync.Once
)

func init() {
	initParserRegistry()
}

// initParserRegistry initializes the parser registry using sync.Once to prevent race conditions.
// This ensures the registry is only initialized once, even when called from multiple init() functions.
func initParserRegistry() {
	registryOnce.Do(func() {
		parsers = make(map[string]JSONParser)
		setupParserRegistry()
	})
}

// setupParserRegistry sets up the parser registry with the default parser.
// Both json/v1 and json/v2 are supported.
func setupParserRegistry() {
	var stdParser JSONParser = &jsonParser{}

	parsers[stdParser.Name()] = stdParser
	parsers["encoding/json/v1"] = stdParser

	atomic.StorePointer(&defaultParser, unsafe.Pointer(&stdParser))
	defaultParserName.Store(stdParser.Name())
}

// RegisterParser adds a JSONParser implementation to the global parser registry.
// This allows the jitjson library to use alternative JSON parsing implementations
// for just-in-time marshalling and unmarshalling operations.
//
// The parser must provide a unique name via its Name() method, which will be used
// as the identifier for selecting this parser. If a parser with the same name
// already exists, an error is returned.
//
// Example:
//
//	parser := &myCustomParser{}
//	err := jitjson.RegisterParser(parser)
//	if err != nil {
//	    log.Fatal(err)
//	}
func RegisterParser(parser JSONParser) error {
	if parser == nil {
		return fmt.Errorf("parser is nil")
	}
	name := parser.Name()
	if name == "" {
		return fmt.Errorf("parser name is empty")
	}
	if _, exists := parsers[name]; exists {
		return fmt.Errorf("parser %s already registered", name)
	}
	parsers[name] = parser
	return nil
}

// MustRegisterParser registers a JSONParser implementation and panics on failure.
// This is a convenience wrapper around [RegisterParser].
//
// Example:
//
//	parser := &myCustomParser{}
//	jitjson.MustRegisterParser(parser) // panics if registration fails
func MustRegisterParser(parser JSONParser) {
	if parser == nil {
		panic("parser is nil")
	}
	if err := RegisterParser(parser); err != nil {
		panic(err)
	}
}

// DefaultParser returns the name of the currently configured default parser.
// The default parser is initially set to "encoding/json" but can be changed
// using [SetDefaultParser].
func DefaultParser() string {
	return defaultParserName.Load().(string)
}

// SetDefaultParser changes the global default parser used by all new JitJSON instances.
// The parser must already be registered via [RegisterParser].
//
// By default, the library uses "encoding/json" as the default parser. You can change
// this to any registered parser name, including custom implementations or alternative
// JSON libraries.
//
// Example:
//
//	// Assuming a custom parser was registered with name "fast-json"
//	err := jitjson.SetDefaultParser("fast-json")
//	if err != nil {
//	    log.Fatal(err)
//	}
func SetDefaultParser(name string) error {
	initParserRegistry() // Ensure registry is initialized
	parser, exists := parsers[name]
	if !exists {
		return fmt.Errorf("parser %s not registered", name)
	}
	defaultParserName.Store(name)
	atomic.StorePointer(&defaultParser, unsafe.Pointer(&parser))
	return nil
}

// MustSetDefaultParser changes the global default parser and panics on failure.
// This is a convenience wrapper around [SetDefaultParser].
//
// Example:
//
//	parser := &jsonV2Parser{}
//	jitjson.MustRegisterParser(parser) // panics if registration fails
//	jitjson.MustSetDefaultParser("encoding/json/v2") // panics if parser not registered
func MustSetDefaultParser(name string) {
	if err := SetDefaultParser(name); err != nil {
		panic(err)
	}
}

func getDefaultParser() JSONParser {
	return *(*JSONParser)(atomic.LoadPointer(&defaultParser))
}
