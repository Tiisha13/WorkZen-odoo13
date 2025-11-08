// Package models contains database model definitions for the accounts service.
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeStamp struct {
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
	DeletedAt primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
