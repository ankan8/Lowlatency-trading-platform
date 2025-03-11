package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"

    "github.com/joho/godotenv"

    pb "github.com/ankan8/swapsync/backend/services/market-data-service/proto"
    "github.com/ankan8/swapsync/backend/services/market-data-service/service"
    "google.golang.org/grpc"
)

// server implements the MarketDataServiceServer interface.
type server struct {
    pb.UnimplementedMarketDataServiceServer
}

// GetQuote calls service.FetchQuote to get the real price, then returns it.
func (s *server) GetQuote(ctx context.Context, req *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
    sym, price, ts, err := service.FetchQuote(req.GetSymbol())
    if err != nil {
        return nil, err
    }

    // Debug
    fmt.Printf("DEBUG: Market Data returning price=%.2f for symbol=%s\n", price, sym)

    return &pb.GetQuoteResponse{
        Symbol:    sym,
        Price:     price,
        Timestamp: ts,
    }, nil
}

func init() {
    // Optional: Load .env file so we can read ALPHA_VANTAGE_KEY, etc.
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, relying on system environment variables.")
    }

    // (Optional) Debug: print ALPHA_VANTAGE_KEY if needed
    if key := os.Getenv("ALPHA_VANTAGE_KEY"); key == "" {
        log.Println("ALPHA_VANTAGE_KEY not set; real quotes won't work!")
    } else {
        log.Println("ALPHA_VANTAGE_KEY is set (not printing for security).")
    }
}

func main() {
    // Listen on port 50054 (or whichever you prefer).
    lis, err := net.Listen("tcp", ":50054") // or any free port
    if err != nil {
        log.Fatalf("Failed to listen on port 50054: %v", err)
    }

    // Create a new gRPC server.
    grpcServer := grpc.NewServer()

    // Register the MarketDataServiceServer implementation.
    pb.RegisterMarketDataServiceServer(grpcServer, &server{})

    log.Printf("Market Data Service listening on %v", lis.Addr())
    // Start serving gRPC
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
