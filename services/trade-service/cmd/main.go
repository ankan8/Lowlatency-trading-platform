package main

import (
    "context"
    "log"
    "net"

    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/internal/middleware"
    pb "github.com/ankan8/swapsync/backend/services/trade-service/proto"
    "github.com/ankan8/swapsync/backend/services/trade-service/service"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

// server implements the gRPC interface for the TradeService.
type server struct {
    pb.UnimplementedTradeServiceServer
}

// GlobalOrderBooks is a map of symbol -> OrderBook (in-memory).
// We'll create or retrieve an OrderBook for each symbol as needed.
var GlobalOrderBooks = map[string]*service.OrderBook{}

// PlaceOrder is called by clients to place a trade order.
func (s *server) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
    // 1) Extract the JWT token from incoming metadata (if present).
    var token string
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        arr := md["authorization"]
        if len(arr) > 0 {
            token = arr[0]
        }
    }

    // 2) Ensure there's an OrderBook for this symbol (if you plan to do in-memory matching).
    symbol := req.GetSymbol()
    if _, ok := GlobalOrderBooks[symbol]; !ok {
        GlobalOrderBooks[symbol] = service.NewOrderBook(symbol)
        log.Printf("Created a new OrderBook for symbol=%s\n", symbol)
    }

    // 3) Call PlaceOrder in the service layer. 
    //    (Currently, this references wallet checks, commission, etc.)
    orderID, err := service.PlaceOrder(
        req.GetUserId(),
        symbol,
        req.GetPrice(),    // Price first
        req.GetQuantity(), // Quantity second
        req.GetOrderType(),
        token,
    )
    if err != nil {
        return &pb.PlaceOrderResponse{Success: false}, err
    }
    return &pb.PlaceOrderResponse{Success: true, OrderId: orderID}, nil
}

// GetTradeHistory returns all past trades for a given user.
func (s *server) GetTradeHistory(ctx context.Context, req *pb.GetTradeHistoryRequest) (*pb.GetTradeHistoryResponse, error) {
    trades, err := service.GetTradeHistory(req.GetUserId())
    if err != nil {
        return nil, err
    }
    var tradeRecords []*pb.TradeRecord
    for _, t := range trades {
        tradeRecords = append(tradeRecords, &pb.TradeRecord{
            TradeId:   t.TradeID,
            Symbol:    t.Symbol,
            Quantity:  t.Quantity,
            Price:     t.Price,
            OrderType: t.OrderType,
            Timestamp: t.Timestamp,
        })
    }
    return &pb.GetTradeHistoryResponse{Trades: tradeRecords}, nil
}

func main() {
    // 1) Connect to MongoDB
    config.ConnectDB()

    // For demonstration, pre-create an OrderBook for AAPL
    GlobalOrderBooks["AAPL"] = service.NewOrderBook("AAPL")
    log.Println("Initialized OrderBook for AAPL")

    // 2) Listen on port 50053
    lis, err := net.Listen("tcp", ":50053")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // 3) Create a gRPC server with the JWT interceptor
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(middleware.UnaryJWTInterceptor),
    )

    // 4) Register the TradeService
    pb.RegisterTradeServiceServer(grpcServer, &server{})

    log.Printf("Trade Service gRPC server is listening on %v", lis.Addr())

    // 5) Serve
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
