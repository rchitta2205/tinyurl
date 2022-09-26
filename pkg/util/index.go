package util

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndex(ctx context.Context, collection *mongo.Collection, fieldName string) error {
	index := mongo.IndexModel{
		Keys:    bson.M{fieldName: 1},
		Options: options.Index().SetUnique(true),
	}

	// creating index for the input collection on input field to speed up search
	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
