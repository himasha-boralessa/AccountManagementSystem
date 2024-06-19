// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"sync"
// 	"time"
// )

// var (
// 	balance int
// 	mu      sync.Mutex
// )

// func handleTransaction(w http.ResponseWriter, r *http.Request) {
// 	amountStr := r.URL.Query().Get("amount")
// 	amount, err := strconv.Atoi(amountStr)
// 	if err != nil {
// 		http.Error(w, "Invalid amount", http.StatusBadRequest)
// 		return
// 	}

// 	mu.Lock()
// 	defer mu.Unlock()
// 	balance += amount

// 	log.Printf("%s: Transaction of %d, new balance: %d\n", time.Now().Format(time.RFC3339), amount, balance)
// 	fmt.Fprintf(w, "Transaction successful, new balance: %d", balance)
// }

// func main() {
// 	http.HandleFunc("/transaction", handleTransaction)
// 	log.Fatal(http.ListenAndServe(":8082", nil))
// }

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
	mu       sync.Mutex
	balance1 int
	balance2 int
)

type Transaction struct {
	Time     string `json:"time"`
	Amount   int    `json:"amount"`
	Balance  int    `json:"balance"`
	ClientID string `json:"client_id"`
}

// const dataFilePath = "/app/account-data/account.txt"
const dataFilePath = "../account-data.txt"

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

	// balance := calculateBalance(clientID)
	transaction := Transaction{
		Time:     time.Now().Format(time.RFC3339),
		Amount:   amount,
		ClientID: clientID,
	}

	if clientID == "client1" {
		balance1 += transaction.Amount
		transaction.Balance = balance1
	}
	if clientID == "client2" {
		balance2 += transaction.Amount
		transaction.Balance = balance2
	}

	appendTransactionToFile(transaction)

	log.Printf("Transaction of %d, new balance: %d, Client-ID: %s\n", amount, transaction.Balance, clientID)
	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"balance":     balance,
	// 	"transaction": transaction,
	// })
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
	// os.MkdirAll("/app/account-data", os.ModePerm) // Ensure the directory exists
	os.MkdirAll("../", os.ModePerm) // Ensure the parent directory exists
	http.HandleFunc("/transaction", handleTransaction)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
