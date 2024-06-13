package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Account struct {
	ID      string   `json:"id"`
	Balance int      `json:"balance"`
	Logs    []string `json:"logs"`
}

var (
	accounts = make(map[string]*Account)
	mu       sync.Mutex
)

func getAccounts(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(accounts)
}

func updateAccount(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	account, exists := accounts[id]
	if !exists {
		account = &Account{ID: id, Logs: []string{}}
		accounts[id] = account
	}
	account.Balance += amount
	account.Logs = append(account.Logs, fmt.Sprintf("%s: %d", time.Now().Format(time.RFC3339), amount))
	fmt.Fprintf(w, "Account %s updated, new balance: %d", id, account.Balance)
}

func main() {
	http.HandleFunc("/accounts", getAccounts)
	http.HandleFunc("/update", updateAccount)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
