package protoconf_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bufbuild/protovalidate-go"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/gosynergy/protoconf"
	mocks "github.com/gosynergy/protoconf/mocks/github.com/gosynergy/protoconf"
	"github.com/gosynergy/protoconf/test/conf"
	"github.com/gosynergy/protoconf/transform/expandenv"
)

type LoadTestSuite struct {
	suite.Suite
}

func (s *LoadTestSuite) SetupTest() {}

func TestLoadTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(LoadTestSuite))
}

func (s *LoadTestSuite) TestLoad() {
	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	s.Require().NoError(err)

	err = loader.Load()
	s.Require().NoError(err)

	var cfg conf.Config
	err = loader.Scan(&cfg)
	s.Require().NoError(err)

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
		s.Failf("config mismatch (-want +got):\n%s", confDiff)
	}
}

//nolint:funlen
func (s *LoadTestSuite) TestLoadWithEnvExpand() {
	const (
		httpAddr = "localhost:8080"
		grpcAddr = "localhost:9000"
	)

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/config-env-expand.yaml")),
		protoconf.WithParser(yaml.Parser()),
		protoconf.WithTransformer(
			expandenv.New(
				expandenv.WithGetenv(func(name string) string {
					switch name {
					case "HTTP_ADDR":
						return httpAddr
					case "GRPC_ADDR":
						return grpcAddr
					}

					return ""
				}),
			),
		),
	)
	s.Require().NoError(err)

	err = loader.Load()
	s.Require().NoError(err)

	var cfg conf.Config
	err = loader.Scan(&cfg)
	s.Require().NoError(err)

	expectedConf := conf.Config{
		Server: &conf.Config_Server{
			Http: &conf.Config_Server_Http{
				Addr: httpAddr,
				Timeout: &durationpb.Duration{
					Seconds: 1,
				},
			},
			Grpc: &conf.Config_Server_Grpc{
				Addr: grpcAddr,
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
		s.Failf("config mismatch (-want +got):\n%s", confDiff)
	}
}

func (s *LoadTestSuite) TestLoadWithValidation() {
	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/invalid-config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	s.Require().NoError(err)

	err = loader.Load()
	s.Require().NoError(err)

	var cfg conf.Config
	err = loader.Scan(&cfg)
	s.Require().Error(err)

	var validationErr *protovalidate.ValidationError
	isValidationError := errors.As(err, &validationErr)
	s.True(isValidationError)

	s.Equal("server.http.addr", validationErr.Violations[0].GetFieldPath())
	s.Equal("required", validationErr.Violations[0].GetConstraintId())
	s.Equal("value is required", validationErr.Violations[0].GetMessage())
}

func (s *LoadTestSuite) TestLoadWithInvalidType() {
	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/invalid-type-config.yaml")),
		protoconf.WithParser(yaml.Parser()),
	)
	s.Require().NoError(err)

	err = loader.Load()
	s.Require().NoError(err)

	var cfg conf.Config
	err = loader.Scan(&cfg)
	s.Require().Error(err)
	s.True(strings.Contains(err.Error(), "invalid google.protobuf.Duration"))
}

func (s *LoadTestSuite) TestLoadWithoutProvider() {
	_, err := protoconf.New(
		protoconf.WithParser(yaml.Parser()),
	)
	s.Require().Error(err)
	s.Require().ErrorIs(err, protoconf.ErrNoProvider)
}

func (s *LoadTestSuite) TestLoadWithCustomTransformer() {
	transformer := mocks.NewMockTransformer(s.T())
	transformer.
		EXPECT().
		Transform(mock.Anything).
		Once().
		Return(map[string]interface{}{}, nil)

	loader, err := protoconf.New(
		protoconf.WithProvider(file.Provider("conf/config.yaml")),
		protoconf.WithParser(yaml.Parser()),
		protoconf.WithTransformer(transformer),
	)
	s.Require().NoError(err)

	err = loader.Load()
	s.Require().NoError(err)
}

func (s *LoadTestSuite) TestLoadWithProviderWithoutParser() {
	provider := mocks.NewMockProvider(s.T())
	provider.EXPECT().
		Read().
		Return(map[string]interface{}{}, nil)

	loader, err := protoconf.New(
		protoconf.WithProvider(provider),
	)
	s.Require().NoError(err)

	err = loader.Load()
	s.Require().NoError(err)
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
