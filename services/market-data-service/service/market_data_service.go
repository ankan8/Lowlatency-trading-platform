package service

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    notificationpb "github.com/ankan8/swapsync/backend/services/notification-service/proto"
    "google.golang.org/grpc"
)

// alphaVantageResponse is the structure to parse Alpha Vantage's GLOBAL_QUOTE response.
type alphaVantageResponse struct {
    GlobalQuote struct {
        Symbol        string `json:"01. symbol"`
        Open          string `json:"02. open"`
        High          string `json:"03. high"`
        Low           string `json:"04. low"`
        Price         string `json:"05. price"`
        Volume        string `json:"06. volume"`
        LatestTrading string `json:"07. latest trading day"`
        PreviousClose string `json:"08. previous close"`
        Change        string `json:"09. change"`
        ChangePercent string `json:"10. change percent"`
    } `json:"Global Quote"`
}

// FetchQuote gets a real quote from Alpha Vantage for the given symbol.
func FetchQuote(symbol string) (string, float64, string, error) {
    // 1) Get your API key from environment or set it here
    //    It's best to store it in an env variable ALPHA_VANTAGE_KEY
    apiKey := os.Getenv("ALPHA_VANTAGE_KEY")
    if apiKey == "" {
        // For demonstration, you could hardcode your key or handle error
        // apiKey = "YOUR_ALPHAVANTAGE_API_KEY"
        return symbol, 0, "", fmt.Errorf("ALPHA_VANTAGE_KEY not set in environment")
    }

    // 2) Construct the Alpha Vantage URL
    url := fmt.Sprintf(
        "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
        symbol,
        apiKey,
    )

    // 3) Make the HTTP request
    resp, err := http.Get(url)
    if err != nil {
        return symbol, 0, "", fmt.Errorf("failed to fetch quote from Alpha Vantage: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return symbol, 0, "", fmt.Errorf("non-200 response from Alpha Vantage: %d", resp.StatusCode)
    }

    // 4) Parse the JSON
    var avResp alphaVantageResponse
    if err := json.NewDecoder(resp.Body).Decode(&avResp); err != nil {
        return symbol, 0, "", fmt.Errorf("failed to parse Alpha Vantage JSON: %v", err)
    }

    // 5) Extract price
    priceStr := avResp.GlobalQuote.Price
    if priceStr == "" {
        return symbol, 0, "", fmt.Errorf("no price returned for symbol=%s", symbol)
    }

    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        return symbol, 0, "", fmt.Errorf("failed to parse price string: %v", err)
    }

    // 6) Optionally notify if above threshold
    fmt.Printf("DEBUG: Fetched real price=%.2f for symbol=%s\n", price, symbol)

    if price > 200 {
        notifyUserMarketData("user123",
            fmt.Sprintf("Price of %s is now %.2f, above your threshold!", symbol, price),
            "PUSH",
        )
    }

    timestamp := time.Now().Format(time.RFC3339)
    return symbol, price, timestamp, nil
}

// notifyUserMarketData calls the Notification Service to send a message.
func notifyUserMarketData(userID, message, channel string) {
    // Dial the Notification Service (assuming it runs on port 50056)
    conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
    if err != nil {
        log.Printf("Error dialing Notification Service: %v\n", err)
        return
    }
    defer conn.Close()

    notifClient := notificationpb.NewNotificationServiceClient(conn)

    // Build the gRPC request
    req := &notificationpb.SendNotificationRequest{
        UserId:  userID,
        Message: message,
        Channel: channel, // e.g. "EMAIL", "SMS", "PUSH"
    }

    // Call SendNotification
    _, err = notifClient.SendNotification(context.Background(), req)
    if err != nil {
        log.Printf("Error sending market data notification: %v\n", err)
        return
    }

    log.Printf("Market Data notification sent to user=%s, channel=%s, message=%s\n",
        userID, channel, message)
}
