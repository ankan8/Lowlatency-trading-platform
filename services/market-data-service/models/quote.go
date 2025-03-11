package models

type Quote struct {
  Symbol    string  `bson:"symbol"`
  Price     float64 `bson:"price"`
  Timestamp string  `bson:"timestamp"`
}
