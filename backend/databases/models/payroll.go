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
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmployeeID primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	Company    primitive.ObjectID `bson:"company" json:"company"`
	PayrunID   primitive.ObjectID `bson:"payrun_id" json:"payrun_id"`
	Month      string             `bson:"month" json:"month"` // YYYY-MM

	// Salary Breakdown
	BasicSalary          float64 `bson:"basic_salary" json:"basic_salary"`
	HouseRentAllowance   float64 `bson:"house_rent_allowance" json:"house_rent_allowance"`
	StandardAllowance    float64 `bson:"standard_allowance" json:"standard_allowance"`
	PerformanceBonus     float64 `bson:"performance_bonus" json:"performance_bonus"`
	LeaveTravelAllowance float64 `bson:"leave_travel_allowance" json:"leave_travel_allowance"`
	FixedAllowance       float64 `bson:"fixed_allowance" json:"fixed_allowance"`

	// Totals
	GrossSalary     float64 `bson:"gross_salary" json:"gross_salary"`
	TotalDeductions float64 `bson:"total_deductions" json:"total_deductions"`
	NetPay          float64 `bson:"net_pay" json:"net_pay"`

	// Deductions
	PFEmployee      float64 `bson:"pf_employee" json:"pf_employee"`
	PFEmployer      float64 `bson:"pf_employer" json:"pf_employer"`
	ProfessionalTax float64 `bson:"professional_tax" json:"professional_tax"`

	// Attendance Data
	WorkingDays int `bson:"working_days" json:"working_days"`
	PresentDays int `bson:"present_days" json:"present_days"`
	LeaveDays   int `bson:"leave_days" json:"leave_days"`
	AbsentDays  int `bson:"absent_days" json:"absent_days"`

	// Warnings
	HasBankAccount bool `bson:"has_bank_account" json:"has_bank_account"`
	HasManager     bool `bson:"has_manager" json:"has_manager"`

	GeneratedBy primitive.ObjectID `bson:"generated_by" json:"generated_by"`
	GeneratedAt string             `bson:"generated_at" json:"generated_at"`
	Status      PayrollStatus      `bson:"status" json:"status"` // pending | processed | paid
	PaidAt      string             `bson:"paid_at,omitempty" json:"paid_at,omitempty"`
	PayslipURL  string             `bson:"payslip_url,omitempty" json:"payslip_url,omitempty"`

	TimeStamp
}
