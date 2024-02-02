# expandenv

The `expandenv` transformer expands environment variables in the configuration data.

This transformer uses [mvdan.cc/sh/v3/shell](https://pkg.go.dev/mvdan.cc/sh/v3/shell) internally.

## Usage

[//]: @formatter:off

```go
import (
    "github.com/gosynergy/protoconf/transform/expandenv"
)

loader, err := protoconf.New(
  ...
  protoconf.WithTransformers(
    expandenv.NewTransformer(),
  ),
)

```

[//]: @formatter:on
