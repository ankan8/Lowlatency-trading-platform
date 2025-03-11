package models

type Holding struct {
  Symbol       string  `bson:"symbol"`
  Quantity     float64 `bson:"quantity"`
  AveragePrice float64 `bson:"average_price"`
}

type Portfolio struct {
  UserID    string    `bson:"user_id"`
  Holdings  []Holding `bson:"holdings"`
}
