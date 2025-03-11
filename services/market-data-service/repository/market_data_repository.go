package repository

import (
    "context"
    "github.com/ankan8/swapsync/backend/internal/config"
    "github.com/ankan8/swapsync/backend/services/market-data-service/models"
)

func SaveQuote(q *models.Quote) error {
    coll := config.DB.Collection("quotes")
    _, err := coll.InsertOne(context.Background(), q)
    return err
}
