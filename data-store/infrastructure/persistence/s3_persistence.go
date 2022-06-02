package persistence

import (
	"bytes"
	"compress/gzip"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/repository"
	"github.com/pkg/errors"
)

type s3Persistence struct {
	bucket string
	*s3.Client
}

func NewS3Persistence(ctx context.Context, region, bucket string) (repository.PayloadStoreRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure aws")
	}

	return &s3Persistence{
		bucket,
		s3.NewFromConfig(cfg),
	}, nil
}

func (sp *s3Persistence) Store(ctx context.Context, payload *model.Payload) error {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if _, err := zw.Write(payload.GetMessage()); err != nil {
		return errors.Wrapf(err, "failed to gzip message")
	}

	if err := zw.Close(); err != nil {
		return errors.Wrapf(err, "failed to close gzip writer")
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(strings.TrimSuffix(sp.bucket, "dat") + "gz"),
		Key:    aws.String(payload.GetKey()),
		Body:   &buf,
	}

	if _, err := sp.PutObject(ctx, input); err != nil {
		return errors.Wrapf(err, "failed to put object")
	}

	return nil
}
