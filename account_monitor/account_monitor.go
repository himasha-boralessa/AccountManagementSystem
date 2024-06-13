package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
    balance int
    mu      sync.Mutex
)

type Transaction struct {
    Time   string `json:"time"`
    Amount int    `json:"amount"`
}

func handleTransaction(w http.ResponseWriter, r *http.Request) {
    amountStr := r.URL.Query().Get("amount")
    amount, err := strconv.Atoi(amountStr)
    if err != nil {
        http.Error(w, "Invalid amount", http.StatusBadRequest)
        return
    }

    mu.Lock()
    balance += amount
    mu.Unlock()

    response := map[string]interface{}{
        "balance": balance,
        "transaction": Transaction{
            Time:   time.Now().Format(time.RFC3339),
            Amount: amount,
        },
    }
    json.NewEncoder(w).Encode(response)
}

func main() {
    http.HandleFunc("/transaction", handleTransaction)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
