package s3util

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type UploaderConfig struct {
	Profile  string
	Region   string
	Bucket   string
	Prefix   string
	Endpoint string
}

func NewUploader(cfg *UploaderConfig) (*s3manager.Uploader, error) {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if cfg.Profile != "" {
		opts.Profile = cfg.Profile
	}

	config := &opts.Config
	config.CredentialsChainVerboseErrors = aws.Bool(true)
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

	sess := session.Must(session.NewSessionWithOptions(opts))
	uploader := s3manager.NewUploader(sess)
	return uploader, nil
}

func UploaderTestConfig() (*UploaderConfig, error) {
	bucket := os.Getenv(S3_TEST_BUCKET)
	if bucket == "" {
		return nil, errors.New("$" + S3_TEST_BUCKET + " not set")
	}

	prefix := os.Getenv(S3_TEST_PREFIX)
	region := os.Getenv(S3_TEST_REGION)
	endpoint := os.Getenv(S3_TEST_ENDPOINT)
	if endpoint != "" && region == "" {
		return nil, errors.New(
			"$" + S3_TEST_REGION + " must be set when $" + S3_TEST_ENDPOINT + " is set",
		)
	}

	return &UploaderConfig{
		Region:   region,
		Bucket:   bucket,
		Prefix:   prefix,
		Endpoint: endpoint,
	}, nil
}
