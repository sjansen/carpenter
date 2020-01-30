package s3util

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestS3Opener(t *testing.T) {
	require := require.New(t)

	if SkipTest() {
		t.Skip("skipping test")
	}

	cfg, err := UploaderTestConfig()
	require.NoError(err)
	require.NotNil(cfg)

	uploader, err := NewUploader(cfg)
	require.NoError(err)
	require.NotNil(uploader)

	uuid := uuid.New().String()
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(uuid),
		Body:   strings.NewReader(uuid),
	})
	require.NoError(err)

	result, err := uploader.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(uuid),
	})
	require.NoError(err)

	b, err := ioutil.ReadAll(result.Body)
	require.NoError(err)
	require.Equal([]byte(uuid), b)
}
