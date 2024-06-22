package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	mu sync.Mutex
)

type Transaction struct {
	Time     string `json:"time"`
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
	ClientID string `json:"client_id"`
}

const dataFilePath = "/app/data/account-data.txt"

// const dataFilePath = "../account-data.txt"

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
		http.Error(w, "Client-ID header missing", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	balance := calculateBalance(clientID)
	transaction := Transaction{
		Time:     time.Now().Format(time.RFC3339),
		Balance:  balance,
		Amount:   amount,
		ClientID: clientID,
	}

	appendTransactionToFile(transaction)

	log.Printf("Transaction of %d, new balance: %d, Client-ID: %s\n", amount, transaction.Balance, clientID)
}

func appendTransactionToFile(transaction Transaction) {
	file, err := os.OpenFile(dataFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open account file: %v", err)
	}
	defer file.Close()

	transactionData, err := json.Marshal(transaction)
	if err != nil {
		log.Fatalf("Failed to marshal transaction data: %v", err)
	}

	if _, err := file.WriteString(string(transactionData) + "\n"); err != nil {
		log.Fatalf("Failed to write to account file: %v", err)
	}
}

func calculateBalance(clientID string) int {

	file, err := os.Open(dataFilePath)
	if err != nil {
		log.Fatalf("Failed to open account file: %v", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	balance := 0
	for scanner.Scan() {
		var transaction Transaction
		if err := json.Unmarshal(scanner.Bytes(), &transaction); err != nil {
			log.Fatalf("Failed to unmarshal transaction data: %v", err)
		}

		if clientID == transaction.ClientID {
			balance += transaction.Amount
		}
	}

	return balance
}

func main() {
	os.MkdirAll("/app", os.ModePerm) // Ensure the directory exists
	// os.MkdirAll("../", os.ModePerm)  // Ensure the parent directory exists
	http.HandleFunc("/transaction", handleTransaction)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
