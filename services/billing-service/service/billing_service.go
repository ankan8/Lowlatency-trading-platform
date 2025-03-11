package service

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/ankan8/swapsync/backend/services/billing-service/models"
    "github.com/ankan8/swapsync/backend/services/billing-service/repository"
    "github.com/google/uuid"
    "github.com/razorpay/razorpay-go"

    notificationpb "github.com/ankan8/swapsync/backend/services/notification-service/proto"
    "google.golang.org/grpc"

    pb "github.com/ankan8/swapsync/backend/services/billing-service/proto"
    walletRepo "github.com/ankan8/swapsync/backend/services/billing-service/repository"
)

// BillingServiceServer is the struct that implements all Billing gRPC methods.
type BillingServiceServer struct {
    pb.UnimplementedBillingServiceServer
}

// CalculateCommission implements the gRPC method for calculating commission.
func (s *BillingServiceServer) CalculateCommission(ctx context.Context, req *pb.CalculateCommissionRequest) (*pb.CalculateCommissionResponse, error) {
    tradeAmount := req.GetTradeAmount()
    commission := calculateCommission(tradeAmount)
    return &pb.CalculateCommissionResponse{Commission: commission}, nil
}

// ProcessPayment implements the gRPC method for processing payment.
func (s *BillingServiceServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
    userID := req.GetUserId()
    amount := req.GetAmount()
    method := req.GetMethod()

    success, txID, err := processPayment(userID, amount, method)
    if err != nil {
        return nil, err
    }
    return &pb.ProcessPaymentResponse{
        Success:       success,
        TransactionId: txID,
    }, nil
}

// DepositFunds implements the gRPC method for depositing funds.
func (s *BillingServiceServer) DepositFunds(ctx context.Context, req *pb.DepositFundsRequest) (*pb.DepositFundsResponse, error) {
    userID := req.GetUserId()
    amount := req.GetAmount()

    if amount <= 0 {
        return &pb.DepositFundsResponse{Success: false}, fmt.Errorf("invalid deposit amount")
    }

    // Try to fetch wallet
    wallet, err := walletRepo.GetWallet(userID)
    if err != nil {
        // If no wallet, create one
        if err.Error() == fmt.Sprintf("no wallet found for user=%s", userID) {
            if createErr := walletRepo.CreateWallet(userID); createErr != nil {
                return nil, createErr
            }
            wallet = &models.Wallet{UserID: userID, Balance: 0}
        } else {
            return nil, err
        }
    }

    newBalance := wallet.Balance + amount
    if err := walletRepo.UpdateWalletBalance(userID, newBalance); err != nil {
        return &pb.DepositFundsResponse{Success: false}, err
    }

    return &pb.DepositFundsResponse{
        Success:    true,
        NewBalance: newBalance,
    }, nil
}

// WithdrawFunds implements the gRPC method for withdrawing funds.
func (s *BillingServiceServer) WithdrawFunds(ctx context.Context, req *pb.WithdrawFundsRequest) (*pb.WithdrawFundsResponse, error) {
    userID := req.GetUserId()
    amount := req.GetAmount()

    if amount <= 0 {
        return &pb.WithdrawFundsResponse{Success: false}, fmt.Errorf("invalid withdraw amount")
    }

    wallet, err := walletRepo.GetWallet(userID)
    if err != nil {
        return &pb.WithdrawFundsResponse{Success: false}, err
    }

    if wallet.Balance < amount {
        return &pb.WithdrawFundsResponse{Success: false}, fmt.Errorf("insufficient funds")
    }

    newBalance := wallet.Balance - amount
    if err := walletRepo.UpdateWalletBalance(userID, newBalance); err != nil {
        return &pb.WithdrawFundsResponse{Success: false}, err
    }

    return &pb.WithdrawFundsResponse{
        Success:     true,
        NewBalance:  newBalance,
    }, nil
}

// GetBalance implements the gRPC method for retrieving wallet balance.
func (s *BillingServiceServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
    userID := req.GetUserId()

    wallet, err := walletRepo.GetWallet(userID)
    if err != nil {
        return &pb.GetBalanceResponse{Success: false}, err
    }

    return &pb.GetBalanceResponse{
        Success: true,
        Balance: wallet.Balance,
    }, nil
}

// calculateCommission is an internal helper for the logic of calculating trade commission.
func calculateCommission(tradeAmount float64) float64 {
    fmt.Printf("Received trade amount: %.2f\n", tradeAmount)
    if tradeAmount <= 0 {
        fmt.Println("Warning: Invalid trade amount received.")
        return 0
    }
    commission := tradeAmount * 0.001
    fmt.Printf("Calculated Commission: %.2f\n", commission)
    return commission
}

// processPayment is an internal helper that creates a Razorpay order, stores a transaction, and notifies the user.
func processPayment(userID string, amount float64, method string) (bool, string, error) {
    if userID == "" || amount <= 0 || method == "" {
        return false, "", fmt.Errorf("invalid payment details: userID=%s, amount=%.2f, method=%s",
            userID, amount, method)
    }

    // Create a Razorpay client
    keyID := os.Getenv("RAZORPAY_KEY_ID")
    keySecret := os.Getenv("RAZORPAY_KEY_SECRET")
    fmt.Printf("DEBUG: Using Razorpay Key ID=%s, Secret length=%d\n", keyID, len(keySecret))

    client := razorpay.NewClient(keyID, keySecret)

    // Convert amount to paise if currency is INR
    razorAmount := int64(amount * 100)

    // Create a Razorpay order
    orderData := map[string]interface{}{
        "amount":          razorAmount,
        "currency":        "INR",
        "receipt":         fmt.Sprintf("tx_%s", userID),
        "payment_capture": 1,
    }
    order, err := client.Order.Create(orderData, nil)
    if err != nil {
        return false, "", fmt.Errorf("failed to create Razorpay order: %v", err)
    }
    razorOrderID := order["id"].(string)
    fmt.Printf("Razorpay Order Created: %s\n", razorOrderID)

    txID := uuid.NewString()
    success := true

    tx := &models.Transaction{
        TransactionID: txID,
        UserID:        userID,
        Amount:        amount,
        Method:        method,
        Timestamp:     time.Now().Format(time.RFC3339),
        Success:       success,
    }

    err = repository.InsertTransaction(tx)
    if err != nil {
        fmt.Printf("Error: Failed to record transaction: %v\n", err)
        return false, "", fmt.Errorf("failed to record transaction: %v", err)
    }
    fmt.Printf("Transaction Recorded Successfully: %+v\n", tx)

    notifyUserBilling(userID, amount, method)

    return success, txID, nil
}

// notifyUserBilling sends a notification to the user about a successful payment.
func notifyUserBilling(userID string, amount float64, method string) {
    conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
    if err != nil {
        log.Printf("Error dialing Notification Service: %v\n", err)
        return
    }
    defer conn.Close()

    notifClient := notificationpb.NewNotificationServiceClient(conn)

    message := fmt.Sprintf("Your payment of %.2f via %s was processed successfully!", amount, method)

    // Call SendNotification
    _, err = notifClient.SendNotification(context.Background(), &notificationpb.SendNotificationRequest{
        UserId:  userID,
        Message: message,
        Channel: "EMAIL",
    })
    if err != nil {
        log.Printf("Error sending billing notification: %v\n", err)
        return
    }
    log.Printf("Billing notification sent to user=%s, message=%s\n", userID, message)
}
