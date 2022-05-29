package downloader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

const retry_config = "\nconnect-retry-max 10"

type Downloader struct {
	*s3.Client
}

func NewDownloader(ctx context.Context, region string) (*Downloader, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load config")
	}

	return &Downloader{
		s3.NewFromConfig(cfg),
	}, nil
}

func (d *Downloader) DownloadNewVpnConfig(ctx context.Context, bucket, key, distDir, fileName string) error {

	if err := d.download(ctx, bucket, key, distDir, fileName); err != nil {
		return err
	}

	if err := d.delete(ctx, bucket, key); err != nil {
		return err
	}

	return nil
}

func (d *Downloader) download(ctx context.Context, bucket string, key string, distDir string, fileName string) error {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	output, err := d.GetObject(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "failed to get object")
	}
	defer output.Body.Close()

	b, err := io.ReadAll(output.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to read object")
	}

	filePath := filepath.Join(distDir, fileName)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		return errors.Wrapf(err, "failed to create file")
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, string(b), retry_config); err != nil {
		return errors.Wrapf(err, "failed to write file")
	}
	return nil
}

func (d *Downloader) delete(ctx context.Context, bucket string, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := d.DeleteObject(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "failed to delete object")
	}

	return nil
}
