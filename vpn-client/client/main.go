package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/ktr03rtk/touring-log-service/vpn-client/client/pkg/downloader"
)

var (
	region string
	bucket string
	key    string
)

const (
	distDir  = "/opt/vpn_client/"
	fileName = "client.ovpn"
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

	client, err := downloader.NewDownloader(ctx, region)
	if err != nil {
		log.Print(err)
		os.Exit(0)
	}

	if err := client.DownloadNewVpnConfig(ctx, bucket, key, distDir, fileName); err != nil {
		log.Print(err)
		os.Exit(0)
	}

	if err := exec.Command("openvpn", "--config", distDir+fileName).Run(); err != nil {
		log.Print(err)
		os.Exit(0)
	}
}
