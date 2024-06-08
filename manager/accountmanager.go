// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"sync"
// 	"time"
// )

// type TransactionRequest struct {
// 	Amount float64 `json:"amount"`
// }

// type AccountManager struct {
// 	balance float64
// 	mu      sync.Mutex
// 	logFile *os.File
// }

// func NewAccountManager(logFilePath string) (*AccountManager, error) {
// 	// file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	return &AccountManager{
// 		balance: 0,
// 		// logFile: file,
// 	}, nil
// }

// func (am *AccountManager) handleTransaction(w http.ResponseWriter, r *http.Request) {
// 	var req TransactionRequest

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	am.mu.Lock()
// 	defer am.mu.Unlock()

// 	am.balance += req.Amount
// 	logEntry := fmt.Sprintf("%s - Transaction: %.2f, New Balance: %.2f\n", time.Now().Format(time.RFC3339), req.Amount, am.balance)
// 	if _, err := am.logFile.WriteString(logEntry); err != nil {
// 		http.Error(w, "Failed to log transaction", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func (am *AccountManager) Start() {
// 	http.HandleFunc("/transaction", am.handleTransaction)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func main() {
// 	accountManager, _ := NewAccountManager("/var/log/account_manager.log")
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to create account manager: %v", err)
// 	// }
// 	// defer accountManager.logFile.Close()

// 	fmt.Println("Account manager started")
// 	accountManager.Start()
// }

package main

import (
	"fmt"
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

func handleTransaction(w http.ResponseWriter, r *http.Request) {
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	balance += amount

	log.Printf("%s: Transaction of %d, new balance: %d\n", time.Now().Format(time.RFC3339), amount, balance)
	fmt.Fprintf(w, "Transaction successful, new balance: %d", balance)
}

func main() {
	http.HandleFunc("/transaction", handleTransaction)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
