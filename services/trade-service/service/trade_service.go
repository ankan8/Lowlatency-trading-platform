package service

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/ankan8/swapsync/backend/services/trade-service/models"
    "github.com/ankan8/swapsync/backend/services/trade-service/repository"

    // Portfolio proto
    pbPortfolio "github.com/ankan8/swapsync/backend/services/portfolio-service/proto"

    // Market Data proto
    pbMarketData "github.com/ankan8/swapsync/backend/services/market-data-service/proto"

    // Billing proto
    pbBilling "github.com/ankan8/swapsync/backend/services/billing-service/proto"

    // Notification proto
    notificationpb "github.com/ankan8/swapsync/backend/services/notification-service/proto"

    "github.com/google/uuid"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

// PlaceOrder inserts a trade record and then calls:
// 1) Market Data Service to get the current price
// 2) Billing Service to deduct the trade cost from the user's wallet (only for BUY) + calculate/charge commission
// 3) Portfolio Service to update holdings (negative quantity for SELL)
// 4) Notification Service to alert the user of a successful trade
func PlaceOrder(userID, symbol string, userSuppliedPrice, quantity float64, orderType string, token string) (string, error) {
    // 1) Fetch the current quote from Market Data Service
    realPrice, err := fetchCurrentPrice(symbol, token)
    fmt.Printf("DEBUG: Fetched realPrice=%.2f\n", realPrice)
    if err != nil {
        return "", fmt.Errorf("failed to fetch market price for %s: %v", symbol, err)
    }

    // Decide whether to use userSuppliedPrice or realPrice
    // If userSuppliedPrice > 0 => treat as limit price, else realPrice (market).
    finalPrice := realPrice
    if userSuppliedPrice > 0 {
        finalPrice = userSuppliedPrice
    }

    // 1.5) Calculate total cost for the trade (user paying from wallet).
    // For a BUY, cost = finalPrice * quantity. For SELL, skip wallet deduction.
    if orderType == "BUY" {
        tradeCost := finalPrice * quantity
        // Check user has enough funds, then withdraw
        if err := checkAndWithdrawTradeCost(userID, tradeCost, token); err != nil {
            return "", err // insufficient funds or billing error
        }
    }

    // 2) Insert the trade record into MongoDB
    trade := &models.TradeRecord{
        TradeID:   uuid.NewString(),
        UserID:    userID,
        Symbol:    symbol,
        Quantity:  quantity,
        Price:     finalPrice,
        OrderType: orderType,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    if err := repository.InsertTradeRecord(trade); err != nil {
        return "", err
    }
    fmt.Printf("Trade executed: %s %.2f shares of %s at %.2f\n", orderType, quantity, symbol, finalPrice)

    // 2.1) Call the Billing Service to handle commission (for both BUY and SELL)
    tradeAmount := finalPrice * quantity
    if err := callBillingService(userID, tradeAmount, token); err != nil {
        // The trade record is already inserted, so handle/log the error
        return trade.TradeID, fmt.Errorf("failed to charge commission: %v", err)
    }

    // 3) Determine how much to update holdings by (negative for SELL)
    updateQuantity := quantity
    if orderType == "SELL" {
        updateQuantity = -quantity
    }

    // 4) Call the Portfolio Service to update holdings
    err = updatePortfolioHoldings(userID, symbol, updateQuantity, finalPrice, token)
    if err != nil {
        return trade.TradeID, fmt.Errorf("failed to update portfolio holdings: %v", err)
    }

    // 5) Notify the user about the successful trade
    notifyUserTrade(userID, symbol, quantity, finalPrice, orderType)

    return trade.TradeID, nil
}

// GetTradeHistory returns all trades for a user.
func GetTradeHistory(userID string) ([]models.TradeRecord, error) {
    return repository.GetTradesByUserID(userID)
}

// fetchCurrentPrice dials the Market Data Service to get the real-time quote.
func fetchCurrentPrice(symbol, token string) (float64, error) {
    conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure())
    if err != nil {
        return 0, fmt.Errorf("failed to dial Market Data Service: %v", err)
    }
    defer conn.Close()

    mdClient := pbMarketData.NewMarketDataServiceClient(conn)

    ctx := context.Background()
    if token != "" {
        md := metadata.New(map[string]string{"authorization": token})
        ctx = metadata.NewOutgoingContext(ctx, md)
    }

    resp, err := mdClient.GetQuote(ctx, &pbMarketData.GetQuoteRequest{
        Symbol: symbol,
    })
    if err != nil {
        return 0, fmt.Errorf("GetQuote RPC failed: %v", err)
    }
    return resp.GetPrice(), nil
}

// updatePortfolioHoldings dials the Portfolio Service's UpdateHoldings RPC.
func updatePortfolioHoldings(userID, symbol string, quantity, price float64, token string) error {
    conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
    if err != nil {
        return fmt.Errorf("failed to dial portfolio service: %v", err)
    }
    defer conn.Close()

    portfolioClient := pbPortfolio.NewPortfolioServiceClient(conn)

    ctx := context.Background()
    if token != "" {
        md := metadata.New(map[string]string{"authorization": token})
        ctx = metadata.NewOutgoingContext(ctx, md)
    }

    res, err := portfolioClient.UpdateHoldings(ctx, &pbPortfolio.UpdateHoldingsRequest{
        UserId:   userID,
        Symbol:   symbol,
        Quantity: quantity,
        Price:    price,
    })
    if err != nil {
        return fmt.Errorf("UpdateHoldings RPC failed: %v", err)
    }
    if !res.Success {
        return fmt.Errorf("UpdateHoldings responded with success=false")
    }

    log.Printf("Holdings updated in Portfolio Service for user %s, symbol %s\n", userID, symbol)
    return nil
}

