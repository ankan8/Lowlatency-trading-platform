package main

import (
    
    "fmt"
    "log"
    "net"
    "os"

    "github.com/joho/godotenv"

    "github.com/ankan8/swapsync/backend/internal/config"
    pb "github.com/ankan8/swapsync/backend/services/billing-service/proto"
    "github.com/ankan8/swapsync/backend/services/billing-service/service"
    "google.golang.org/grpc"
)

func main() {
    // 1) Optionally load .env for Razorpay keys
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Warning: .env file not found, using system environment variables")
    }

    // Debugging
    fmt.Printf("DEBUG: RAZORPAY_KEY_ID=%s, SECRET length=%d\n",
        os.Getenv("RAZORPAY_KEY_ID"), len(os.Getenv("RAZORPAY_KEY_SECRET")))

    // 2) Connect to MongoDB if storing transactions/wallet data
    config.ConnectDB()

    // 3) Listen on port 50055
    lis, err := net.Listen("tcp", ":50055")
    if err != nil {
        log.Fatalf("Failed to listen on port 50055: %v", err)
    }

    // 4) Create gRPC server
    grpcServer := grpc.NewServer()

    // 5) Register our BillingServiceServer
    //    This is the struct that implements all the methods (CalculateCommission, ProcessPayment, DepositFunds, etc.)
    pb.RegisterBillingServiceServer(grpcServer, &service.BillingServiceServer{})

    log.Printf("Billing Service listening on %v", lis.Addr())

    // 6) Serve
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
