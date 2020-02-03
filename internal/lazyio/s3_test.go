package lazyio_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/require"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/s3util"
)

const battlecry = "Spoon!"

func TestS3Opener(t *testing.T) {
	require := require.New(t)

	if s3util.SkipTest() {
		t.Skip("skipping test")
	}

	cfg, err := s3util.NewTestConfig()
	require.NoError(err)
	require.NotNil(cfg)

	uploader, err := s3util.NewUploader(cfg)
	require.NoError(err)
	require.NotNil(uploader)

	o := lazyio.S3Writer{
		Bucket:   cfg.Bucket,
		Prefix:   cfg.Prefix,
		Uploader: uploader,
	}

	w, err := o.Open("battlecry")
	require.NoError(err)

	_, err = w.Write([]byte(battlecry))
	require.NoError(err)

	err = w.Close()
	require.NoError(err)

	var key string
	if strings.HasSuffix(cfg.Prefix, "/") {
		key = cfg.Prefix + "battlecry"
	} else {
		key = cfg.Prefix + "/battlecry"
	}
	result, err := uploader.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(key),
	})
	require.NoError(err)

	b, err := ioutil.ReadAll(result.Body)
	require.NoError(err)
	require.Equal([]byte(battlecry), b)
}
