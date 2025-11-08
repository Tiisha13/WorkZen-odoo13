// Package models contains database model definitions for the accounts service.
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`

	CreatedAt primitive.DateTime  `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime  `bson:"updated_at" json:"updated_at"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
