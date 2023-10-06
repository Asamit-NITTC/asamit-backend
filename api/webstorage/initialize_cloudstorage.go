package webstorage

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

func InitializeCloudStorage() (context.Context, *storage.BucketHandle) {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name for the new bucket.
	bucketName := "asamit-image"

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	return ctx, bucket
}
