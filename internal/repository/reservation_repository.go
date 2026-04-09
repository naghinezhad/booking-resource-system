package repository

import (
	"context"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationRepository struct {
	collection *mongo.Collection
}

func NewReservationRepository(db *mongo.Database) *ReservationRepository {
	return &ReservationRepository{
		collection: db.Collection("reservations"),
	}
}

func (r *ReservationRepository) Create(
	ctx context.Context,
	reservation *model.Reservation,
) error {

	_, err := r.collection.InsertOne(ctx, reservation)
	return err
}

func (r *ReservationRepository) CheckAvailability(
	ctx context.Context,
	resourceID primitive.ObjectID,
	start time.Time,
	end time.Time,
) (bool, error) {

	filter := bson.M{
		"resource_id": resourceID,
		"start_time": bson.M{
			"$lt": end,
		},
		"end_time": bson.M{
			"$gt": start,
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *ReservationRepository) GetByResourceID(
	ctx context.Context,
	resourceID primitive.ObjectID,
) ([]model.Reservation, error) {

	filter := bson.M{
		"resource_id": resourceID,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var reservations []model.Reservation

	for cursor.Next(ctx) {
		var res model.Reservation
		if err := cursor.Decode(&res); err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}

	return reservations, nil
}
