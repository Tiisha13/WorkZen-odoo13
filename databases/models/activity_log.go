package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ActivityLog stores user activities for audit trail
type ActivityLog struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Action   string             `bson:"action" json:"action"`
	Module   string             `bson:"module" json:"module"`
	Metadata interface{}        `bson:"metadata,omitempty" json:"metadata,omitempty"`

	TimeStamp
}
