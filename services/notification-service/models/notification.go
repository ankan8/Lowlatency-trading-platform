package models

type Notification struct {
  NotificationID string `bson:"notification_id"`
  UserID         string `bson:"user_id"`
  Message        string `bson:"message"`
  Channel        string `bson:"channel"`
  Timestamp      string `bson:"timestamp"`
  // e.g. "EMAIL", "SMS", "PUSH"
}
