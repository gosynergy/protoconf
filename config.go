package protoconf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var ErrNoProvider = errors.New("no provider")

type ConfigLoader struct {
	opts      options
	validator *protovalidate.Validator
	values    map[string]interface{}
}

var _ Loader = (*ConfigLoader)(nil)

type Loader interface {
	Load() error
	Scan(message proto.Message) error
}

func New(opts ...Option) (*ConfigLoader, error) {
	confOpts := options{}
	for _, opt := range opts {
		opt(&confOpts)
	}

	if confOpts.provider == nil {
		return nil, ErrNoProvider
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

func (c *ConfigLoader) Scan(message proto.Message) error {
	var err error

	err = c.unmarshal(message)
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	err = c.validator.Validate(message)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func (c *ConfigLoader) Load() error {
	var err error

	err = c.parse()
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	err = c.transform()
	if err != nil {
		return fmt.Errorf("transform config: %w", err)
	}

	return nil
}

func (c *ConfigLoader) parse() error {
	var err error

	if c.opts.parser == nil {
		c.values, err = c.opts.provider.Read()
		if err != nil {
			return fmt.Errorf("read config: %w", err)
		}

		return nil
	}

	data, err := c.opts.provider.ReadBytes()
	if err != nil {
		return fmt.Errorf("read config bytes: %w", err)
	}

	c.values, err = c.opts.parser.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}

func (c *ConfigLoader) transform() error {
	var err error

	for _, t := range c.opts.transformers {
		c.values, err = t.Transform(c.values)
		if err != nil {
			return fmt.Errorf("transform config: %w", err)
		}
	}

	return nil
}

func (c *ConfigLoader) unmarshal(message proto.Message) error {
	var err error

	data, err := json.Marshal(c.values)
	if err != nil {
		return fmt.Errorf("json marshal config: %w", err)
	}

	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("protojson unmarshal config: %w", err)
	}

	return nil
}
