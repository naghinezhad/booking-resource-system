package repository

import (
	"context"

	"github.com/naghinezhad/BookingResourceSystem/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResourceRepository struct {
	collection *mongo.Collection
}

func NewResourceRepository(db *mongo.Database) *ResourceRepository {
	return &ResourceRepository{
		collection: db.Collection("resources"),
	}
}

func (r *ResourceRepository) GetByID(
	ctx context.Context,
	id string,
) (*model.Resource, error) {

	filter := bson.M{
		"id": id,
	}

	var resource model.Resource

	err := r.collection.FindOne(ctx, filter).Decode(&resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}
