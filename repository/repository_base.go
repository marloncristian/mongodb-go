package repository

import (
	"context"
	"errors"
	"reflect"

	"github.com/marloncristian/mongodb-go-helper/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryBase struct {
	collection  *mongo.Collection
	keyProperty string
}

func (base RepositoryBase) fillSlice(slice interface{}, cursor *mongo.Cursor) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return errors.New("parameter slice must be a pointer")
	}
	for cursor.Next(context.Background()) {
		spt := reflect.ValueOf(slice)
		svl := spt.Elem()
		sl := reflect.Indirect(spt)
		tT := sl.Type().Elem()
		ptr := reflect.New(tT).Interface()

		err := cursor.Decode(ptr)
		if err != nil {
			return err
		}
		s := reflect.ValueOf(ptr).Elem()
		svl.Set(reflect.Append(svl, s))
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

// Aggregate executes a aggregated command in the database
func (base RepositoryBase) Aggregate(pipeline interface{}, slice interface{}) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return errors.New("parameter slice must be a pointer")
	}

	cur, err := base.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())
	if err := base.fillSlice(slice, cur); err != nil {
		return err
	}

	return nil
}

// query retrieves documents by query or all
func (base RepositoryBase) Find(query interface{}, slice interface{}) error {
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		return errors.New("parameter slice must be a pointer")
	}

	cur, err := base.collection.Find(context.Background(), query)
	if err != nil {
		return err
	}

	defer cur.Close(context.Background())
	if err := base.fillSlice(slice, cur); err != nil {
		return err
	}

	return nil
}

// Count returns a count of all documents in repository
func (base RepositoryBase) Count(query interface{}) (int64, error) {
	cnt, err := base.collection.CountDocuments(context.Background(), query)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

//InsertOne : inserts a new object in repository
func (base RepositoryBase) InsertOne(value interface{}) (primitive.ObjectID, error) {
	res, err := base.collection.InsertOne(context.Background(), value)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

// UpdateOne : updates an document
func (base RepositoryBase) UpdateOne(query interface{}, update interface{}, result interface{}) error {
	doc := base.collection.FindOneAndUpdate(context.Background(), query, update)
	if doc.Err() != nil {
		return doc.Err()
	}
	if result == nil {
		return nil
	}
	if err := doc.Decode(result); err != nil {
		return err
	}
	return nil
}

// ReplaceOne replace an entire document
func (base RepositoryBase) ReplaceOne(query interface{}, entity interface{}) (err error) {
	_, err = base.collection.ReplaceOne(context.Background(), query, entity)
	return
}

// DeleteOne removes an elemento from database
func (base RepositoryBase) DeleteOne(id primitive.ObjectID) error {
	_, err := base.collection.DeleteOne(context.Background(), bson.M{base.keyProperty: id})
	if err != nil {
		return err
	}
	return nil
}

// NewRepositoryBase creates a new service base
func NewRepositoryBase(collectionName string) RepositoryBase {
	return RepositoryBase{
		collection: database.Database.Collection(collectionName),
	}
}
