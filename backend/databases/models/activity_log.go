package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ActivityLog stores user activities for audit trail
type ActivityLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Company   primitive.ObjectID `bson:"company,omitempty" json:"company,omitempty"`
	Action    string             `bson:"action" json:"action"`                         // login | create_user | update_payroll | etc
	Module    string             `bson:"module" json:"module"`                         // auth | user | payroll | attendance | etc
	Resource  string             `bson:"resource,omitempty" json:"resource,omitempty"` // Resource ID affected
	IPAddress string             `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	Metadata  interface{}        `bson:"metadata,omitempty" json:"metadata,omitempty"`

	TimeStamp
}
