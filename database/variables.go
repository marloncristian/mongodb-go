package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)
