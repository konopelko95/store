package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

	log.Fatal(http.ListenAndServe(":8080", nil))

	db, err := sql.Open("postgres", "postgres://user:password@hostname:port/database?sslmode=require")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.Handle("/", fs)

	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var order Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		amount := order.Items * 50

		// Log order information to console
		fmt.Printf("Order: %v\n", order)
		fmt.Printf("Amount: %d %s\n", amount, order.Currency)

		// Log order information to database
		statement, err := db.Prepare("INSERT INTO orders(date, time, sum, card_number, name, surname, expiry, cvv) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer statement.Close()

		now := time.Now()
		date := now.Format("2006-01-02")
		t := now.Format("15:04:05")

		payment := Payment{}
		err = json.NewDecoder(r.Body).Decode(&payment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = statement.Exec(date, t, amount, payment.CardNumber, payment.Name, payment.Surname, payment.Expiry, payment.CVV)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Purchase successful"})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
