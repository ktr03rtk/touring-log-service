package persistence

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ktr03rtk/touring-log-service/api-backend/domain/repository"
	"github.com/pkg/errors"
)

type photoImagePersistence struct {
	bucket string
	*s3.Client
}

func NewPhotoImagePersistence(ctx context.Context, region, bucket string) (repository.PhotoImageRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, errors.Wrap(err, "failed to configure aws")
	}

	return &photoImagePersistence{
		bucket,
		s3.NewFromConfig(cfg),
	}, nil
}

func (pp *photoImagePersistence) Store(reader io.Reader, key string) error {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if _, err := io.Copy(zw, reader); err != nil {
		return errors.Wrapf(err, "failed to gzip write")
	}

	if err := zw.Close(); err != nil {
		return errors.Wrapf(err, "failed to close gzip writer")
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(pp.bucket),
		Key:    aws.String(key),
		Body:   &buf,
	}

	if _, err := pp.PutObject(context.Background(), input); err != nil {
		return errors.Wrapf(err, "failed to put object")
	}

	return nil
}
