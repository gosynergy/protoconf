# protoconf

[![Go Reference](https://pkg.go.dev/badge/github.com/gosynergy/protoconf.svg)](https://pkg.go.dev/github.com/gosynergy/protoconf)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Test](https://github.com/gosynergy/protoconf/actions/workflows/go.yml/badge.svg)](https://github.com/gosynergy/protoconf/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/gh/gosynergy/protoconf/graph/badge.svg?token=D2XW7B9I6J)](https://codecov.io/gh/gosynergy/protoconf)
[![Maintainability](https://api.codeclimate.com/v1/badges/8b29f40a200ca4f0af04/maintainability)](https://codeclimate.com/github/gosynergy/protoconf/maintainability)

`protoconf` is a Go-based project that provides an opinionated way to define configuration using Protocol Buffers. It
leverages the power of protobuf to create structured and strongly-typed configuration files, ensuring type
safety and reducing the risk of configuration errors.

## Features

- **Strongly Typed Configuration**: By using protobuf, `protoconf` ensures that your configuration is strongly typed,
  reducing the risk of configuration errors.
- **Flexible Configuration Providers**: `protoconf` supports different configuration providers, allowing you to read
  configuration from various sources such as files, HTTP, etc.
- **Customizable Parsers**: You can use different parsers to parse your configuration data, providing flexibility in how
  you structure your configuration files.
- **Transformers**: `protoconf` allows you to transform your configuration data as needed, providing an additional layer
  of flexibility.
- **Validation**: `protoconf` includes built-in validation using `protovalidate`, ensuring that your configuration
  adheres to the defined protobuf structure.

## Usage

To use `protoconf`, you need to define your configuration structure using protobuf. Then, you can create
a `ConfigLoader` with the desired options, such as the configuration provider and parser. You can also add transformers
if needed.

Here is a basic example:

[//]: @formatter:off

```go
import (
    "github.com/knadh/koanf/parsers/yaml"
    "github.com/knadh/koanf/providers/file"
    "github.com/gosynergy/protoconf/transform/expandenv"

    "github.com/gosynergy/protoconf/conf"
)

loader, err := protoconf.New(
  protoconf.WithProvider(file.Provider("conf/config.yaml")),
  protoconf.WithParser(yaml.Parser()),
  protoconf.WithTransformers(
    expandenv.NewTransformer(),
  ),
)
if err != nil {
  // handle error
}

err = loader.Load()
if err != nil {
  // handle error
}

var cfg conf.Config
err = loader.Scan(&cfg)
if err != nil {
  // handle error
}
```

[//]: @formatter:on

See the [config_test.go](config_test.go) for more examples.

### Provider

The provider is responsible for reading the configuration data from a source. `protoconf` supports different providers,
such as file, HTTP, etc. In the example above, `file.Provider` is used to read the configuration from a YAML file.

[//]: @formatter:off

```go
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
```

[//]: @formatter:on

`protoconf` Provider compatible with [koanf](https://github.com/knadh/koanf?tab=readme-ov-file#api) providers.

### Parser

The parser is responsible for parsing the configuration data into a format that can be scanned into a protobuf
message. `protoconf` supports different parsers, such as JSON, YAML, etc. In the example above, `yaml.Parser` is used to
parse the configuration data from a YAML file.

In this example, `protoconf` reads the configuration from a YAML file, parses it, and scans it into a `conf.Config`
protobuf message.

[//]: @formatter:off

```go
// Parser represents a configuration format parser.
type Parser interface {
	Unmarshal(data []byte) (map[string]interface{}, error)
}
```

[//]: @formatter:on

`protoconf` Parser compatible with [koanf](https://github.com/knadh/koanf?tab=readme-ov-file#api) parsers.

### Transformer

Transformers are used to transform the configuration data as needed. `protoconf` supports different transformers, such
as `expandenv`, `mapkeys`, etc. In the example above, `expandenv.NewTransformer()` is used to expand environment
variables in the configuration data.

[//]: @formatter:off

```go
type Transformer interface {
	Transform(values map[string]interface{}) (map[string]interface{}, error)
}
```

[//]: @formatter:on

Built-in [expandenv](transform/expandenv) is a transformer that expands environment variables in the configuration data.

## Contributing

Contributions to `protoconf` are welcome! Please submit a pull request or create an issue if you have any improvements
or features you'd like to add.

## License

`protoconf` is licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for more information.
