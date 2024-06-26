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
	bucketName = "qwiklabs-gcp-02-7e208e1db0ca-bucket"
	objectName = "data.txt"
)

var (
	client *storage.Service
	mu     sync.Mutex
)

type Transaction struct {
	Time     string `json:"time"`
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
	ClientID string `json:"client_id"`
}

type AccountData struct {
	Transactions []Transaction `json:"transactions"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("/app/public")))
	http.HandleFunc("/monitor", handleMonitor)
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func handleMonitor(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	ctx := context.Background()
	initializeGCSClient(ctx)

	accountData, err := readFromGCS(ctx, client, bucketName, objectName)
	if err != nil {
		http.Error(w, "Failed to read from GCS", http.StatusInternalServerError)
		log.Fatalf("Failed to read from GCS: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accountData); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
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
func readFromGCS(ctx context.Context, storageService *storage.Service, bucketName, objectName string) (AccountData, error) {
	var accountData AccountData

	resp, err := storageService.Objects.Get(bucketName, objectName).Context(ctx).Download()
	if err != nil {
		return accountData, fmt.Errorf("failed to download object: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}
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
}
