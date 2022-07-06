package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueryBuilder struct {
	Options QueryBuilderOptions
}

type QueryBuilderOptions struct {
	keyProperty string
}

func (qb QueryBuilder) Equals(property string, value interface{}) interface{} {
	return bson.M{property: value}
}

func (qb QueryBuilder) EqualsID(value primitive.ObjectID) interface{} {
	return qb.Equals(qb.Options.keyProperty, value)
}

func (qb QueryBuilder) EqualsHexID(value string) interface{} {
	objID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		return nil
	}
	return qb.Equals(qb.Options.keyProperty, objID)
}

func (qb QueryBuilder) Match(value map[string]interface{}) interface{} {
	var b []bson.M
	for k, v := range value {
		b = append(b, bson.M{k: v})
	}
	return bson.M{
		"$match": b,
	}
}

func (qb QueryBuilder) MatchSingle(property string, value interface{}) interface{} {
	return bson.M{
		"$match": bson.M{property: value},
	}
}

func (qb QueryBuilder) Sort(value map[string]int) interface{} {
	var b []bson.M
	for k, v := range value {
		b = append(b, bson.M{k: v})
	}
	return bson.M{
		"$sort": b,
	}
}

func (qb QueryBuilder) SortSingle(property string, value int) interface{} {
	return bson.M{
		"$sort": bson.M{property: value},
	}
}

func (qb QueryBuilder) Limit(value int) interface{} {
	return bson.M{
		"$limit": value,
	}
}

func (qb QueryBuilder) Skip(value int) interface{} {
	return bson.M{
		"$skip": value,
	}
}

func (qb QueryBuilder) Pipeline(args ...interface{}) interface{} {
	var a bson.A
	for _, arg := range args {
		a = append(a, arg)
	}
	return a
}
