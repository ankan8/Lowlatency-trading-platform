package config

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/joho/godotenv"
)

var DB *mongo.Database

func ConnectDB() {
    
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found!!")
    }

    uri := os.Getenv("MONGO_URI")
    if uri == "" {
        log.Fatal("MONGO_URI not found in .env file")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatal("Error connecting to MongoDB:", err)
    }

    log.Println("Connected to MongoDB!")

    
    DB = client.Database("swapsync")
}
