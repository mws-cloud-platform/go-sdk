package mws

import (
	"cmp"
	"fmt"
	"time"

	"go.mws.cloud/util-toolset/pkg/os/env"
)

const (
	// DefaultBaseEndpoint is a default MWS Cloud Platform API base endpoint.
	DefaultBaseEndpoint = "https://api.mwsapis.ru"
	// DefaultZone is a default MWS Cloud Platform availability zone.
	DefaultZone = "ru-central1-a"
	// DefaultTimeout is a default HTTP client request timeout.
	DefaultTimeout = 5 * time.Second

	BaseEndpointEnv                    = "MWS_BASE_ENDPOINT"
	ProjectEnv                         = "MWS_PROJECT"
	ZoneEnv                            = "MWS_ZONE"
	TokenEnv                           = "MWS_TOKEN"
	ServiceAccountAuthorizedKeyPathEnv = "MWS_SERVICE_ACCOUNT_AUTHORIZED_KEY_PATH"
	TimeoutEnv                         = "MWS_TIMEOUT"
	LogLevelEnv                        = "MWS_LOG_LEVEL"
)

// Config is an SDK configuration.
//
// Use [LoadConfig] to load configuration from environment variables and
// sensible defaults.
type Config struct {
	// MWS Cloud Platform API base endpoint. Can be specified using the
	// `MWS_BASE_ENDPOINT` environment variable.
	BaseEndpoint string `yaml:"base_endpoint"`
	// Default project. Can be specified using the `MWS_PROJECT` environment
	// variable.
	Project string `yaml:"project,omitempty"`
	// Default zone. Can be specified using the `MWS_ZONE` environment variable.
	Zone string `yaml:"zone,omitempty"`
	// IAM token for authentication. Can be specified using the `MWS_TOKEN`
	// environment variable.
	Token string `yaml:"-"`
	// Path to the service account authorized key file used for authentication.
	// Has no effect if Token is not empty. Can be specified using the
	// `MWS_SERVICE_ACCOUNT_AUTHORIZED_KEY_PATH` environment variable.
	ServiceAccountAuthorizedKeyPath string `yaml:"service_account_authorized_key_path,omitempty"`
	// HTTP client request timeout. Can be specified using the `MWS_TIMEOUT`
	// environment variable.
	Timeout time.Duration `yaml:"timeout,omitempty"`
	// Log level for the SDK. Can be specified using the `MWS_LOG_LEVEL`
	// environment variable.
	LogLevel string `yaml:"log_level,omitempty"`
}

// LoadConfig loads SDK configuration from environment variables and sensible
// defaults.
func LoadConfig(opts ...LoadConfigOption) (config *Config, err error) {
	o := &loadConfigOptions{
		env: env.RealEnv{},
	}
	for _, opt := range opts {
		opt(o)
	}

	timeout := DefaultTimeout
	if v, ok := o.env.LookupEnv(TimeoutEnv); ok {
		timeout, err = time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("parse %q: %w", TimeoutEnv, err)
		}
	}

	return &Config{
		BaseEndpoint:                    cmp.Or(o.env.Getenv(BaseEndpointEnv), DefaultBaseEndpoint),
		Project:                         o.env.Getenv(ProjectEnv),
		Zone:                            cmp.Or(o.env.Getenv(ZoneEnv), DefaultZone),
		Token:                           o.env.Getenv(TokenEnv),
		ServiceAccountAuthorizedKeyPath: o.env.Getenv(ServiceAccountAuthorizedKeyPathEnv),
		Timeout:                         timeout,
		LogLevel:                        o.env.Getenv(LogLevelEnv),
	}, nil
}

// LoadConfigOption is a functional option for SDK configuration loading.
type LoadConfigOption func(*loadConfigOptions)

// LoadConfigWithEnv sets the environment for SDK configuration loading.
func LoadConfigWithEnv(env env.Env) LoadConfigOption {
	return func(o *loadConfigOptions) {
		o.env = env
	}
}

type loadConfigOptions struct {
	env env.Env
}
