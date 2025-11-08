package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PayrollService struct{}

func NewPayrollService() *PayrollService {
	return &PayrollService{}
}

// CreatePayrollConfigurationRequest for payroll settings
type CreatePayrollConfigurationRequest struct {
	PFEmployeePercent        float64 `json:"pf_employee_percent"`
	PFEmployerPercent        float64 `json:"pf_employer_percent"`
	ProfessionalTax          float64 `json:"professional_tax"`
	DefaultBasicPercent      float64 `json:"default_basic_percent"`
	DefaultHRAPercent        float64 `json:"default_hra_percent"`
	DefaultStandardAllowance float64 `json:"default_standard_allowance"`
	DefaultPerformanceBonus  float64 `json:"default_performance_bonus"`
	DefaultLTA               float64 `json:"default_lta"`
}

// CreateConfiguration creates or updates payroll configuration
func (s *PayrollService) CreateConfiguration(req *CreatePayrollConfigurationRequest, companyID primitive.ObjectID) (*models.PayrollConfiguration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configCollection := databases.MongoDBDatabase.Collection(collections.PayrollConfigurations)

	// Check if configuration exists
	var existing models.PayrollConfiguration
	err := configCollection.FindOne(ctx, bson.M{"company": companyID}).Decode(&existing)

	config := models.PayrollConfiguration{
		Company:                  companyID,
		PFEmployeePercent:        req.PFEmployeePercent,
		PFEmployerPercent:        req.PFEmployerPercent,
		ProfessionalTax:          req.ProfessionalTax,
		DefaultBasicPercent:      req.DefaultBasicPercent,
		DefaultHRAPercent:        req.DefaultHRAPercent,
		DefaultStandardAllowance: req.DefaultStandardAllowance,
		DefaultPerformanceBonus:  req.DefaultPerformanceBonus,
		DefaultLTA:               req.DefaultLTA,
		Currency:                 "INR",
	}

	if err == nil {
		// Update existing
		config.ID = existing.ID
		_, err = configCollection.ReplaceOne(ctx, bson.M{"_id": existing.ID}, config)
		if err != nil {
			return nil, err
		}
	} else {
		// Create new
		config.ID = primitive.NewObjectID()
		config.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		config.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
		_, err = configCollection.InsertOne(ctx, config)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// GetConfiguration retrieves payroll configuration for company
func (s *PayrollService) GetConfiguration(companyID primitive.ObjectID) (*models.PayrollConfiguration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configCollection := databases.MongoDBDatabase.Collection(collections.PayrollConfigurations)

	var config models.PayrollConfiguration
	err := configCollection.FindOne(ctx, bson.M{"company": companyID}).Decode(&config)
	if err != nil {
		return nil, errors.New("payroll configuration not found")
	}

	return &config, nil
}

// CreatePayrunRequest for generating payroll
type CreatePayrunRequest struct {
	Month string `json:"month" validate:"required"` // YYYY-MM
}

// CreatePayrun generates monthly payroll for all active employees
func (s *PayrollService) CreatePayrun(req *CreatePayrunRequest, companyID, generatedByID primitive.ObjectID) (*models.Payrun, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Collections
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	salaryCollection := databases.MongoDBDatabase.Collection(collections.SalaryStructures)
	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)
	payrunCollection := databases.MongoDBDatabase.Collection(collections.Payruns)
	configCollection := databases.MongoDBDatabase.Collection(collections.PayrollConfigurations)

	// Get payroll configuration
	var config models.PayrollConfiguration
	err := configCollection.FindOne(ctx, bson.M{"company": companyID}).Decode(&config)
	if err != nil {
		config = models.PayrollConfiguration{
			PFEmployeePercent: 12.0,
			PFEmployerPercent: 12.0,
			ProfessionalTax:   200.0,
		}
	}

	// Get all active employees
	cursor, err := usersCollection.Find(ctx, bson.M{
		"company": companyID,
		"status":  models.UserActive,
		"role":    bson.M{"$ne": models.RoleSuperAdmin},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var employees []models.User
	if err = cursor.All(ctx, &employees); err != nil {
		return nil, err
	}

	// Create payrun
	payrun := models.Payrun{
		ID:             primitive.NewObjectID(),
		Company:        companyID,
		Month:          req.Month,
		TotalEmployees: len(employees),
		Status:         models.PayrunGenerated,
		GeneratedBy:    generatedByID,
		GeneratedAt:    helpers.FormatDateTime(time.Now()),
	}
	payrun.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	payrun.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	var totalPayroll float64
	processedCount := 0
	missingBankCount := 0
	missingManagerCount := 0

	// Generate payroll for each employee
	for _, emp := range employees {
		// Get salary structure
		var salary models.SalaryStructure
		err := salaryCollection.FindOne(ctx, bson.M{
			"employee_id": emp.ID,
			"is_active":   true,
		}).Decode(&salary)
		if err != nil {
			continue // Skip if no salary structure
		}

		// Calculate deductions
		pfEmployee, pfEmployer, profTax := helpers.CalculateDeductions(salary.BasicSalary.Amount, &config)

		// Check warnings
		hasBankAccount := emp.BankDetails != nil && emp.BankDetails.AccountNumber != ""
		hasManager := !emp.ManagerID.IsZero()

		if !hasBankAccount {
			missingBankCount++
		}
		if !hasManager {
			missingManagerCount++
		}

		// Create payroll record
		payroll := models.Payroll{
			ID:                   primitive.NewObjectID(),
			EmployeeID:           emp.ID,
			Company:              companyID,
			PayrunID:             payrun.ID,
			Month:                req.Month,
			BasicSalary:          salary.BasicSalary.Amount,
			HouseRentAllowance:   salary.HouseRentAllowance.Amount,
			StandardAllowance:    salary.StandardAllowance.Amount,
			PerformanceBonus:     salary.PerformanceBonus.Amount,
			LeaveTravelAllowance: salary.LeaveTravelAllowance.Amount,
			FixedAllowance:       salary.FixedAllowance.Amount,
			GrossSalary:          salary.TotalEarnings,
			PFEmployee:           pfEmployee,
			PFEmployer:           pfEmployer,
			ProfessionalTax:      profTax,
			TotalDeductions:      pfEmployee + profTax,
			NetPay:               salary.TotalEarnings - (pfEmployee + profTax),
			HasBankAccount:       hasBankAccount,
			HasManager:           hasManager,
			Status:               models.PayrollProcessed,
			GeneratedBy:          generatedByID,
			GeneratedAt:          helpers.FormatDateTime(time.Now()),
		}
		payroll.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
		payroll.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

		_, err = payrollCollection.InsertOne(ctx, payroll)
		if err == nil {
			processedCount++
			totalPayroll += payroll.NetPay
		}
	}

	// Update payrun with totals
	payrun.ProcessedCount = processedCount
	payrun.TotalPayroll = totalPayroll
	payrun.MissingBankCount = missingBankCount
	payrun.MissingManagerCount = missingManagerCount

	_, err = payrunCollection.InsertOne(ctx, payrun)
	if err != nil {
		return nil, fmt.Errorf("failed to create payrun: %w", err)
	}

	return &payrun, nil
}

// ListPayruns retrieves payruns with pagination
func (s *PayrollService) ListPayruns(companyID primitive.ObjectID, page, limit int64) ([]models.Payrun, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payrunCollection := databases.MongoDBDatabase.Collection(collections.Payruns)

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := payrunCollection.Find(ctx, bson.M{"company": companyID}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var payruns []models.Payrun
	if err = cursor.All(ctx, &payruns); err != nil {
		return nil, 0, err
	}

	total, err := payrunCollection.CountDocuments(ctx, bson.M{"company": companyID})
	if err != nil {
		return nil, 0, err
	}

	return payruns, total, nil
}

// GetEmployeePayroll retrieves payroll for specific employee and month
func (s *PayrollService) GetEmployeePayroll(employeeID primitive.ObjectID, month string) (*models.Payroll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)

	var payroll models.Payroll
	err := payrollCollection.FindOne(ctx, bson.M{
		"employee_id": employeeID,
		"month":       month,
	}).Decode(&payroll)
	if err != nil {
		return nil, errors.New("payroll not found")
	}

	return &payroll, nil
}

// MarkAsPaid marks a payroll record as paid
func (s *PayrollService) MarkAsPaid(payrollID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)

	result, err := payrollCollection.UpdateOne(
		ctx,
		bson.M{"_id": payrollID},
		bson.M{
			"$set": bson.M{
				"status":     models.PayrollPaid,
				"paid_at":    helpers.FormatDateTime(time.Now()),
				"updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("payroll not found")
	}

	return nil
}
