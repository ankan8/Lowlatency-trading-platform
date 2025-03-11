package repository

import (
    "context"
    

    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/billing-service/models"
    "go.mongodb.org/mongo-driver/bson"
)

func InsertTransaction(tx *models.Transaction) error {
    coll := config.DB.Collection("transactions")
    _, err := coll.InsertOne(context.Background(), tx)
    return err
}

func GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
    coll := config.DB.Collection("transactions")
    var txs []models.Transaction
    cursor, err := coll.Find(context.Background(), bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var t models.Transaction
        if err := cursor.Decode(&t); err != nil {
            return nil, err
        }
        txs = append(txs, t)
    }
    return txs, nil
}
