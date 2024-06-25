package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/storage/v1"
)

const (
	bucketName = "qwiklabs-gcp-01-afc66f30517d-bucket"
	objectName = "accounts-data.txt" // Name of the file/object in the bucket
)

var (
	client *storage.Service
)

type Transaction struct {
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
	ClientID string `json:"client_id"`
}

type AccountData struct {
	Transactions []Transaction `json:"transactions"`
}

var mu sync.Mutex

func handleMonitor(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	ctx := context.Background()

	// Initialize Google Cloud Storage client
	initializeGCSClient(ctx)

	// Read from the object in GCS
	var accountData AccountData
	accountData, err := readFromGCS(ctx, client, bucketName, objectName)
	if err != nil {
		log.Fatalf("Failed to read from GCS: %v", err)
	}

	// // Print the parsed data
	// fmt.Printf("Name: %s\nAge: %d\nEmail: %s\n", transaction.ClientID, transaction.Amount, transaction.Balance)

	// fmt.Printf("Contents of %s:\n%s\n", objectName, string(data))
	json.NewEncoder(w).Encode(accountData)
}

// initializeGCSClient initializes the Google Cloud Storage client
func initializeGCSClient(ctx context.Context) {
	var err error
	client, err = storage.NewService(ctx, option.WithCredentialsFile("../service-account-file.json"))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}
	log.Println("Google Cloud Storage client initialized")
}

// readFromGCS reads data from a file in GCS bucket
func readFromGCS(ctx context.Context, storageService *storage.Service, bucketName, objectName string) (AccountData, error) {

	var accountData AccountData
	resp, err := storageService.Objects.Get(bucketName, objectName).Context(ctx).Download()
	if err != nil {
		return accountData, fmt.Errorf("failed to download object: %v", err)
	}
	defer resp.Body.Close()

	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read object body: %v", err)
	// }
	scanner := bufio.NewScanner(resp.Body)
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

	return accountData, nil
	// return data, nil
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/app/public")))
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/monitor", handleMonitor)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
