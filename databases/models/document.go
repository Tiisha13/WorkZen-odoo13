package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DocumentCategory string

const (
	DocumentCategoryResume  DocumentCategory = "resume"
	DocumentCategoryIDProof DocumentCategory = "id_proof"
	DocumentCategoryPayslip DocumentCategory = "payslip"
	DocumentCategoryPolicy  DocumentCategory = "policy"
	DocumentCategoryReport  DocumentCategory = "report"
	DocumentCategoryOther   DocumentCategory = "other"
)

// Document represents an uploaded file or stored HR document in the system
type Document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FileName    string             `bson:"file_name" json:"file_name"`                         // Original file name
	FilePath    string             `bson:"file_path" json:"file_path"`                         // Local file path
	FileURL     string             `bson:"file_url" json:"file_url"`                           // Public access URL
	FileType    string             `bson:"file_type" json:"file_type"`                         // e.g. pdf, jpg, png, docx
	Category    DocumentCategory   `bson:"category" json:"category"`                           // resume | id_proof | payslip | policy | report | other
	UploadedBy  primitive.ObjectID `bson:"uploaded_by,omitempty" json:"uploaded_by,omitempty"` // User who uploaded the file
	Company     primitive.ObjectID `bson:"company,omitempty" json:"company,omitempty"`         // Company context
	EmployeeID  primitive.ObjectID `bson:"employee_id,omitempty" json:"employee_id,omitempty"` // Optional (if document belongs to an employee)
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	IsPrivate   bool               `bson:"is_private" json:"is_private"`         // If true, restricted access (e.g., payslip)
	Size        int64              `bson:"size,omitempty" json:"size,omitempty"` // File size in bytes

	TimeStamp
}
