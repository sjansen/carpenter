package s3util

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/sjansen/carpenter/internal/sys"
)

type Config struct {
	Profile  string
	Region   string
	Bucket   string
	Prefix   string
	Endpoint string
}

func (cfg *Config) newSession(io *sys.IO) (*session.Session, error) {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if cfg.Profile != "" {
		opts.Profile = cfg.Profile
	}

	config := &opts.Config
	config.CredentialsChainVerboseErrors = aws.Bool(true)
	config.LogLevel = aws.LogLevel(aws.LogDebugWithRequestRetries | aws.LogDebugWithRequestErrors)
	config.Logger = aws.LoggerFunc(func(args ...interface{}) {
		io.Log.Debug(args)
	})
	if cfg.Endpoint != "" {
		config.Region = aws.String(cfg.Region)
		config.Endpoint = aws.String(cfg.Endpoint)
		if strings.HasPrefix(cfg.Endpoint, "http://") {
			config.DisableSSL = aws.Bool(true)
		}
		config.S3ForcePathStyle = aws.Bool(true)
	} else if cfg.Region != "" {
		config.Region = aws.String(cfg.Region)
	}

	return session.NewSessionWithOptions(opts)
}

func NewTestConfig() (*Config, error) {
	uri := os.Getenv(s3TestURI)
	if uri == "" {
		return nil, errors.New("$" + s3TestURI + " not set")
	}

	parsed, err := ParseURI(uri)
	switch {
	case err != nil:
		return nil, err
	case parsed.Endpoint != "" && parsed.Region == "":
		return nil, errors.New(
			"invalid $" + s3TestURI + ": region must be set when endpoint is set",
		)
	}

	cfg := parsed.ToConfig()
	return cfg, nil
}
