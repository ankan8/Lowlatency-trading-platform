package models

type TradeRecord struct {
  TradeID   string  `bson:"trade_id"`
  Symbol    string  `bson:"symbol"`
  Quantity  float64 `bson:"quantity"`
  Price     float64 `bson:"price"`
  OrderType string  `bson:"order_type"`
  Timestamp string  `bson:"timestamp"`
  UserID    string  `bson:"user_id"`
}
