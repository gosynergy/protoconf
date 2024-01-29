package protoconf

import (
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
	opts options
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

	if confOpts.validator == nil {
		validator, err := protovalidate.New()
		if err != nil {
			return nil, fmt.Errorf("protovalidate new: %w", err)
		}

		confOpts.validator = validator
	}

	if confOpts.protoParser == nil {
		confOpts.protoParser = &protojson.UnmarshalOptions{DiscardUnknown: true}
	}

	if confOpts.jsonParser == nil {
		confOpts.jsonParser = defaultJSONParser{}
	}

	return &ConfigLoader{
		opts: confOpts,
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

	data, err := c.opts.jsonParser.Marshal(values)
	if err != nil {
		return fmt.Errorf("json marshal config: %w", err)
	}

	err = c.opts.protoParser.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("protojson unmarshal config: %w", err)
	}

	err = c.opts.validator.Validate(message)
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

	content, err := shell.Expand(string(data), os.Getenv)
	if err != nil {
		return nil, fmt.Errorf("shell expand: %w", err)
	}

	values, err := c.opts.parser.Unmarshal([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return values, nil
}
