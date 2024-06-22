package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
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

	accountData := loadAccountData()
	json.NewEncoder(w).Encode(accountData)
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

func main() {
	http.Handle("/", http.FileServer(http.Dir("/app/public")))
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/monitor", handleMonitor)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
