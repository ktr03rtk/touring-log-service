package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	region string
	bucket string
	key    string
)

const (
	distDir      = "/opt/vpn_client/"
	fileName     = "client.ovpn"
	retry_config = "\nconnect-retry-max 10"
)

func init() {
	if err := getEnv(); err != nil {
		log.Fatal(err)
	}
}

func getEnv() error {
	r, ok := os.LookupEnv("REGION")
	if !ok {
		return errors.New("env REGION is not found")
	}

	region = r

	b, ok := os.LookupEnv("BUCKET")
	if !ok {
		return errors.New("env BUCKET is not found")
	}

	bucket = b

	k, ok := os.LookupEnv("KEY")
	if !ok {
		return errors.New("env KEY is not found")
	}

	key = k

	return nil
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	output, err := client.GetObject(ctx, input)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Body.Close()

	b, err := io.ReadAll(output.Body)
	if err != nil {
		log.Fatal(err)
		// return nil, errors.Wrapf(err, "failed to read file")
	}

	filePath := filepath.Join(distDir, fileName)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
		// return errors.Wrapf(err, "failed to create file")
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, string(b), retry_config); err != nil {
		log.Fatal(err)
		// return errors.Wrapf(err, "failed to write file")
	}

	if err := exec.Command("openvpn", "--config", distDir+fileName).Run(); err != nil {
		log.Fatal(err)
	}
}
