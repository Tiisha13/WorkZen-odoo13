package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Department represents an organizational department
type Department struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Company     primitive.ObjectID `bson:"company" json:"company"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	HeadID      primitive.ObjectID `bson:"head_id,omitempty" json:"head_id,omitempty"`

	TimeStamp
}
