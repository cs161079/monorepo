package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnection(ctx context.Context) (*mongo.Database, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading enviroment variables.[%s]", err.Error())
	}
	ip := "127.0.0.1"
	if os.Getenv("application.mongo.ip") != "" {
		ip = os.Getenv("application.mongo.ip")
	}

	port := "27017"
	if os.Getenv("application.mongo.port") != "" {
		port = os.Getenv("application.mongo.port")
	}
	// ====================== Connect to Mongo Database ======================
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", ip, port)).
		SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	// ======================================================================

	// ================ Create or Use exist Database and Create Collection ==============
	return client.Database("testdb"), nil
}
