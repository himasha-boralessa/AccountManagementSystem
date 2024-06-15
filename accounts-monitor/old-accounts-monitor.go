// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"sync"
// 	"time"
// )

// var (
// 	balance int
// 	mu      sync.Mutex
// )

// type Transaction struct {
// 	Time   string `json:"time"`
// 	Amount int    `json:"amount"`
// }

// func handleTransaction(w http.ResponseWriter, r *http.Request) {
// 	// amountStr := r.URL.Query().Get("amount")
// 	// amount, err := strconv.Atoi(amountStr)
// 	// if err != nil {
// 	// 	http.Error(w, "Invalid amount", http.StatusBadRequest)
// 	// 	return
// 	// }
// 	amount := 10000
// 	balance := 500

// 	mu.Lock()
// 	balance += amount
// 	mu.Unlock()

// 	// response := map[string]interface{}{
// 	// 	"balance": balance,
// 	// 	"transactions": Transaction{
// 	// 		Time:   time.Now().Format(time.RFC3339),
// 	// 		Amount: amount,
// 	// 	},
// 	// }

// 	response := map[string]interface{}{
// 		"balance": balance,
// 		"transactions": []Transaction{
// 			{
// 				Time:   time.Now().Format(time.RFC3339),
// 				Amount: 500000},
// 			{
// 				Time:   time.Now().Format(time.RFC3339),
// 				Amount: 60000000},
// 			{
// 				Time:   time.Now().Format(time.RFC3339),
// 				Amount: 800000},
// 		},
// 	}

// 	json.NewEncoder(w).Encode(response)
// }

// func main() {
// 	http.HandleFunc("/transaction", handleTransaction)
// 	http.Handle("/", http.FileServer(http.Dir("./public")))
// 	log.Fatal(http.ListenAndServe(":8083", nil))
// }
