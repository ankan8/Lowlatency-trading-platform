package models

type Wallet struct {
    UserID  string  `bson:"user_id"`
    Balance float64 `bson:"balance"`
}
