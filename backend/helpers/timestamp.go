package helpers

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SetCreatedTimestamp sets the creation timestamp and user
func SetCreatedTimestamp(userID primitive.ObjectID) (primitive.DateTime, primitive.ObjectID) {
	return primitive.NewDateTimeFromTime(time.Now()), userID
}

// SetUpdatedTimestamp sets the update timestamp and user
func SetUpdatedTimestamp(userID primitive.ObjectID) (primitive.DateTime, primitive.ObjectID) {
	return primitive.NewDateTimeFromTime(time.Now()), userID
}

// SetDeletedTimestamp sets the delete timestamp and user for soft delete
func SetDeletedTimestamp(userID primitive.ObjectID) (primitive.DateTime, primitive.ObjectID) {
	return primitive.NewDateTimeFromTime(time.Now()), userID
}
