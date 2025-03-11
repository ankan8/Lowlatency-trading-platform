package repository

import (
    "context"
    "fmt"

    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/billing-service/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// GetWallet fetches the wallet document for a user.
func GetWallet(userID string) (*models.Wallet, error) {
    coll := config.DB.Collection("wallets")
    var w models.Wallet
    err := coll.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&w)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, fmt.Errorf("no wallet found for user=%s", userID)
        }
        return nil, err
    }
    return &w, nil
}

// CreateWallet initializes a new wallet with balance=0.
func CreateWallet(userID string) error {
    coll := config.DB.Collection("wallets")
    _, err := coll.InsertOne(context.Background(), models.Wallet{
        UserID:  userID,
        Balance: 0,
    })
    return err
}

// UpdateWalletBalance sets the wallet's balance to newBalance.
func UpdateWalletBalance(userID string, newBalance float64) error {
    coll := config.DB.Collection("wallets")
    _, err := coll.UpdateOne(
        context.Background(),
        bson.M{"user_id": userID},
        bson.M{"$set": bson.M{"balance": newBalance}},
    )
    return err
}
