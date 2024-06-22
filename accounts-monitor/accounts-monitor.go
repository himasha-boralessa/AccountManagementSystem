package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/storage"
)

type Transaction struct {
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
	ClientID string `json:"client_id"`
}

type AccountData struct {
	Transactions []Transaction `json:"transactions"`
}

const dataFilePath = "/app/data/account-data.txt"

// const dataFilePath = "../account-data.txt"

var mu sync.Mutex

func handleMonitor(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer client.Close()
	// Specify your bucket and object (file) name
	bucketName := "my-test-bucket"
	objectName := "example.txt"
	// Read from the object in GCS
	data, err := readFromGCS(ctx, client, bucketName, objectName)
	if err != nil {
		log.Fatalf("Failed to read from GCS: %v", err)
	}
	fmt.Printf("Contents of %s:\n%s\n", objectName, string(data))

	// accountData := loadAccountData()
	// json.NewEncoder(w).Encode(accountData)
}

func loadAccountData() AccountData {
	file, err := os.Open(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return AccountData{}
		}
		log.Fatalf("Failed to open account file: %v", err)
	}
	defer file.Close()

	var accountData AccountData
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var transaction Transaction
		if err := json.Unmarshal(scanner.Bytes(), &transaction); err != nil {
			log.Fatalf("Failed to unmarshal transaction data: %v", err)
		}
		accountData.Transactions = append(accountData.Transactions, transaction)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read account file: %v", err)
	}

	return accountData
}

// readFromGCS reads data from a file in GCS bucket.
func readFromGCS(ctx context.Context, client *storage.Client, bucketName, objectName string) ([]byte, error) {
	// Get a bucket handle
	bucket := client.Bucket(bucketName)

	// Get a reader for the object in GCS
	obj := bucket.Object(objectName)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS reader: %v", err)
	}
	defer reader.Close()

	// Read data from the object in GCS
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read from GCS object: %v", err)
	}

	return data, nil
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/app/public")))
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/monitor", handleMonitor)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
