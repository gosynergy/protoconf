package protoconf

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"mvdan.cc/sh/v3/shell"
)

var (
	ErrTypeAssert = errors.New("type assert error")
	ErrNoProvider = errors.New("no provider")
)

type ConfigLoader struct {
	opts      options
	validator *protovalidate.Validator
}

var _ Loader = (*ConfigLoader)(nil)

type Loader interface {
	Load(v interface{}) error
}

func New(opts ...Option) (*ConfigLoader, error) {
	confOpts := options{}
	for _, opt := range opts {
		opt(&confOpts)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("protovalidate new: %w", err)
	}

	return &ConfigLoader{
		opts:      confOpts,
		validator: validator,
	}, nil
}

func (c *ConfigLoader) Load(v interface{}) error {
	message, ok := v.(proto.Message)
	if !ok {
		return ErrTypeAssert
	}

	values, err := c.parse()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("json marshal config: %w", err)
	}

	expandedData, err := shell.Expand(string(data), os.Getenv)
	if err != nil {
		return fmt.Errorf("shell expand config: %w", err)
	}

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal([]byte(expandedData), message)
	if err != nil {
		return fmt.Errorf("protojson unmarshal config: %w", err)
	}

	err = c.validator.Validate(message)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func (c *ConfigLoader) parse() (map[string]interface{}, error) {
	if c.opts.provider == nil {
		return nil, ErrNoProvider
	}

	if c.opts.parser == nil {
		values, err := c.opts.provider.Read()
		if err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}

		return values, nil
	}

	data, err := c.opts.provider.ReadBytes()
	if err != nil {
		return nil, fmt.Errorf("read config bytes: %w", err)
	}

	values, err := c.opts.parser.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return values, nil
}
