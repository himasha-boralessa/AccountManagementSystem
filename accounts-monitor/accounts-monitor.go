package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Transaction struct {
	Time   string `json:"time"`
	Amount int    `json:"amount"`
}

type Account struct {
	Balance      int           `json:"balance"`
	Transactions []Transaction `json:"transactions"`
}

var (
	mu sync.Mutex
)

func getTransactions(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://account-system:8082/transactions")
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var account Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		http.Error(w, "Failed to decode transactions", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(account)
}

func main() {
	http.HandleFunc("/transactions", getTransactions)
	log.Fatal(http.ListenAndServe(":8083", nil))
}
