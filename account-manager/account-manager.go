package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const storagePath = "/app/storage"

func saveJSONToFile(jsonData, filename string) error {
	filePath := storagePath + "/" + filename
	return ioutil.WriteFile(filePath, []byte(jsonData), 0644)
}

func readJSONFromFile(filename string) (string, error) {
	filePath := storagePath + "/" + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", nil // Or return an error or default value
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type Transaction struct {
	Time   string `json:"time"`
	Amount int    `json:"amount"`
}

type Account struct {
	Balance      int           `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

var (
	accountFile = "account.json"
	account     = Account{
		Balance:      0,
		Transactions: []Transaction{},
	}
	mu sync.Mutex
)

func loadAccount() {
	data, err := readJSONFromFile(accountFile)
	if err != nil {
		log.Fatalf("Failed to read account file: %v", err)
	}
	if data != "" {
		err = json.Unmarshal([]byte(data), &account)
		if err != nil {
			log.Fatalf("Failed to parse account file: %v", err)
		}
	}
}

func saveAccount() {
	data, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize account: %v", err)
	}
	err = saveJSONToFile(string(data), accountFile)
	if err != nil {
		log.Fatalf("Failed to write account file: %v", err)
	}
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	mu.Lock()
	account.Balance += amount
	transaction := Transaction{
		Time:   time.Now().Format(time.RFC3339),
		Amount: amount,
	}
	account.Transactions = append(account.Transactions, transaction)
	saveAccount()
	mu.Unlock()

	response := map[string]interface{}{
		"balance":      account.Balance,
		"transactions": account.Transactions,
	}
	json.NewEncoder(w).Encode(response)
}

func getTransactions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(account)
}

func main() {
	loadAccount()

	http.HandleFunc("/transaction", handleTransaction)
	http.HandleFunc("/transactions", getTransactions)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
