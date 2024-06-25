package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	bucketName = "qwiklabs-gcp-01-afc66f30517d-bucket"
	objectName = "accounts-data.txt" // Name of the file/object in the bucket
)

var (
	ctx    context.Context
	client *storage.Client
)

func main() {
	// Initialize Google Cloud Storage client
	initializeGCSClient()

	// Example data to append
	dataToAppend := `{"time":"2024-06-19T07:57:01+02:00","amount":8888888888,"balance":647,"client_id":"client1"}`

	// Append data to the text file in the bucket
	err := appendToGCS(ctx, client, bucketName, objectName, []byte(dataToAppend))
	if err != nil {
		log.Fatalf("Failed to append to GCS: %v", err)
	}

	fmt.Printf("Data appended to %s in bucket %s\n", objectName, bucketName)
}

// initializeGCSClient initializes the Google Cloud Storage client
func initializeGCSClient() {
	var err error
	ctx = context.Background()

	// Use service account key for authentication
	client, err = storage.NewClient(ctx, option.WithCredentialsFile("../service-account-file.json"))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	log.Println("Google Cloud Storage client initialized")
}

// appendToGCS appends data to a file in GCS bucket
func appendToGCS(ctx context.Context, storageClient *storage.Client, bucketName, objectName string, data []byte) error {
	// Read existing data from the object
	existingData, err := readFromGCS(ctx, storageClient, bucketName, objectName)
	if err != nil {
		// If the file does not exist, create it
		if err == storage.ErrObjectNotExist {
			existingData = []byte{}
		} else {
			return fmt.Errorf("failed to read existing data: %v", err)
		}
	}

	// Append the new data to the existing data
	updatedData := append(existingData, data...)

	// Write the updated data back to the object in GCS
	return writeToGCS(ctx, storageClient, bucketName, objectName, updatedData)
}

// readFromGCS reads data from a file in GCS bucket
func readFromGCS(ctx context.Context, storageClient *storage.Client, bucketName, objectName string) ([]byte, error) {
	// Open a reader for the object in GCS
	rc, err := storageClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %v", err)
	}
	defer rc.Close()

	// Read the object data
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %v", err)
	}

	return data, nil
}

// writeToGCS writes data to a file in GCS bucket
func writeToGCS(ctx context.Context, storageClient *storage.Client, bucketName, objectName string, data []byte) error {
	// Open a writer for the object in GCS
	wc := storageClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	// Write data to the object
	if _, err := wc.Write(data); err != nil {
		return fmt.Errorf("failed to write object: %v", err)
	}

	// Close the writer to flush the data to GCS
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}
