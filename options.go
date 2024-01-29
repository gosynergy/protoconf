package protoconf

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Option is config option.
type Option func(*options)

// Provider represents a configuration provider. Providers can
// read configuration from a source (file, HTTP etc.)
type Provider interface {
	// ReadBytes Read returns the entire configuration as raw []bytes to be parsed.
	// with a Parser.
	ReadBytes() ([]byte, error)

	// Read returns the parsed configuration as a nested map[string]interface{}.
	// It is important to note that the string keys should not be flat delimited
	// keys like `parent.child.key`, but nested like `{parent: {child: {key: 1}}}`.
	Read() (map[string]interface{}, error)
}

// Parser represents a configuration format parser.
type Parser interface {
	Unmarshal(data []byte) (map[string]interface{}, error)
}

// Validator represents a configuration validator.
type Validator interface {
	Validate(message proto.Message) error
}

// ProtoParser represents a protobuf parser.
type ProtoParser interface {
	Unmarshal(data []byte, message proto.Message) error
}

// JSONParser represents a json parser.
type JSONParser interface {
	Marshal(values map[string]interface{}) ([]byte, error)
}

type options struct {
	provider    Provider
	parser      Parser
	validator   Validator
	protoParser ProtoParser
	jsonParser  JSONParser
}

func WithProvider(p Provider) Option {
	return func(o *options) {
		o.provider = p
	}
}

func WithParser(p Parser) Option {
	return func(o *options) {
		o.parser = p
	}
}

func WithValidator(v Validator) Option {
	return func(o *options) {
		o.validator = v
	}
}

func WithProtoParser(p ProtoParser) Option {
	return func(o *options) {
		o.protoParser = p
	}
}

func WithJSONParser(p JSONParser) Option {
	return func(o *options) {
		o.jsonParser = p
	}
}

type defaultJSONParser struct{}

func (p defaultJSONParser) Marshal(values map[string]interface{}) ([]byte, error) {
	data, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("json marshal config: %w", err)
	}

	return data, nil
}
