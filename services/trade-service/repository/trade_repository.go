package repository

import (
    "context"
    

    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/trade-service/models"
    
    "go.mongodb.org/mongo-driver/bson"
)

func InsertTradeRecord(trade *models.TradeRecord) error {
    coll := config.DB.Collection("trades")
    _, err := coll.InsertOne(context.Background(), trade)
    return err
}

func GetTradesByUserID(userID string) ([]models.TradeRecord, error) {
    coll := config.DB.Collection("trades")
    var trades []models.TradeRecord
    cursor, err := coll.Find(context.Background(), bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var record models.TradeRecord
        if err := cursor.Decode(&record); err != nil {
            return nil, err
        }
        trades = append(trades, record)
    }
    return trades, nil
}
