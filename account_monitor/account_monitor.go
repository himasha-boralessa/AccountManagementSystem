// accounts_monitor.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type AccountStatus struct {
	Balance float64  `json:"balance"`
	Logs    []string `json:"logs"`
}

type AccountsMonitor struct {
	mu       sync.Mutex
	accounts map[string]AccountStatus
}

func (am *AccountsMonitor) updateAccounts() {
	// Assuming the account manager is running on account-manager:8080
	urls := []string{"http://account-manager-1:8080/logs", "http://account-manager-2:8080/logs"}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching logs:", err)
			continue
		}
		defer resp.Body.Close()

		var logs []string
		if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
			fmt.Println("Error decoding logs:", err)
			continue
		}

		balance := 0.0
		for _, logEntry := range logs {
			// Simple parsing to extract balance from logs
			fmt.Sscanf(logEntry, "%*s - Transaction: %*f, New Balance: %f", &balance)
		}

		am.mu.Lock()
		am.accounts[url] = AccountStatus{Balance: balance, Logs: logs}
		am.mu.Unlock()
	}
}

func (am *AccountsMonitor) handleStatus(w http.ResponseWriter, r *http.Request) {
	am.mu.Lock()
	defer am.mu.Unlock()

	json.NewEncoder(w).Encode(am.accounts)
}

func main() {
	monitor := &AccountsMonitor{accounts: make(map[string]AccountStatus)}

	go func() {
		for {
			monitor.updateAccounts()
			time.Sleep(10 * time.Second)
		}
	}()

	http.HandleFunc("/status", monitor.handleStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
