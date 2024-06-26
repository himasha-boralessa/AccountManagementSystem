package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type Transaction struct {
	Time     string `json:"time"`
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
	ClientID string `json:"client_id"`
}

var (
	balance1 int
	balance2 int
	mu       sync.Mutex
)

const (
	bucketName = "qwiklabs-gcp-00-8af698a525cb-bucket"
	objectName = "data.txt"
)

func main() {
	http.HandleFunc("/transaction", handleTransaction)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Retrieve Client-ID from header
	clientID := r.Header.Get("Client-ID")
	if clientID == "" {
		http.Error(w, "clientid query parameter is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	transaction := Transaction{
		Time:     time.Now().Format(time.RFC3339),
		Amount:   amount,
		ClientID: clientID,
	}

	if clientID == "client1" {
		balance1 += transaction.Amount
		transaction.Balance = balance1
	} else if clientID == "client2" {
		balance2 += transaction.Amount
		transaction.Balance = balance2
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		http.Error(w, "Failed to create storage client", http.StatusInternalServerError)
		return
	}

	err = appendToGCS(ctx, client, bucketName, objectName, transaction)
	if err != nil {
		http.Error(w, "Failed to write to GCS", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Transaction of %d, new balance: %d, Client-ID: %s\n", transaction.Amount, transaction.Balance, transaction.ClientID)
}

func appendToGCS(ctx context.Context, client *storage.Client, bucketName, objectName string, transaction Transaction) error {
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// Read existing content
	r, err := obj.NewReader(ctx)
	if err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("failed to read object: %v", err)
	}
	var content string
	if err == nil {
		defer r.Close()
		body, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to read object body: %v", err)
		}
		content = string(body)
		// log.Println("Read existing content from GCS:", content)
	}

	// // Append new data
	// newData := fmt.Sprintf("{\"time\":\"%s\",\"amount\":%d,\"client_id\":\"%s\"}\n", time.Now().Format(time.RFC3339), amount, clientID)
	// content += newData

	// Marshal transaction to JSON
	// transactionData, err := json.Marshal(transaction)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal transaction data: %v", err)
	// }

	// Append new transaction data
	newData := fmt.Sprintf("{\"amount\":\"%d\",\"balance\":%d,\"client_id\":\"%s\"}\n", transaction.Amount, transaction.Balance, transaction.ClientID)
	content += string(newData) + "\n"

	// Write back to GCS
	w := obj.NewWriter(ctx)
	defer w.Close()

	_, err = w.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("failed to write to GCS: %v", err)
	}

	return nil
}
