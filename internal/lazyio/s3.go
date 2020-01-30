package lazyio

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var _ Opener = &S3Opener{}

type S3Opener struct {
	Bucket   string
	Prefix   string
	Uploader *s3manager.Uploader
}

func (o *S3Opener) Open(path string) (io.WriteCloser, error) {
	ch := make(chan error)
	r, w := io.Pipe()
	obj := &s3object{
		ch: ch,
		w:  w,
	}

	go o.upload(path, r, ch)
	return obj, nil
}

func (o *S3Opener) upload(path string, r io.Reader, ch chan<- error) {
	_, err := o.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(o.Bucket),
		Key:    aws.String(o.Prefix + path),
		Body:   r,
	})
	ch <- err
}

/*
* s3object
 */
type s3object struct {
	ch <-chan error
	w  io.WriteCloser
}

func (o *s3object) Close() error {
	o.w.Close()
	return <-o.ch
}

func (o *s3object) Write(p []byte) (int, error) {
	return o.w.Write(p)
}
