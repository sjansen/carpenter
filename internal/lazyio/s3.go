package lazyio

import (
	"io"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sjansen/carpenter/internal/s3util"
)

var _ OutputOpener = &S3Writer{}

type S3Writer struct {
	Bucket   string
	Prefix   string
	Uploader *s3manager.Uploader
}

func NewS3Opener(uri string) (*S3Writer, error) {
	parsed, err := s3util.ParseURI(uri)
	if err != nil {
		return nil, err
	}
	cfg := parsed.ToUploaderConfig()
	uploader, err := s3util.NewUploader(cfg)
	if err != nil {
		return nil, err
	}
	opener := &S3Writer{
		Bucket:   cfg.Bucket,
		Prefix:   cfg.Prefix,
		Uploader: uploader,
	}
	return opener, nil
}

func (o *S3Writer) Open(path string) (io.WriteCloser, error) {
	ch := make(chan error)
	r, w := io.Pipe()
	obj := &s3object{
		ch: ch,
		w:  w,
	}

	go o.upload(path, r, ch)
	return obj, nil
}

func (o *S3Writer) upload(suffix string, r io.Reader, ch chan<- error) {
	key := path.Join(o.Prefix, suffix)
	_, err := o.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(o.Bucket),
		Key:    aws.String(key),
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
	err1 := o.w.Close()
	err2 := <-o.ch
	if err2 != nil {
		return err2
	}
	return err1
}

func (o *s3object) Write(p []byte) (int, error) {
	return o.w.Write(p)
}
