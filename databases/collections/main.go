// Package collections defines MongoDB collection names used across the WorkZen HRMS system
package collections

const (
	// Core Collections
	Companies   = "companies"
	Users       = "users"
	Departments = "departments"

	// Attendance & Leave
	Attendances = "attendances"
	Leaves      = "leaves"

	// Payroll & Salary
	SalaryStructures      = "salary_structures"
	PayrollConfigurations = "payroll_configurations"
	Payruns               = "payruns"
	Payrolls              = "payrolls"

	// Documents
	Documents = "documents"

	// Audit & Logs
	ActivityLogs = "activity_logs"
)
