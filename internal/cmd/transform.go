package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/sjansen/carpenter/internal/lazyio"
	"github.com/sjansen/carpenter/internal/patterns"
	"github.com/sjansen/carpenter/internal/pipeline"
	"github.com/sjansen/carpenter/internal/s3util"
	"github.com/sjansen/carpenter/internal/sys"
	"github.com/sjansen/carpenter/internal/tokenizer"
	"github.com/sjansen/carpenter/internal/uaparser"
)

type TransformCmd struct {
	Patterns string
	SrcURI   string
	DstURI   string
	ErrURI   string
}

func (c *TransformCmd) Run(base *Base) error {
	io := &base.IO
	log := io.Log
	log.Debugw("creating pipeline")
	pipeline, walker, err := c.newPipeline(io)
	if err != nil {
		return err
	}

	log.Debugw("starting pipeline")
	pipeline.Start()
	defer log.Debugw("pipeline finished")
	defer pipeline.Wait()
	defer log.Debugw("waiting for pipeline to finish")

	return walker.Walk(func(path string) error {
		log.Debugw("adding task to pipeline", "path", path)
		pipeline.AddTask(path)
		return nil
	})
}

func (c *TransformCmd) newPipeline(io *sys.IO) (*pipeline.Pipeline, lazyio.InputWalker, error) {
	log := io.Log
	patterns, err := loadPatterns(io, c.Patterns)
	if err != nil {
		return nil, nil, err
	}

	log.Debugw("loading user-agent parser")
	uaparser, err := uaparser.UserAgentParser()
	if err != nil {
		return nil, nil, err
	}

	pipeline := &pipeline.Pipeline{
		IO:        io,
		Patterns:  patterns,
		Tokenizer: tokenizer.ALB,
		UAParser:  uaparser,
	}

	input, err := newInputOpenWalker(io, c.SrcURI)
	if err != nil {
		return nil, nil, err
	}
	pipeline.Source = input

	output, err := newOutputOpener(io, c.DstURI)
	if err != nil {
		return nil, nil, err
	}
	pipeline.Result = output

	if c.ErrURI != "" {
		opener, err := newOutputOpener(io, c.ErrURI)
		if err != nil {
			return nil, nil, err
		}
		pipeline.Debug = opener
	}

	return pipeline, input, nil
}

func loadPatterns(io *sys.IO, uri string) (*patterns.Patterns, error) {
	log := io.Log
	switch {
	case strings.HasPrefix(uri, "s3://") || strings.HasPrefix(uri, "S3://"):
		log.Debugw("loading patterns from S3", "uri", uri)
		parsed, err := s3util.ParseURI(uri)
		if err != nil {
			return nil, err
		}
		cfg := parsed.ToConfig()
		downloader, err := s3util.NewDownloader(io, cfg)
		if err != nil {
			return nil, err
		}
		result, err := downloader.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(parsed.Bucket),
			Key:    aws.String(parsed.Key),
		})
		if err != nil {
			return nil, err
		}
		return patterns.Load(uri, result.Body)
	default:
		log.Debugw("loading patterns from FS", "uri", uri)
		r, err := os.Open(uri)
		if err != nil {
			return nil, err
		}
		return patterns.Load(uri, r)
	}
}

type inputOpenWalker interface {
	lazyio.InputOpener
	lazyio.InputWalker
}

func newInputOpenWalker(io *sys.IO, uri string) (inputOpenWalker, error) {
	log := io.Log
	switch {
	case strings.HasPrefix(uri, "s3://") || strings.HasPrefix(uri, "S3://"):
		log.Debugw("creating S3 input walker", "uri", uri)
		return lazyio.NewS3Reader(io, uri)
	default:
		log.Debugw("creating FS input walker", "uri", uri)
		uri = filepath.Clean(uri)
		if info, err := os.Stat(uri); err != nil {
			return nil, err
		} else if !info.IsDir() {
			return nil, fmt.Errorf("error: not a directory %q", uri)
		}
		return &lazyio.FileReader{Dir: uri}, nil
	}
}

func newOutputOpener(io *sys.IO, uri string) (lazyio.OutputOpener, error) {
	log := io.Log
	switch {
	case strings.HasPrefix(uri, "s3://") || strings.HasPrefix(uri, "S3://"):
		log.Debugw("creating S3 output writer", "uri", uri)
		return lazyio.NewS3Writer(io, uri)
	default:
		log.Debugw("creating FS output writer", "uri", uri)
		uri = filepath.Clean(uri)
		if info, err := os.Stat(uri); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
			if err = os.MkdirAll(uri, 0777); err != nil {
				return nil, err
			}
		} else if !info.IsDir() {
			return nil, fmt.Errorf("error: not a directory %q", uri)
		}
		return &lazyio.FileWriter{Dir: uri}, nil
	}
}
