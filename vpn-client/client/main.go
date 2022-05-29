package main

import (
	"context"
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

// var (
// 	sourceDir      string
// 	uploadInterval time.Duration
// )

// func init() {
// 	if err := getEnv(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func getEnv() error {
// 	d, ok := os.LookupEnv("SOURCE_DIRECTORY")
// 	if !ok {
// 		return errors.New("env SOURCE_DIRECTORY is not found")
// 	}

// 	sourceDir = d

// 	r, ok := os.LookupEnv("UPLOAD_INTERVAL_SECONDS")
// 	if !ok {
// 		return errors.New("env UPLOAD_INTERVAL_SECONDS is not found")
// 	}

// 	i, err := strconv.Atoi(r)
// 	if err != nil {
// 		return errors.New("env UPLOAD_INTERVAL_SECONDS is not integer")
// 	}

// 	uploadInterval = time.Duration(i) * time.Second

// 	return nil
// }

func main() {
	region := "ap-northeast-1"
	bucket := "foo"
	key := "client-ovpn/client.ovpn"
	distDir := "/opt/vpn_client/"
	fileName := "client.ovpn"
	retry_config := "\nconnect-retry-max 10"

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
