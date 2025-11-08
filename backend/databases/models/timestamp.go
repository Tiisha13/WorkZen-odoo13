// Package models contains database model definitions for the accounts service.
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TimeStamp struct {
	CreatedAt primitive.DateTime  `bson:"created_at" json:"created_at"`
	CreatedBy primitive.ObjectID  `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedAt primitive.DateTime  `bson:"updated_at" json:"updated_at"`
	UpdatedBy primitive.ObjectID  `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	DeletedAt *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	DeletedBy *primitive.ObjectID `bson:"deleted_by,omitempty" json:"deleted_by,omitempty"`
	IsDeleted bool                `bson:"is_deleted" json:"is_deleted"`
}
