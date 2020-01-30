package s3util

import "os"

const S3_TEST_BUCKET = "CARPENTER_TEST_S3_BUCKET"     //nolint:golint
const S3_TEST_PREFIX = "CARPENTER_TEST_S3_PREFIX"     //nolint:golint
const S3_TEST_REGION = "CARPENTER_TEST_S3_REGION"     //nolint:golint
const S3_TEST_ENDPOINT = "CARPENTER_TEST_S3_ENDPOINT" //nolint:golint

func SkipTest() bool {
	bucket := os.Getenv(S3_TEST_BUCKET)
	return bucket == ""
}
