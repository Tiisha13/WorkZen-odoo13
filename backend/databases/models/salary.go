package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WageType string
type ComponentType string

const (
	WageTypeFixed    WageType = "fixed"
	WageTypeVariable WageType = "variable"

	ComponentTypePercentage ComponentType = "percentage"
	ComponentTypeFixed      ComponentType = "fixed"
)

// SalaryComponent represents a single component of salary structure
type SalaryComponent struct {
	Name   string        `bson:"name" json:"name"`     // Basic, HRA, etc.
	Type   ComponentType `bson:"type" json:"type"`     // percentage | fixed
	Value  float64       `bson:"value" json:"value"`   // percentage value or fixed amount
	Amount float64       `bson:"amount" json:"amount"` // calculated amount
}

// SalaryStructure defines the salary breakdown for an employee
type SalaryStructure struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmployeeID    primitive.ObjectID `bson:"employee_id" json:"employee_id"`
	Company       primitive.ObjectID `bson:"company" json:"company"`
	WageType      WageType           `bson:"wage_type" json:"wage_type"`           // fixed | variable
	MonthlyWage   float64            `bson:"monthly_wage" json:"monthly_wage"`     // Total monthly wage
	YearlyWage    float64            `bson:"yearly_wage" json:"yearly_wage"`       // Total yearly wage (MonthlyWage * 12)
	Currency      string             `bson:"currency" json:"currency"`             // INR, USD, etc.
	EffectiveFrom string             `bson:"effective_from" json:"effective_from"` // YYYY-MM-DD

	// Salary Components
	BasicSalary          SalaryComponent `bson:"basic_salary" json:"basic_salary"`
	HouseRentAllowance   SalaryComponent `bson:"house_rent_allowance" json:"house_rent_allowance"`
	StandardAllowance    SalaryComponent `bson:"standard_allowance" json:"standard_allowance"`
	PerformanceBonus     SalaryComponent `bson:"performance_bonus" json:"performance_bonus"`
	LeaveTravelAllowance SalaryComponent `bson:"leave_travel_allowance" json:"leave_travel_allowance"`
	FixedAllowance       SalaryComponent `bson:"fixed_allowance" json:"fixed_allowance"`

	// Computed Values
	TotalEarnings   float64 `bson:"total_earnings" json:"total_earnings"`     // Sum of all components
	TotalDeductions float64 `bson:"total_deductions" json:"total_deductions"` // PF + Tax
	NetPay          float64 `bson:"net_pay" json:"net_pay"`                   // TotalEarnings - TotalDeductions

	IsActive bool `bson:"is_active" json:"is_active"` // Only one active structure per employee

	TimeStamp
}

// PayrollConfiguration stores company-wide payroll computation parameters
type PayrollConfiguration struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Company           primitive.ObjectID `bson:"company" json:"company"`
	PFEmployeePercent float64            `bson:"pf_employee_percent" json:"pf_employee_percent"` // default 12%
	PFEmployerPercent float64            `bson:"pf_employer_percent" json:"pf_employer_percent"` // default 12%
	ProfessionalTax   float64            `bson:"professional_tax" json:"professional_tax"`       // default ₹200

	// Default Component Ratios (as percentage of wage)
	DefaultBasicPercent      float64 `bson:"default_basic_percent" json:"default_basic_percent"`           // 50%
	DefaultHRAPercent        float64 `bson:"default_hra_percent" json:"default_hra_percent"`               // 50% of Basic
	DefaultStandardAllowance float64 `bson:"default_standard_allowance" json:"default_standard_allowance"` // 16.67%
	DefaultPerformanceBonus  float64 `bson:"default_performance_bonus" json:"default_performance_bonus"`   // 8.33%
	DefaultLTA               float64 `bson:"default_lta" json:"default_lta"`                               // 8.33%

	Currency string `bson:"currency" json:"currency"` // INR, USD, etc.

	TimeStamp
}
