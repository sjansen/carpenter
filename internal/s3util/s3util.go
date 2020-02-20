package s3util

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sjansen/carpenter/internal/sys"
)

const s3TestURI = "CARPENTER_TEST_S3_URI"

func SkipTest() bool {
	bucket := os.Getenv(s3TestURI)
	return bucket == ""
}

func NewDownloader(io *sys.IO, cfg *Config) (*s3.S3, error) {
	sess, err := cfg.newSession(io)
	if err != nil {
		return nil, err
	}
	downloader := s3.New(sess)
	return downloader, nil
}

func NewUploader(io *sys.IO, cfg *Config) (*s3manager.Uploader, error) {
	sess, err := cfg.newSession(io)
	if err != nil {
		return nil, err
	}
	uploader := s3manager.NewUploader(sess)
	return uploader, nil
}
