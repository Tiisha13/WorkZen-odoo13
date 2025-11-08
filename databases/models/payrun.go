package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PayrunStatus string

const (
	PayrunDraft      PayrunStatus = "draft"
	PayrunProcessing PayrunStatus = "processing"
	PayrunCompleted  PayrunStatus = "completed"
)

// Payrun represents a payroll batch for a specific period
type Payrun struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Month          string             `bson:"month" json:"month"`
	GeneratedBy    primitive.ObjectID `bson:"generated_by" json:"generated_by"`
	StartDate      string             `bson:"start_date" json:"start_date"`
	EndDate        string             `bson:"end_date" json:"end_date"`
	TotalEmployees int                `bson:"total_employees" json:"total_employees"`
	TotalAmount    float64            `bson:"total_amount" json:"total_amount"`
	Status         PayrunStatus       `bson:"status" json:"status"` // draft | processing | completed

	TimeStamp
}
