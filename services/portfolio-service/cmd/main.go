package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"

    pb "github.com/ankan8/swapsync/backend/services/portfolio-service/proto"
    "github.com/ankan8/swapsync/backend/services/portfolio-service/service"
    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/internal/middleware"
)

type server struct {
    pb.UnimplementedPortfolioServiceServer
}

func (s *server) GetPortfolio(ctx context.Context, req *pb.GetPortfolioRequest) (*pb.GetPortfolioResponse, error) {
    portfolio, err := service.GetPortfolio(req.GetUserId())
    if err != nil {
        return nil, err
    }
    // Convert portfolio holdings to proto holdings
    var holdings []*pb.Holding
    for _, h := range portfolio.Holdings {
        holdings = append(holdings, &pb.Holding{
            Symbol:       h.Symbol,
            Quantity:     h.Quantity,
            AveragePrice: h.AveragePrice,
        })
    }
    return &pb.GetPortfolioResponse{Holdings: holdings}, nil
}

func (s *server) UpdateHoldings(ctx context.Context, req *pb.UpdateHoldingsRequest) (*pb.UpdateHoldingsResponse, error) {
    err := service.UpdateHoldings(req.GetUserId(), req.GetSymbol(), req.GetQuantity(), req.GetPrice())
    if err != nil {
        return &pb.UpdateHoldingsResponse{Success: false}, err
    }
    return &pb.UpdateHoldingsResponse{Success: true}, nil
}

func main() {
    // Connect to MongoDB
    config.ConnectDB()

    lis, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    // Use the same JWT interceptor if you want to protect these endpoints
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(middleware.UnaryJWTInterceptor),
    )

    pb.RegisterPortfolioServiceServer(grpcServer, &server{})

    log.Printf("Portfolio Service gRPC server is listening on %v", lis.Addr())
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
