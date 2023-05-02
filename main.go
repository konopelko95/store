package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Order struct {
	Currency string `json:"currency"`
	Items    int    `json:"items"`
}

type Payment struct {
	CardNumber   string `json:"cardNumber"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Expiry       string `json:"expiry"`
	CVV          string `json:"cvv"`
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current working directory:", dir)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	db, err := sql.Open("postgres", "postgres://hrtivrrt:KghCpKhtC8lJj-6LK6nzn5ZfRxXroMM5@john.db.elephantsql.com/hrtivrrt")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Create a wrapper struct to correctly decode the JSON request body
    type RequestWrapper struct {
        Order   *Order   `json:"Order"`
        Payment *Payment `json:"Payment"`
    }

    // Decode JSON request body into the RequestWrapper struct
    var reqWrapper RequestWrapper
    err := json.NewDecoder(r.Body).Decode(&reqWrapper)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    order := reqWrapper.Order
    payment := reqWrapper.Payment
    amount := order.Items * 50


    // Log order information to console
	fmt.Printf("Payment: %v\n", *payment)
    fmt.Printf("Amount: %d %s\n", amount, order.Currency)



    // Send the response to the client
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": fmt.Sprintf("Ordered %d items\n Total amount: %d %s", order.Items, amount, order.Currency),
    })
})


	log.Fatal(http.ListenAndServe(":8080", nil))
}