// package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"time"
// )

// const accountManagerURL = "http://account-manager-service/transaction"

// type TransactionRequest struct {
// 	Amount float64 `json:"amount"`
// }

// func sendTransaction(amount float64) error {
// 	requestBody, err := json.Marshal(TransactionRequest{Amount: amount})
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := http.Post(accountManagerURL, "application/json", bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to send transaction: %s", resp.Status)
// 	}

// 	return nil
// }

// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	for {
// 		amount := (rand.Float64() * 200) - 100 // Random amount between -100 and 100
// 		fmt.Printf("Sending transaction of amount: %.2f\n", amount)

// 		err := sendTransaction(amount)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error sending transaction: %v\n", err)
// 		}

// 		sleepDuration := time.Duration(rand.Intn(10)) * time.Second
// 		time.Sleep(sleepDuration)
// 	}
// }

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{}
	// Read the CLIENT_ID environment variable
	clientID := os.Getenv("CLIENT_ID")
	if clientID == "" {
		// fmt.Println("CLIENT_ID environment variable not set")
		// return
		clientID = "client2"
	}

	for {
		amount := rand.Intn(200) - 100 // Random amount between -100 and +100
		req, err := http.NewRequest("POST", "http://localhost:8082/transaction", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}
		q := req.URL.Query()
		q.Add("amount", fmt.Sprintf("%d", amount))
		req.URL.RawQuery = q.Encode()

		// Add Client-ID header to identify the client
		req.Header.Set("Client-ID", clientID)

		_, err = client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
		}

		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}
