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
		fmt.Println("CLIENT_ID environment variable not set")
		return
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
