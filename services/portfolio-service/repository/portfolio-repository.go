package repository

import (
    "context"

    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/portfolio-service/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func GetPortfolioByUserID(userID string) (*models.Portfolio, error) {
    coll := config.DB.Collection("portfolios")
    var portfolio models.Portfolio
    err := coll.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&portfolio)
    if err != nil {
        return nil, err
    }
    return &portfolio, nil
}

func UpdateHoldings(userID, symbol string, quantity, price float64) error {
    coll := config.DB.Collection("portfolios")

    // Upsert logic: if portfolio doesn't exist, create it.
    filter := bson.M{"user_id": userID, "holdings.symbol": symbol}
    update := bson.M{
        "$set": bson.M{
            "user_id": userID,
        },
        "$inc": bson.M{
            "holdings.$.quantity": quantity,
        },
    }

    // Attempt to update existing holding
    res, err := coll.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }

    if res.MatchedCount == 0 {
        // Symbol not found in holdings, push a new entry
        newHolding := models.Holding{
            Symbol:       symbol,
            Quantity:     quantity,
            AveragePrice: price,
        }
        // Upsert the entire doc if user not found
        upsertFilter := bson.M{"user_id": userID}
        upsertUpdate := bson.M{
            "$setOnInsert": bson.M{"user_id": userID},
            "$push":        bson.M{"holdings": newHolding},
        }

        // Use UpdateOptions to set upsert=true
        upsertOpts := options.Update().SetUpsert(true)

        _, err = coll.UpdateOne(context.Background(), upsertFilter, upsertUpdate, upsertOpts)
        if err != nil {
            return err
        }
    }

    return nil
}
