package models

type Transaction struct {
  TransactionID string  `bson:"transaction_id"`
  UserID        string  `bson:"user_id"`
  Amount        float64 `bson:"amount"`
  Method        string  `bson:"method"`
  Timestamp     string  `bson:"timestamp"`
  Success       bool    `bson:"success"`
}
