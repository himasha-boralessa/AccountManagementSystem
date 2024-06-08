package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const accountManagerURL = "http://account-manager-service/transaction"

type TransactionRequest struct {
	Amount float64 `json:"amount"`
}

func sendTransaction(amount float64) error {
	requestBody, err := json.Marshal(TransactionRequest{Amount: amount})
	if err != nil {
		return err
	}

	resp, err := http.Post(accountManagerURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send transaction: %s", resp.Status)
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		amount := (rand.Float64() * 200) - 100 // Random amount between -100 and 100
		fmt.Printf("Sending transaction of amount: %.2f\n", amount)

		err := sendTransaction(amount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending transaction: %v\n", err)
		}

		sleepDuration := time.Duration(rand.Intn(10)) * time.Second
		time.Sleep(sleepDuration)
	}
}
