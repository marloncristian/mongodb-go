package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ENV_CONNECTIONSTRING = "MONGODB_CONNECTIONSTRING"

func ConnectDefault() (*mongo.Client, error) {
	return Connect(os.Getenv(ENV_CONNECTIONSTRING))
}

func Connect(connectionString string) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*25)
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString).SetDirect(true)

	Client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to create client %v", err)
	}
	err = Client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize connection %v", err)
	}
	err = Client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to connect %v", err)
	}

	Database = Client.Database(os.Getenv("COSMOSDB_DATABASE"))
	return Client, nil
}
