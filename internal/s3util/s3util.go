package s3util

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const S3_TEST_BUCKET = "CARPENTER_TEST_S3_BUCKET"     //nolint:golint
const S3_TEST_PREFIX = "CARPENTER_TEST_S3_PREFIX"     //nolint:golint
const S3_TEST_REGION = "CARPENTER_TEST_S3_REGION"     //nolint:golint
const S3_TEST_ENDPOINT = "CARPENTER_TEST_S3_ENDPOINT" //nolint:golint

func SkipTest() bool {
	bucket := os.Getenv(S3_TEST_BUCKET)
	return bucket == ""
}

func NewDownloader(cfg *Config) (*s3.S3, error) {
	sess, err := cfg.newSession()
	if err != nil {
		return nil, err
	}
	downloader := s3.New(sess)
	return downloader, nil
}

func NewUploader(cfg *Config) (*s3manager.Uploader, error) {
	sess, err := cfg.newSession()
	if err != nil {
		return nil, err
	}
	uploader := s3manager.NewUploader(sess)
	return uploader, nil
}
