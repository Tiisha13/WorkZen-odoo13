// Package models contains database model definitions for the accounts service.
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Role defines different user roles across the WorkZen HRMS system
type Role string
type UserStatus string

const (
	RoleSuperAdmin Role = "superadmin" // Platform-level access
	RoleAdmin      Role = "admin"      // Company-level admin
	RoleHR         Role = "hr"         // Human Resource Officer
	RolePayroll    Role = "payroll"    // Payroll Officer
	RoleEmployee   Role = "employee"   // Regular employee

	UserActive   UserStatus = "active"
	UserInactive UserStatus = "inactive"
)

// User represents a registered person in the HRMS system.
// Each user belongs to a company (except SuperAdmin).
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username         string             `bson:"username" json:"username"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"password"`
	FirstName        string             `bson:"first_name" json:"first_name"`
	LastName         string             `bson:"last_name" json:"last_name"`
	Role             Role               `bson:"role" json:"role"` // superadmin | admin | hr | payroll | employee
	IsSuperAdmin     bool               `bson:"is_super_admin,omitempty" json:"is_super_admin,omitempty"`
	Designation      string             `bson:"designation,omitempty" json:"designation,omitempty"`
	DepartmentID     primitive.ObjectID `bson:"department_id,omitempty" json:"department_id,omitempty"`
	ManagerID        primitive.ObjectID `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	EmployeeCode     string             `bson:"employee_code,omitempty" json:"employee_code,omitempty"`
	DateOfJoin       string             `bson:"date_of_join,omitempty" json:"date_of_join,omitempty"` // YYYY-MM-DD
	Status           UserStatus         `bson:"status" json:"status"`                                 // active | inactive
	Phone            string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Address          Address            `bson:"address,omitempty" json:"address,omitempty"`
	ProfilePic       string             `bson:"profile_pic,omitempty" json:"profile_pic,omitempty"`
	ResumeURL        string             `bson:"resume_url,omitempty" json:"resume_url,omitempty"`
	Company          primitive.ObjectID `bson:"company,omitempty" json:"company,omitempty"`
	BankDetails      *BankDetails       `bson:"bank_details,omitempty" json:"bank_details,omitempty"`
	LastLogin        primitive.DateTime `bson:"last_login,omitempty" json:"last_login,omitempty"`
	EmailVerified    bool               `bson:"email_verified" json:"email_verified"`
	TwoFactorEnabled bool               `bson:"two_factor_enabled" json:"two_factor_enabled"`

	TimeStamp
}

// BankDetails stores employee banking information for payroll
type BankDetails struct {
	AccountNumber string `bson:"account_number,omitempty" json:"account_number,omitempty"`
	BankName      string `bson:"bank_name,omitempty" json:"bank_name,omitempty"`
	IFSCCode      string `bson:"ifsc_code,omitempty" json:"ifsc_code,omitempty"`
	BranchName    string `bson:"branch_name,omitempty" json:"branch_name,omitempty"`
	PANNo         string `bson:"pan_no,omitempty" json:"pan_no,omitempty"`
	UANNo         string `bson:"uan_no,omitempty" json:"uan_no,omitempty"` // Universal Account Number for PF
}

// Address structure embedded in User
type Address struct {
	City    string `bson:"city,omitempty" json:"city,omitempty"`
	State   string `bson:"state,omitempty" json:"state,omitempty"`
	Country string `bson:"country,omitempty" json:"country,omitempty"`
}
