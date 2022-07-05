package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ENV_CONNECTIONSTRING = "MONGODB_CONNECTIONSTRING"
const ENV_DATABASENAME = "MONGODB_DATABASE"

func ConnectDefault() (*mongo.Client, error) {
	return Connect(os.Getenv(ENV_CONNECTIONSTRING), os.Getenv(ENV_DATABASENAME))
}

func Connect(connectionString string, databaseName string) (*mongo.Client, error) {

	if len(connectionString) == 0 {
		return errors.New("invalid connectionstring")
	}
	if len(databaseName) == 0 {
		return errors.New("invalid database name")
	}

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

	Database = Client.Database(databaseName)
	return Client, nil
}
