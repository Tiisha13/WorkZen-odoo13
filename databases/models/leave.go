package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type LeaveStatus string
type LeaveType string

const (
	LeavePending  LeaveStatus = "pending"
	LeaveApproved LeaveStatus = "approved"
	LeaveRejected LeaveStatus = "rejected"

	LeaveSick     LeaveType = "sick"
	LeaveCasual   LeaveType = "casual"
	LeaveVacation LeaveType = "vacation"
)

// Leave represents a leave application request
type Leave struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmployeeID primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	Company    primitive.ObjectID `bson:"company" json:"company"`
	LeaveType  LeaveType          `bson:"leave_type" json:"leave_type"` // sick | casual | vacation
	Reason     string             `bson:"reason" json:"reason"`
	StartDate  string             `bson:"start_date" json:"start_date"` // YYYY-MM-DD
	EndDate    string             `bson:"end_date" json:"end_date"`     // YYYY-MM-DD
	Days       int                `bson:"days" json:"days"`
	Status     LeaveStatus        `bson:"status" json:"status"` // pending | approved | rejected
	ApprovedBy primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	RejectedBy primitive.ObjectID `bson:"rejected_by,omitempty" json:"rejected_by,omitempty"`
	ReviewedAt string             `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"` // YYYY-MM-DD HH:MM:SS

	TimeStamp
}
