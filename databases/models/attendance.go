package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// AttendanceStatus represents employee presence status
type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "present"
	StatusOnLeave AttendanceStatus = "on_leave"
	StatusAbsent  AttendanceStatus = "absent"
)

// Attendance represents daily employee attendance log
type Attendance struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmployeeID primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	Date       string             `bson:"date" json:"date"` // YYYY-MM-DD
	CheckIn    string             `bson:"check_in,omitempty" json:"check_in,omitempty"`
	CheckOut   string             `bson:"check_out,omitempty" json:"check_out,omitempty"`
	Status     AttendanceStatus   `bson:"status" json:"status"` // present | on_leave | absent
	WorkHours  float64            `bson:"work_hours,omitempty" json:"work_hours,omitempty"`
	Remarks    string             `bson:"remarks,omitempty" json:"remarks,omitempty"`

	TimeStamp
}
