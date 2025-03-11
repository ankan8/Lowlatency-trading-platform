package repository

import (
    "fmt"
    "context"
    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/auth-service/models"
)


func FindUserByEmail(email string) (*models.User, error) {
	collection := config.DB.Collection("users")
    var user models.User
    err := collection.FindOne(context.Background(), map[string]interface{}{"email": email}).Decode(&user)
    return &user, err
}
func InsertUser(user *models.User) error {
	collection := config.DB.Collection("users")
    result, err := collection.InsertOne(context.Background(), user)
    if err != nil {
        return err
    }
    fmt.Printf("Inserted user with ID: %v\n", result.InsertedID)
    return nil
}