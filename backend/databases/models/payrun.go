package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PayrunStatus string

const (
	PayrunDraft     PayrunStatus = "draft"
	PayrunGenerated PayrunStatus = "generated"
	PayrunCompleted PayrunStatus = "completed"
)

// Payrun represents a payroll batch for a specific period
type Payrun struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Company        primitive.ObjectID `bson:"company" json:"company"`
	Month          string             `bson:"month" json:"month"` // YYYY-MM
	GeneratedBy    primitive.ObjectID `bson:"generated_by" json:"generated_by"`
	GeneratedAt    string             `bson:"generated_at" json:"generated_at"`
	StartDate      string             `bson:"start_date" json:"start_date"` // YYYY-MM-DD
	EndDate        string             `bson:"end_date" json:"end_date"`     // YYYY-MM-DD
	TotalEmployees int                `bson:"total_employees" json:"total_employees"`
	ProcessedCount int                `bson:"processed_count" json:"processed_count"`
	TotalPayroll   float64            `bson:"total_payroll" json:"total_payroll"` // Sum of all net pay
	Status         PayrunStatus       `bson:"status" json:"status"`               // draft | generated | completed

	// Warning Counts
	MissingBankCount    int `bson:"missing_bank_count" json:"missing_bank_count"`
	MissingManagerCount int `bson:"missing_manager_count" json:"missing_manager_count"`

	TimeStamp
}
