package jitjson

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"unsafe"
)

// JSONParser is an interface for allowing use of custom JSON parsers.
// This is useful for applications that may use multiple JSON parsers,
// such as encoding/json and encoding/json/v2.
// The default parser is "encoding/json".
type JSONParser interface {
	Name() string
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

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

// defaultParser is the default JSONParser used by the library.
var defaultParser unsafe.Pointer
var defaultParserName atomic.Value

var parsers = make(map[string]JSONParser)

func init() {
	var stdParser JSONParser = &jsonParser{}
	parsers[stdParser.Name()] = stdParser
	atomic.StorePointer(&defaultParser, unsafe.Pointer(&stdParser))
	defaultParserName.Store(stdParser.Name())
}

// RegisterParser adds a new JSONParser to the registry.
// If the parser is nil, provides no name, or is already registered, an error will be returned.
// The default pre-registered parser is "encoding/json".
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

// MustRegisterParser panics if RegisterParser fails.
func MustRegisterParser(parser JSONParser) {
	if parser == nil {
		panic("parser is nil")
	}
	if err := RegisterParser(parser); err != nil {
		panic(err)
	}
}

// SetDefaultParser changes the global default parser.
// Returns an error if the parser is not pre-registered.
// The default pre-registered parser is "encoding/json" which can be changed, or reset.
func SetDefaultParser(name string) error {
	parser, exists := parsers[name]
	if !exists {
		return fmt.Errorf("parser %s not registered", name)
	}
	defaultParserName.Store(name)
	atomic.StorePointer(&defaultParser, unsafe.Pointer(&parser))
	return nil
}

// MustSetDefaultParser panics if SetDefaultParser fails.
func MustSetDefaultParser(name string) {
	if err := SetDefaultParser(name); err != nil {
		panic(err)
	}
}

// DefaultParser returns the name of the current default parser.
func DefaultParser() string {
	return defaultParserName.Load().(string)
}

func getDefaultParser() JSONParser {
	return *(*JSONParser)(atomic.LoadPointer(&defaultParser))
}
