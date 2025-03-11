package repository

import (
    "context"
    "fmt"

    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/notification-service/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// InsertNotification inserts a notification record into the "notifications" collection.
func InsertNotification(n *models.Notification) error {
    coll := config.DB.Collection("notifications")
    _, err := coll.InsertOne(context.Background(), n)
    return err
}

// Optional: If you want to fetch a user's email from a "users" collection, define a User struct:
type User struct {
    UserID string `bson:"_id"`
    Email  string `bson:"email"`
}

// GetUserEmail fetches the user's email from a "users" collection by userID (optional).
func GetUserEmail(userID string) (string, error) {
    coll := config.DB.Collection("users") // or wherever you store user data
    var user User
    err := coll.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return "", fmt.Errorf("no user found with user_id=%s", userID)
        }
        return "", fmt.Errorf("db error: %v", err)
    }
    return user.Email, nil
}
