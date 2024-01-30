# protoconf

`protoconf` is a Go-based project that provides an opinionated way to define configuration using Protocol Buffers (protobuf). It leverages the power of protobuf to create structured and strongly-typed configuration files, ensuring type safety and reducing the risk of configuration errors.

## Features

- **Strongly Typed Configuration**: By using protobuf, `protoconf` ensures that your configuration is strongly typed, reducing the risk of configuration errors.
- **Flexible Configuration Providers**: `protoconf` supports different configuration providers, allowing you to read configuration from various sources such as files, HTTP, etc.
- **Customizable Parsers**: You can use different parsers to parse your configuration data, providing flexibility in how you structure your configuration files.
- **Transformers**: `protoconf` allows you to transform your configuration data as needed, providing an additional layer of flexibility.
- **Validation**: `protoconf` includes built-in validation using `protovalidate`, ensuring that your configuration adheres to the defined protobuf structure.

## Usage

To use `protoconf`, you need to define your configuration structure using protobuf. Then, you can create a `ConfigLoader` with the desired options, such as the configuration provider and parser. You can also add transformers if needed.

Here is a basic example:

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
    expandenv.New(),
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

In this example, `protoconf` reads the configuration from a YAML file, parses it, and scans it into a `conf.Config` protobuf message.

## Contributing

Contributions to `protoconf` are welcome! Please submit a pull request or create an issue if you have any improvements or features you'd like to add.

## License

`protoconf` is licensed under the MIT License. See the `LICENSE` file for more information.
