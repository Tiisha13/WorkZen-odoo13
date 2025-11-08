package helpers

import "go.mongodb.org/mongo-driver/bson/primitive"

// ObjectID is a helper function that converts a string to a primitive.ObjectID
func ObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objectID, nil
}
