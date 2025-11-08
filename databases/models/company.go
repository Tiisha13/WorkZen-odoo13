package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Company represents an organization registered in the WorkZen HRMS system
type Company struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`                       // Company name
	Email    string             `bson:"email,omitempty" json:"email,omitempty"` // Optional company contact email
	Phone    string             `bson:"phone,omitempty" json:"phone,omitempty"` // Optional contact number
	LogoURL  string             `bson:"logo_url,omitempty" json:"logo_url,omitempty"`
	Website  string             `bson:"website,omitempty" json:"website,omitempty"`
	Industry string             `bson:"industry,omitempty" json:"industry,omitempty"`
	Address  Address            `bson:"address,omitempty" json:"address,omitempty"`
	OwnerID  primitive.ObjectID `bson:"owner_id,omitempty" json:"owner_id,omitempty"` // User who created or owns the company
	IsActive bool               `bson:"is_active" json:"is_active"`                   // Active status

	TimeStamp
}
