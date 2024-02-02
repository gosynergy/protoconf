package expandenv

import (
	"encoding/json"
	"fmt"
	"os"

	"mvdan.cc/sh/v3/shell"
)

type Transformer struct {
	opts options
}

func (t *Transformer) Transform(values map[string]interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(values)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	content, err := shell.Expand(string(data), t.opts.getenv)
	if err != nil {
		return nil, fmt.Errorf("expand: %w", err)
	}

	var expanded map[string]interface{}

	err = json.Unmarshal([]byte(content), &expanded)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return expanded, nil
}

func NewTransformer(opts ...Option) *Transformer {
	confOpts := options{}

	for _, opt := range opts {
		opt(&confOpts)
	}

	if confOpts.getenv == nil {
		confOpts.getenv = os.Getenv
	}

	return &Transformer{
		opts: confOpts,
	}
}
