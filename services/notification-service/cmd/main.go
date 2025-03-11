package main

import (
    "context"
    "fmt"
    "log"
    "net"
	"os"
	

    "github.com/joho/godotenv"

    "github.com/ankan8/swapsync/backend/internal/config"
    pb "github.com/ankan8/swapsync/backend/services/notification-service/proto"
    "github.com/ankan8/swapsync/backend/services/notification-service/service"
    "google.golang.org/grpc"
)

type server struct {
    pb.UnimplementedNotificationServiceServer
}

// SendNotification calls our service.SendNotification logic.
func (s *server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
    userID := req.GetUserId()
    message := req.GetMessage()
    channel := req.GetChannel()

    success, notifID, err := service.SendNotification(userID, message, channel)
    if err != nil {
        return nil, err
    }
    return &pb.SendNotificationResponse{
        Success:         success,
        NotificationId:  notifID,
    }, nil
}

// StreamNotifications is an optional server-side streaming method for real-time notifications.
func (s *server) StreamNotifications(req *pb.StreamNotificationsRequest, stream pb.NotificationService_StreamNotificationsServer) error {
    // Implementation depends on your real-time strategy. 
    // You might watch a DB change stream or keep a channel of messages to push.
    return fmt.Errorf("StreamNotifications not implemented yet")
}

func main() {
    // Optionally
	//  load .env
	
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Warning: No .env file found, using system environment variables")
    }
	fmt.Printf("DEBUG: SENDGRID_API_KEY length=%d\n", len(os.Getenv("SENDGRID_API_KEY")))

    // Connect to DB if you want to store notifications
    config.ConnectDB()

    lis, err := net.Listen("tcp", ":50056") // e.g., port 50056
    if err != nil {
        log.Fatalf("Failed to listen on port 50056: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterNotificationServiceServer(grpcServer, &server{})

    log.Printf("Notification Service listening on %v", lis.Addr())
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