// callBillingService calculates commission for the tradeAmount, then processes the payment.
func callBillingService(userID string, tradeAmount float64, token string) error {
    // 1) Connect to the Billing Service (assuming localhost:50055)
    conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
    if err != nil {
        return fmt.Errorf("failed to dial billing service: %v", err)
    }
    defer conn.Close()

    billingClient := pbBilling.NewBillingServiceClient(conn)

    // 2) Set up a context (including JWT if needed)
    ctx := context.Background()
    if token != "" {
        md := metadata.New(map[string]string{"authorization": token})
        ctx = metadata.NewOutgoingContext(ctx, md)
    }

    // 3) First, calculate the commission
    commResp, err := billingClient.CalculateCommission(ctx, &pbBilling.CalculateCommissionRequest{
        TradeAmount: tradeAmount,
    })
    if err != nil {
        return fmt.Errorf("CalculateCommission RPC failed: %v", err)
    }
    commission := commResp.GetCommission()
    fmt.Printf("Calculated commission for tradeAmount=%.2f is %.2f\n", tradeAmount, commission)

    // 3.5) If commission < 1 INR, skip Razorpay (avoid "Order amount less than minimum allowed").
    if commission < 1.0 {
        log.Printf("Commission=%.2f < 1.0 INR, skipping external payment\n", commission)
        return nil
    }

    // 4) Then process the payment for the commission
    payResp, err := billingClient.ProcessPayment(ctx, &pbBilling.ProcessPaymentRequest{
        UserId: userID,
        Amount: commission,
        Method: "WALLET", // or "CREDIT_CARD", "UPI", etc.
    })
    if err != nil {
        return fmt.Errorf("ProcessPayment RPC failed: %v", err)
    }
    if !payResp.GetSuccess() {
        return fmt.Errorf("ProcessPayment responded with success=false")
    }

    fmt.Printf("Commission of %.2f charged. Transaction ID=%s\n", commission, payResp.GetTransactionId())
    return nil
}

// checkAndWithdrawTradeCost ensures user has enough wallet balance and withdraws the cost for a BUY order.
func checkAndWithdrawTradeCost(userID string, cost float64, token string) error {
    if cost <= 0 {
        return nil // no cost to deduct if cost is zero or negative
    }

    conn, err := grpc.Dial("localhost:50055", grpc.WithInsecure())
    if err != nil {
        return fmt.Errorf("failed to dial billing service: %v", err)
    }
    defer conn.Close()

    billingClient := pbBilling.NewBillingServiceClient(conn)

    ctx := context.Background()
    if token != "" {
        md := metadata.New(map[string]string{"authorization": token})
        ctx = metadata.NewOutgoingContext(ctx, md)
    }

    // 1) Check current balance
    balResp, err := billingClient.GetBalance(ctx, &pbBilling.GetBalanceRequest{UserId: userID})
    if err != nil {
        return fmt.Errorf("failed to get wallet balance: %v", err)
    }
    if !balResp.Success {
        return fmt.Errorf("GetBalance responded with success=false")
    }
    if balResp.Balance < cost {
        return fmt.Errorf("insufficient wallet funds: need %.2f, have %.2f", cost, balResp.Balance)
    }

    // 2) Withdraw cost
    wdrResp, err := billingClient.WithdrawFunds(ctx, &pbBilling.WithdrawFundsRequest{
        UserId: userID,
        Amount: cost,
    })
    if err != nil {
        return fmt.Errorf("failed to withdraw trade cost: %v", err)
    }
    if !wdrResp.Success {
        return fmt.Errorf("WithdrawFunds responded with success=false")
    }

    log.Printf("Deducted trade cost of %.2f from user %s. New balance=%.2f\n", cost, userID, wdrResp.NewBalance)
    return nil
}

// notifyUserTrade calls the Notification Service to alert the user about the executed trade.
func notifyUserTrade(userID, symbol string, quantity, finalPrice float64, orderType string) {
    // 1) Connect to the Notification Service (assuming it runs on localhost:50056)
    conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
    if err != nil {
        log.Printf("Error dialing Notification Service: %v\n", err)
        return
    }
    defer conn.Close()

    notifClient := notificationpb.NewNotificationServiceClient(conn)

    // 2) Construct a message
    message := fmt.Sprintf("Your %s order for %.2f shares of %s is executed at %.2f",
        orderType, quantity, symbol, finalPrice)

    // 3) Call SendNotification
    _, err = notifClient.SendNotification(context.Background(), &notificationpb.SendNotificationRequest{
        UserId:  userID,
        Message: message,
        Channel: "EMAIL", // or "SMS", "PUSH"
    })
    if err != nil {
        log.Printf("Error sending trade notification: %v\n", err)
        return
    }

    log.Printf("Trade notification sent to user=%s, message=%s\n", userID, message)
}
