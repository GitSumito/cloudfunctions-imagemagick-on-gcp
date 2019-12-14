package imagemagick

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
)

func init() {
	var err error

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// ImageConvert blurs offensive images uploaded to GCS.
func BlurOffensiveImages(ctx context.Context, e GCSEvent) error {
	outputBucket := os.Getenv("THUMBNAILED")
	if outputBucket == "" {
		return errors.New("THUMBNAILED must be set")
	}
	return blur(ctx, e.Bucket, outputBucket, e.Name)
}

// [END functions_imagemagick_analyze]

// image convert from gs://inputBucket/name to gs://outputBucket/name.
func blur(ctx context.Context, inputBucket, outputBucket, name string) error {
	inputBlob := storageClient.Bucket(inputBucket).Object(name)
	r, err := inputBlob.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("NewReader: %v", err)
	}

	outputBlob := storageClient.Bucket(outputBucket).Object(name)
	w := outputBlob.NewWriter(ctx)
	defer w.Close()

	// Use - as input and output to use stdin and stdout.
	cmd := exec.Command("convert", "-", "-resize", "30%", "-")
	cmd.Stdin = r
	cmd.Stdout = w

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %v", err)
	}

	log.Printf("Blurred image uploaded to gs://%s/%s", outputBlob.BucketName(), outputBlob.ObjectName())

	return nil
}

// [END functions_imagemagick_blur]
