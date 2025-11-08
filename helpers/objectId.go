package helpers

import (
	"api.workzen.odoo/encryptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectID is a helper function that converts a string to a primitive.ObjectID
func ObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objectID, nil
}

// DecryptObjectID decrypts an encrypted ID and converts it to primitive.ObjectID
func DecryptObjectID(encryptedID string) (primitive.ObjectID, error) {
	// Decrypt the ID
	decryptedID, err := encryptions.DecryptID(encryptedID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// Convert to ObjectID
	objectID, err := primitive.ObjectIDFromHex(decryptedID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return objectID, nil
}
