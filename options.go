package protoconf

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

type Transformer interface {
	Transform(values map[string]interface{}) (map[string]interface{}, error)
}

type options struct {
	provider     Provider
	parser       Parser
	transformers []Transformer
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

func WithTransformers(t ...Transformer) Option {
	return func(o *options) {
		o.transformers = append(o.transformers, t...)
	}
}
