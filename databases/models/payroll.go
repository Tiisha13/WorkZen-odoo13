package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PayrollStatus string

const (
	PayrollPending   PayrollStatus = "pending"
	PayrollProcessed PayrollStatus = "processed"
	PayrollPaid      PayrollStatus = "paid"
)

// Payroll represents monthly salary details of an employee
type Payroll struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmployeeID  primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	PayrunID    primitive.ObjectID `bson:"payrun_id" json:"payrun_id"`
	Month       string             `bson:"month" json:"month"` // YYYY-MM
	BasicSalary float64            `bson:"basic_salary" json:"basic_salary"`
	Allowances  Allowances         `bson:"allowances" json:"allowances"`
	Deductions  Deductions         `bson:"deductions" json:"deductions"`
	GrossSalary float64            `bson:"gross_salary" json:"gross_salary"`
	NetPay      float64            `bson:"net_pay" json:"net_pay"`
	GeneratedBy primitive.ObjectID `bson:"generated_by" json:"generated_by"`
	GeneratedAt string             `bson:"generated_at" json:"generated_at"`
	Status      PayrollStatus      `bson:"status" json:"status"` // pending | processed | paid
	PayslipURL  string             `bson:"payslip_url" json:"payslip_url"`

	TimeStamp
}

// Allowances defines additional salary components
type Allowances struct {
	HRA    float64 `bson:"hra" json:"hra"`
	Travel float64 `bson:"travel" json:"travel"`
	Other  float64 `bson:"other,omitempty" json:"other,omitempty"`
}

// Deductions defines salary deductions
type Deductions struct {
	PF              float64 `bson:"pf" json:"pf"`
	ProfessionalTax float64 `bson:"professional_tax" json:"professional_tax"`
	LeaveDeduction  float64 `bson:"leave_deduction" json:"leave_deduction"`
	Other           float64 `bson:"other,omitempty" json:"other,omitempty"`
}
