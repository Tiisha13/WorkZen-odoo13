package services

import (
	"fmt"

	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/encryptions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AttendanceResponse represents attendance data with encrypted IDs
type AttendanceResponse struct {
	ID         string                  `json:"id,omitempty"`
	EmployeeID string                  `json:"employee_id"`
	Company    string                  `json:"company"`
	Date       string                  `json:"date"`
	CheckIn    string                  `json:"check_in,omitempty"`
	CheckOut   string                  `json:"check_out,omitempty"`
	Status     models.AttendanceStatus `json:"status"`
	WorkHours  float64                 `json:"work_hours,omitempty"`
	Remarks    string                  `json:"remarks,omitempty"`
	CreatedAt  primitive.DateTime      `json:"created_at,omitempty"`
	UpdatedAt  primitive.DateTime      `json:"updated_at,omitempty"`
}

// LeaveResponse represents leave data with encrypted IDs
type LeaveResponse struct {
	ID         string             `json:"id,omitempty"`
	EmployeeID string             `json:"employee_id"`
	Company    string             `json:"company"`
	LeaveType  models.LeaveType   `json:"leave_type"`
	Reason     string             `json:"reason"`
	StartDate  string             `json:"start_date"`
	EndDate    string             `json:"end_date"`
	Days       int                `json:"days"`
	Status     models.LeaveStatus `json:"status"`
	ApprovedBy string             `json:"approved_by,omitempty"`
	RejectedBy string             `json:"rejected_by,omitempty"`
	ReviewedAt string             `json:"reviewed_at,omitempty"`
	User       *UserResponse      `json:"user,omitempty"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty"`
	UpdatedAt  primitive.DateTime `json:"updated_at,omitempty"`
}

// SalaryStructureResponse represents salary structure data with encrypted IDs
type SalaryStructureResponse struct {
	ID                   string                 `json:"id,omitempty"`
	EmployeeID           string                 `json:"employee_id"`
	Company              string                 `json:"company"`
	WageType             models.WageType        `json:"wage_type"`
	MonthlyWage          float64                `json:"monthly_wage"`
	YearlyWage           float64                `json:"yearly_wage"`
	Currency             string                 `json:"currency"`
	EffectiveFrom        string                 `json:"effective_from"`
	BasicSalary          models.SalaryComponent `json:"basic_salary"`
	HouseRentAllowance   models.SalaryComponent `json:"house_rent_allowance"`
	StandardAllowance    models.SalaryComponent `json:"standard_allowance"`
	PerformanceBonus     models.SalaryComponent `json:"performance_bonus"`
	LeaveTravelAllowance models.SalaryComponent `json:"leave_travel_allowance"`
	FixedAllowance       models.SalaryComponent `json:"fixed_allowance"`
	TotalEarnings        float64                `json:"total_earnings"`
	TotalDeductions      float64                `json:"total_deductions"`
	NetPay               float64                `json:"net_pay"`
	IsActive             bool                   `json:"is_active"`
	CreatedAt            primitive.DateTime     `json:"created_at,omitempty"`
	UpdatedAt            primitive.DateTime     `json:"updated_at,omitempty"`
}

// DocumentResponse represents document data with encrypted IDs
type DocumentResponse struct {
	ID          string                  `json:"id,omitempty"`
	FileName    string                  `json:"file_name"`
	FilePath    string                  `json:"file_path"`
	FileURL     string                  `json:"file_url"`
	FileType    string                  `json:"file_type"`
	Category    models.DocumentCategory `json:"category"`
	UploadedBy  string                  `json:"uploaded_by,omitempty"`
	Company     string                  `json:"company,omitempty"`
	EmployeeID  string                  `json:"employee_id,omitempty"`
	Description string                  `json:"description,omitempty"`
	IsPrivate   bool                    `json:"is_private"`
	Size        int64                   `json:"size,omitempty"`
	CreatedAt   primitive.DateTime      `json:"created_at,omitempty"`
	UpdatedAt   primitive.DateTime      `json:"updated_at,omitempty"`
}

// PayrollResponse represents payroll data with encrypted IDs
type PayrollResponse struct {
	ID                   string               `json:"id,omitempty"`
	EmployeeID           string               `json:"employee_id"`
	Company              string               `json:"company"`
	PayrunID             string               `json:"payrun_id"`
	Month                string               `json:"month"`
	BasicSalary          float64              `json:"basic_salary"`
	HouseRentAllowance   float64              `json:"house_rent_allowance"`
	StandardAllowance    float64              `json:"standard_allowance"`
	PerformanceBonus     float64              `json:"performance_bonus"`
	LeaveTravelAllowance float64              `json:"leave_travel_allowance"`
	FixedAllowance       float64              `json:"fixed_allowance"`
	GrossSalary          float64              `json:"gross_salary"`
	TotalDeductions      float64              `json:"total_deductions"`
	NetPay               float64              `json:"net_pay"`
	PFEmployee           float64              `json:"pf_employee"`
	PFEmployer           float64              `json:"pf_employer"`
	ProfessionalTax      float64              `json:"professional_tax"`
	WorkingDays          int                  `json:"working_days"`
	PresentDays          int                  `json:"present_days"`
	LeaveDays            int                  `json:"leave_days"`
	AbsentDays           int                  `json:"absent_days"`
	HasBankAccount       bool                 `json:"has_bank_account"`
	HasManager           bool                 `json:"has_manager"`
	GeneratedBy          string               `json:"generated_by"`
	GeneratedAt          string               `json:"generated_at"`
	Status               models.PayrollStatus `json:"status"`
	PaidAt               string               `json:"paid_at,omitempty"`
	PayslipURL           string               `json:"payslip_url,omitempty"`
	CreatedAt            primitive.DateTime   `json:"created_at,omitempty"`
	UpdatedAt            primitive.DateTime   `json:"updated_at,omitempty"`
}

// PayrunResponse represents payrun data with encrypted IDs
type PayrunResponse struct {
	ID                  string              `json:"id,omitempty"`
	Company             string              `json:"company"`
	Month               string              `json:"month"`
	GeneratedBy         string              `json:"generated_by"`
	GeneratedAt         string              `json:"generated_at"`
	StartDate           string              `json:"start_date"`
	EndDate             string              `json:"end_date"`
	TotalEmployees      int                 `json:"total_employees"`
	ProcessedCount      int                 `json:"processed_count"`
	TotalPayroll        float64             `json:"total_payroll"`
	Status              models.PayrunStatus `json:"status"`
	MissingBankCount    int                 `json:"missing_bank_count"`
	MissingManagerCount int                 `json:"missing_manager_count"`
	CreatedAt           primitive.DateTime  `json:"created_at,omitempty"`
	UpdatedAt           primitive.DateTime  `json:"updated_at,omitempty"`
}

// Converter functions

// ConvertAttendanceToResponse converts Attendance model to AttendanceResponse with encrypted IDs
func ConvertAttendanceToResponse(attendance *models.Attendance) (*AttendanceResponse, error) {
	if attendance == nil {
		return nil, nil
	}

	response := &AttendanceResponse{
		Date:      attendance.Date,
		CheckIn:   attendance.CheckIn,
		CheckOut:  attendance.CheckOut,
		Status:    attendance.Status,
		WorkHours: attendance.WorkHours,
		Remarks:   attendance.Remarks,
		CreatedAt: attendance.CreatedAt,
		UpdatedAt: attendance.UpdatedAt,
	}

	if !attendance.ID.IsZero() {
		encID, err := encryptions.EncryptID(attendance.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt attendance ID: %w", err)
		}
		response.ID = encID
	}

	if !attendance.EmployeeID.IsZero() {
		encID, err := encryptions.EncryptID(attendance.EmployeeID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt employee ID: %w", err)
		}
		response.EmployeeID = encID
	}

	if !attendance.Company.IsZero() {
		encID, err := encryptions.EncryptID(attendance.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encID
	}

	return response, nil
}

// ConvertLeaveToResponse converts Leave model to LeaveResponse with encrypted IDs
func ConvertLeaveToResponse(leave *models.Leave) (*LeaveResponse, error) {
	if leave == nil {
		return nil, nil
	}

	response := &LeaveResponse{
		LeaveType:  leave.LeaveType,
		Reason:     leave.Reason,
		StartDate:  leave.StartDate,
		EndDate:    leave.EndDate,
		Days:       leave.Days,
		Status:     leave.Status,
		ReviewedAt: leave.ReviewedAt,
		CreatedAt:  leave.CreatedAt,
		UpdatedAt:  leave.UpdatedAt,
	}

	if !leave.ID.IsZero() {
		encID, err := encryptions.EncryptID(leave.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt leave ID: %w", err)
		}
		response.ID = encID
	}

	if !leave.EmployeeID.IsZero() {
		encID, err := encryptions.EncryptID(leave.EmployeeID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt employee ID: %w", err)
		}
		response.EmployeeID = encID
	}

	if !leave.Company.IsZero() {
		encID, err := encryptions.EncryptID(leave.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encID
	}

	if !leave.ApprovedBy.IsZero() {
		encID, err := encryptions.EncryptID(leave.ApprovedBy.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt approved by ID: %w", err)
		}
		response.ApprovedBy = encID
	}

	if !leave.RejectedBy.IsZero() {
		encID, err := encryptions.EncryptID(leave.RejectedBy.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt rejected by ID: %w", err)
		}
		response.RejectedBy = encID
	}

	return response, nil
}

// ConvertLeaveToResponseWithUser converts LeaveWithUser to LeaveResponse with user data
func ConvertLeaveToResponseWithUser(leave *LeaveWithUser) (*LeaveResponse, error) {
	if leave == nil {
		return nil, nil
	}

	// First convert the base leave
	response, err := ConvertLeaveToResponse(&leave.Leave)
	if err != nil {
		return nil, err
	}

	// Add user data if available
	if leave.User != nil {
		userResp, err := ConvertUserToResponse(leave.User)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user: %w", err)
		}
		response.User = userResp
	}

	return response, nil
}

// ConvertSalaryStructureToResponse converts SalaryStructure model to SalaryStructureResponse with encrypted IDs
func ConvertSalaryStructureToResponse(salary *models.SalaryStructure) (*SalaryStructureResponse, error) {
	if salary == nil {
		return nil, nil
	}

	response := &SalaryStructureResponse{
		WageType:             salary.WageType,
		MonthlyWage:          salary.MonthlyWage,
		YearlyWage:           salary.YearlyWage,
		Currency:             salary.Currency,
		EffectiveFrom:        salary.EffectiveFrom,
		BasicSalary:          salary.BasicSalary,
		HouseRentAllowance:   salary.HouseRentAllowance,
		StandardAllowance:    salary.StandardAllowance,
		PerformanceBonus:     salary.PerformanceBonus,
		LeaveTravelAllowance: salary.LeaveTravelAllowance,
		FixedAllowance:       salary.FixedAllowance,
		TotalEarnings:        salary.TotalEarnings,
		TotalDeductions:      salary.TotalDeductions,
		NetPay:               salary.NetPay,
		IsActive:             salary.IsActive,
		CreatedAt:            salary.CreatedAt,
		UpdatedAt:            salary.UpdatedAt,
	}

	if !salary.ID.IsZero() {
		encID, err := encryptions.EncryptID(salary.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt salary ID: %w", err)
		}
		response.ID = encID
	}

	if !salary.EmployeeID.IsZero() {
		encID, err := encryptions.EncryptID(salary.EmployeeID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt employee ID: %w", err)
		}
		response.EmployeeID = encID
	}

	if !salary.Company.IsZero() {
		encID, err := encryptions.EncryptID(salary.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encID
	}

	return response, nil
}

// ConvertDocumentToResponse converts Document model to DocumentResponse with encrypted IDs
func ConvertDocumentToResponse(doc *models.Document) (*DocumentResponse, error) {
	if doc == nil {
		return nil, nil
	}

	response := &DocumentResponse{
		FileName:    doc.FileName,
		FilePath:    doc.FilePath,
		FileURL:     doc.FileURL,
		FileType:    doc.FileType,
		Category:    doc.Category,
		Description: doc.Description,
		IsPrivate:   doc.IsPrivate,
		Size:        doc.Size,
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}

	if !doc.ID.IsZero() {
		encID, err := encryptions.EncryptID(doc.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt document ID: %w", err)
		}
		response.ID = encID
	}

	if !doc.UploadedBy.IsZero() {
		encID, err := encryptions.EncryptID(doc.UploadedBy.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt uploaded by ID: %w", err)
		}
		response.UploadedBy = encID
	}

	if !doc.Company.IsZero() {
		encID, err := encryptions.EncryptID(doc.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encID
	}

	if !doc.EmployeeID.IsZero() {
		encID, err := encryptions.EncryptID(doc.EmployeeID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt employee ID: %w", err)
		}
		response.EmployeeID = encID
	}

	return response, nil
}

// ConvertPayrollToResponse converts Payroll model to PayrollResponse with encrypted IDs
func ConvertPayrollToResponse(payroll *models.Payroll) (*PayrollResponse, error) {
	if payroll == nil {
		return nil, nil
	}

	response := &PayrollResponse{
		Month:                payroll.Month,
		BasicSalary:          payroll.BasicSalary,
		HouseRentAllowance:   payroll.HouseRentAllowance,
		StandardAllowance:    payroll.StandardAllowance,
		PerformanceBonus:     payroll.PerformanceBonus,
		LeaveTravelAllowance: payroll.LeaveTravelAllowance,
		FixedAllowance:       payroll.FixedAllowance,
		GrossSalary:          payroll.GrossSalary,
		TotalDeductions:      payroll.TotalDeductions,
		NetPay:               payroll.NetPay,
		PFEmployee:           payroll.PFEmployee,
		PFEmployer:           payroll.PFEmployer,
		ProfessionalTax:      payroll.ProfessionalTax,
		WorkingDays:          payroll.WorkingDays,
		PresentDays:          payroll.PresentDays,
		LeaveDays:            payroll.LeaveDays,
		AbsentDays:           payroll.AbsentDays,
		HasBankAccount:       payroll.HasBankAccount,
		HasManager:           payroll.HasManager,
		GeneratedAt:          payroll.GeneratedAt,
		Status:               payroll.Status,
		PaidAt:               payroll.PaidAt,
		PayslipURL:           payroll.PayslipURL,
		CreatedAt:            payroll.CreatedAt,
		UpdatedAt:            payroll.UpdatedAt,
	}

	if !payroll.ID.IsZero() {
		encID, err := encryptions.EncryptID(payroll.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt payroll ID: %w", err)
		}
		response.ID = encID
	}

	if !payroll.EmployeeID.IsZero() {
		encID, err := encryptions.EncryptID(payroll.EmployeeID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt employee ID: %w", err)
		}
		response.EmployeeID = encID
	}

	if !payroll.Company.IsZero() {
		encID, err := encryptions.EncryptID(payroll.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encID
	}

	if !payroll.PayrunID.IsZero() {
		encID, err := encryptions.EncryptID(payroll.PayrunID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt payrun ID: %w", err)
		}
		response.PayrunID = encID
	}

	if !payroll.GeneratedBy.IsZero() {
		encID, err := encryptions.EncryptID(payroll.GeneratedBy.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt generated by ID: %w", err)
		}
		response.GeneratedBy = encID
	}

	return response, nil
}

// ConvertPayrunToResponse converts Payrun model to PayrunResponse with encrypted IDs
func ConvertPayrunToResponse(payrun *models.Payrun) (*PayrunResponse, error) {
	if payrun == nil {
		return nil, nil
	}

	response := &PayrunResponse{
		Month:               payrun.Month,
		GeneratedAt:         payrun.GeneratedAt,
		StartDate:           payrun.StartDate,
		EndDate:             payrun.EndDate,
		TotalEmployees:      payrun.TotalEmployees,
		ProcessedCount:      payrun.ProcessedCount,
		TotalPayroll:        payrun.TotalPayroll,
		Status:              payrun.Status,
		MissingBankCount:    payrun.MissingBankCount,
		MissingManagerCount: payrun.MissingManagerCount,
		CreatedAt:           payrun.CreatedAt,
		UpdatedAt:           payrun.UpdatedAt,
	}

	if !payrun.ID.IsZero() {
		encID, err := encryptions.EncryptID(payrun.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt payrun ID: %w", err)
		}
		response.ID = encID
	}

	if !payrun.Company.IsZero() {
		encID, err := encryptions.EncryptID(payrun.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encID
	}

	if !payrun.GeneratedBy.IsZero() {
		encID, err := encryptions.EncryptID(payrun.GeneratedBy.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt generated by ID: %w", err)
		}
		response.GeneratedBy = encID
	}

	return response, nil
}

