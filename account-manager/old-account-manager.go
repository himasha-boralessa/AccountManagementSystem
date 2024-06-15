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
