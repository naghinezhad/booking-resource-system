package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ResourceID primitive.ObjectID `bson:"resource_id"`
	StartTime  time.Time          `bson:"start_time"`
	EndTime    time.Time          `bson:"end_time"`
	CreatedAt  time.Time          `bson:"created_at"`
}
