package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Document represents an uploaded file or stored HR document in the system
type Document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FileName    string             `bson:"file_name" json:"file_name"`                         // Original file name
	FileURL     string             `bson:"file_url" json:"file_url"`                           // Public or internal access URL
	FileType    string             `bson:"file_type" json:"file_type"`                         // e.g. pdf, jpg, png, docx
	Category    string             `bson:"category" json:"category"`                           // employee_doc | payslip | policy | report
	OwnerID     primitive.ObjectID `bson:"owner_id,omitempty" json:"owner_id,omitempty"`       // User who uploaded the file
	CompanyID   primitive.ObjectID `bson:"company_id,omitempty" json:"company_id,omitempty"`   // Company context
	EmployeeID  primitive.ObjectID `bson:"employee_id,omitempty" json:"employee_id,omitempty"` // Optional (if document belongs to an employee)
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	IsPrivate   bool               `bson:"is_private" json:"is_private"`         // If true, restricted access (e.g., payslip)
	Size        int64              `bson:"size,omitempty" json:"size,omitempty"` // File size in bytes

	TimeStamp
}
