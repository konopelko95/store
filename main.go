package main

import (
    "fmt"
    "net/http"
)

//This handler function takes two parameters: a ResponseWriter object that we can use to send the HTTP response,
// and a Request object that represents the HTTP request.

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

// Here, we register the helloHandler function as the handler for the root URL ("/") using the http.HandleFunc function.

// Then, we start the web server on port 8080 using the http.ListenAndServe function. 
// The second parameter is the handler to use for incoming requests, and we pass nil to use the default handler (which is the http.DefaultServeMux).

func main() {
    http.HandleFunc("/", helloHandler)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("Failed to start server:", err)
    }
}

// In this example, we define a new function called buyBitcoin that takes two parameters: amount (the amount of Bitcoin to buy) and currency (the currency to use for payment).
// The function constructs a JSON request body with the necessary parameters, sets the API key and secret in the request headers, 
// and sends an HTTP POST request to the /api/v1/trade/placeOrder endpoint.
// After sending the request, the function reads the response body and prints it to the console.

// Note that this is a simplified example and you may need to modify the code to fit your specific use case. 
// You should also review the BitForex API documentation to understand the full set of API endpoints and parameters available.
func buyBitcoin(amount float64, currency string) error {
    // Set up the API endpoint URL
    url := "https://api.bitforex.com/api/v1/trade/placeOrder"

    // Set up the request body
    body := fmt.Sprintf(`{"symbol": "BTC/%s", "type": "BUY", "price": 0, "amount": %f, "exchangeType": "RETAIL", "payMethod": "BANK_CARD"}`, currency, amount)

    // Create a new HTTP request
    req, err := http.NewRequest("POST", url, strings.NewReader(body))
    if err != nil {
        return err
    }

    // Set the API key and secret in the request headers
    req.Header.Set("BF-API-KEY", "<YOUR_API_KEY>")
    req.Header.Set("BF-API-SECRET", "<YOUR_API_SECRET>")

    // Send the HTTP request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Read the response body
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // Print the response body
    fmt.Println(string(respBody))

    return nil
}

func buyHandler(w http.ResponseWriter, r *http.Request) {
    // Get the form data from the URL parameters
    currency := r.URL.Query().Get("currency")
    amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
    if err != nil {
        http.Error(w, "Invalid amount", http.StatusBadRequest)
        return
    }

    // Buy the cryptocurrency using the BitForex API
    err = buyBitcoin(amount, currency)
    if err != nil {
        http.Error(w, "Failed to buy cryptocurrency", http.StatusInternalServerError)
        return
    }

    // Return a JSON response with a success message
    response := struct {
        Message string `json:"message"`
    }{
        Message: fmt.Sprintf("Successfully bought %f %s of cryptocurrency", amount, currency),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


