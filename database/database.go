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

	conn := os.Getenv(ENV_CONNECTIONSTRING)
	if len(conn) == 0 {
		return nil, fmt.Errorf("environment variable %v not found", ENV_CONNECTIONSTRING)
	}

	dat := os.Getenv(ENV_DATABASENAME)
	if len(dat) == 0 {
		return nil, fmt.Errorf("environment variable %v not found", ENV_DATABASENAME)
	}

	return Connect(conn, dat)
}

func Connect(connectionString string, databaseName string) (*mongo.Client, error) {

	if len(connectionString) == 0 {
		return nil, errors.New("invalid connectionstring")
	}
	if len(databaseName) == 0 {
		return nil, errors.New("invalid database name")
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
