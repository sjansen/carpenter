package lazyio

import (
	"io"
	pathlib "path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/sjansen/carpenter/internal/s3util"
)

var _ InputOpener = &S3Reader{}
var _ InputWalker = &S3Reader{}
var _ OutputOpener = &S3Writer{}

type S3Reader struct {
	Bucket     string
	Prefix     string
	Downloader *s3.S3
}

func NewS3Reader(uri string) (*S3Reader, error) {
	parsed, err := s3util.ParseURI(uri)
	if err != nil {
		return nil, err
	}
	cfg := parsed.ToConfig()
	downloader, err := s3util.NewDownloader(cfg)
	if err != nil {
		return nil, err
	}
	reader := &S3Reader{
		Bucket:     cfg.Bucket,
		Prefix:     cfg.Prefix,
		Downloader: downloader,
	}
	return reader, nil
}

func (r *S3Reader) Open(path string) (io.ReadCloser, error) {
	key := pathlib.Join(r.Prefix, path)
	result, err := r.Downloader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (r *S3Reader) Walk(fn func(string) error) error {
	prefix := r.Prefix
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	err := r.Downloader.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket:    aws.String(r.Bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			if aws.Int64Value(obj.Size) > 0 {
				key := aws.StringValue(obj.Key)
				suffix := r.stripPrefix(key)
				fn(suffix)
			}
		}
		return lastPage
	})
	return err
}

func (r *S3Reader) stripPrefix(key string) string {
	n := len(r.Prefix)
	if n > 0 && r.Prefix[n-1] == '/' {
		key = key[n:]
	} else {
		key = key[n+1:]
	}
	return key
}

type S3Writer struct {
	Bucket   string
	Prefix   string
	Uploader *s3manager.Uploader
}

func NewS3Writer(uri string) (*S3Writer, error) {
	parsed, err := s3util.ParseURI(uri)
	if err != nil {
		return nil, err
	}
	cfg := parsed.ToConfig()
	uploader, err := s3util.NewUploader(cfg)
	if err != nil {
		return nil, err
	}
	writer := &S3Writer{
		Bucket:   cfg.Bucket,
		Prefix:   cfg.Prefix,
		Uploader: uploader,
	}
	return writer, nil
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
	key := pathlib.Join(o.Prefix, suffix)
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
