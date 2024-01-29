package protoconf

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/gosynergy/protoconf"
	"github.com/gosynergy/protoconf_test/conf"
)

func TestLoad(t *testing.T) {
	unsetEnvs(t)

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	expectedConf := conf.Config{
		Server: &conf.Config_Server{
			Http: &conf.Config_Server_Http{
				Addr: "127.0.0.1:8080",
				Timeout: &durationpb.Duration{
					Seconds: 1,
				},
			},
			Grpc: &conf.Config_Server_Grpc{
				Addr: "0.0.0.0:9000",
				Timeout: &durationpb.Duration{
					Seconds: 1,
				},
			},
		},
		Data: &conf.Config_Data{
			Database: &conf.Config_Data_Database{
				Driver: "mysql",
				Source: "root:root@tcp(127.0.0.1:3306)/test",
			},
			Redis: &conf.Config_Data_Redis{
				Addr: "127.0.0.1:6379",
				ReadTimeout: &durationpb.Duration{
					Nanos: 200000000,
				},
				WriteTimeout: &durationpb.Duration{
					Nanos: 200000000,
				},
			},
		},
	}

	confDiff := diff(&expectedConf, &cfg)
	if confDiff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", confDiff)
	}
}

func TestLoadWithEnvExpand(t *testing.T) {
	unsetEnvs(t)

	err := os.Setenv("HTTP_ADDR", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("GRPC_ADDR", "localhost:9000")
	if err != nil {
		t.Fatal(err)
	}

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	expectedConf := conf.Config{
		Server: &conf.Config_Server{
			Http: &conf.Config_Server_Http{
				Addr: "localhost:8080",
				Timeout: &durationpb.Duration{
					Seconds: 1,
				},
			},
			Grpc: &conf.Config_Server_Grpc{
				Addr: "localhost:9000",
				Timeout: &durationpb.Duration{
					Seconds: 1,
				},
			},
		},
		Data: &conf.Config_Data{
			Database: &conf.Config_Data_Database{
				Driver: "mysql",
				Source: "root:root@tcp(127.0.0.1:3306)/test",
			},
			Redis: &conf.Config_Data_Redis{
				Addr: "127.0.0.1:6379",
				ReadTimeout: &durationpb.Duration{
					Nanos: 200000000,
				},
				WriteTimeout: &durationpb.Duration{
					Nanos: 200000000,
				},
			},
		},
	}

	confDiff := diff(&expectedConf, &cfg)
	if confDiff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", confDiff)
	}
}

func TestLoadWithValidation(t *testing.T) {
	unsetEnvs(t)

	err := os.Setenv("HTTP_ADDR", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("GRPC_ADDR", "localhost:9000")
	if err != nil {
		t.Fatal(err)
	}

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/invalid-config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	assert.Error(t, err)

	var validationErr *protovalidate.ValidationError
	isValidationError := errors.As(err, &validationErr)
	assert.True(t, isValidationError)

	assert.Equal(t, "server.http.addr", validationErr.Violations[0].FieldPath)
	assert.Equal(t, "required", validationErr.Violations[0].ConstraintId)
	assert.Equal(t, "value is required", validationErr.Violations[0].Message)
}

func TestLoadWithInvalidType(t *testing.T) {
	unsetEnvs(t)

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/invalid-type-config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "invalid google.protobuf.Duration value \"invalid\""))
}

func TestLoadNotProtoMessage(t *testing.T) {
	unsetEnvs(t)

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg interface{}
	err = loader.Load(&cfg)
	assert.Error(t, err)

	if !errors.Is(err, protoconf.ErrTypeAssert) {
		t.Fatalf("expected error %v, got %v", protoconf.ErrTypeAssert, err)
	}
}

func TestLoadWithoutProvider(t *testing.T) {
	unsetEnvs(t)

	loader, err := protoconf.New(
		protoconf.WithParser(yaml.Parser()),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	assert.Error(t, err)

	if !errors.Is(err, protoconf.ErrNoProvider) {
		t.Fatalf("expected error %v, got %v", protoconf.ErrNoProvider, err)
	}
}

func TestLoadWithoutParserCallRead(t *testing.T) {
	unsetEnvs(t)

	provider := MockProvider{}

	loader, err := protoconf.New(
		protoconf.WithProvider(&provider),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	assert.Error(t, err)

	assert.True(t, provider.ReadCalled)
}

func TestLoadWithParserCallUnmarshal(t *testing.T) {
	unsetEnvs(t)

	provider := MockProvider{}
	parser := MockParser{}

	loader, err := protoconf.New(
		protoconf.WithProvider(&provider),
		protoconf.WithParser(&parser),
	)
	if err != nil {
		t.Fatal(err)
	}

	var cfg conf.Config
	err = loader.Load(&cfg)
	assert.Error(t, err)

	assert.True(t, provider.ReadBytesCalled)
	assert.True(t, parser.UnmarshalCalled)
}

func unsetEnvs(t *testing.T) {
	err := os.Unsetenv("HTTP_ADDR")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Unsetenv("GRPC_ADDR")
	if err != nil {
		t.Fatal(err)
	}
}

func diff(want, got interface{}) string {
	return cmp.Diff(want, got,
		cmpopts.IgnoreUnexported(
			conf.Config{},
			conf.Config_Server{},
			conf.Config_Server_Http{},
			conf.Config_Server_Grpc{},
			conf.Config_Data{},
			conf.Config_Data_Database{},
			conf.Config_Data_Redis{},
			durationpb.Duration{},
		),
	)
}
