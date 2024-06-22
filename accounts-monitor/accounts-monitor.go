package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

const (
	projectID  = "your-project-id"
	bucketName = "your-bucket-name"
	objectName = "your-object-name" // Name of the file/object in the bucket
)

var (
	client *storage.Service
)

func main() {
	ctx := context.Background()

	// Initialize Google Cloud Storage client
	initializeGCSClient(ctx)

	// Read from the object in GCS
	data, err := readFromGCS(ctx, client, bucketName, objectName)
	if err != nil {
		log.Fatalf("Failed to read from GCS: %v", err)
	}

	fmt.Printf("Contents of %s:\n%s\n", objectName, string(data))
}

// initializeGCSClient initializes the Google Cloud Storage client
func initializeGCSClient(ctx context.Context) {
	var err error
	client, err = storage.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	log.Println("Google Cloud Storage client initialized")
}

// readFromGCS reads data from a file in GCS bucket
func readFromGCS(ctx context.Context, storageService *storage.Service, bucketName, objectName string) ([]byte, error) {
	resp, err := storageService.Objects.Get(bucketName, objectName).Context(ctx).Download()
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %v", err)
	}

	return data, nil
}
