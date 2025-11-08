package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Company represents an organization registered in the WorkZen HRMS system
type Company struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	Phone      string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Industry   string             `bson:"industry,omitempty" json:"industry,omitempty"`
	Website    string             `bson:"website,omitempty" json:"website,omitempty"`
	LogoURL    string             `bson:"logo_url,omitempty" json:"logo_url,omitempty"`
	Address    Address            `bson:"address,omitempty" json:"address,omitempty"`
	OwnerID    primitive.ObjectID `bson:"owner_id,omitempty" json:"owner_id,omitempty"`
	ApprovedBy primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	IsApproved bool               `bson:"is_approved" json:"is_approved"`
	IsActive   bool               `bson:"is_active" json:"is_active"`

	TimeStamp
}
