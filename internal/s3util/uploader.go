package s3util

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type UploaderFactory struct {
	Profile  string
	Region   string
	Bucket   string
	Prefix   string
	Endpoint string
}

func (f *UploaderFactory) New() (*s3manager.Uploader, error) {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if f.Profile != "" {
		opts.Profile = f.Profile
	}

	config := &opts.Config
	config.CredentialsChainVerboseErrors = aws.Bool(true)
	if f.Endpoint != "" {
		config.Region = aws.String(f.Region)
		config.Endpoint = aws.String(f.Endpoint)
		if strings.HasPrefix(f.Endpoint, "http://") {
			config.DisableSSL = aws.Bool(true)
		}
		config.S3ForcePathStyle = aws.Bool(true)
	} else if f.Region != "" {
		config.Region = aws.String(f.Region)
	}

	sess := session.Must(session.NewSessionWithOptions(opts))
	uploader := s3manager.NewUploader(sess)
	return uploader, nil
}

func NewTestUploaderFactory() (*UploaderFactory, error) {
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

	return &UploaderFactory{
		Region:   region,
		Bucket:   bucket,
		Prefix:   prefix,
		Endpoint: endpoint,
	}, nil
}
