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
)

type SalaryService struct{}

func NewSalaryService() *SalaryService {
	return &SalaryService{}
}

// CreateSalaryStructureRequest for creating salary structure
type CreateSalaryStructureRequest struct {
	EmployeeID    string  `json:"employee_id" validate:"required"`
	MonthlyWage   float64 `json:"monthly_wage" validate:"required"`
	EffectiveFrom string  `json:"effective_from"` // YYYY-MM-DD
}

// CreateSalaryStructure creates a new salary structure for an employee
func (s *SalaryService) CreateSalaryStructure(req *CreateSalaryStructureRequest, companyID primitive.ObjectID) (*models.SalaryStructure, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	salaryCollection := databases.MongoDBDatabase.Collection(collections.SalaryStructures)
	configCollection := databases.MongoDBDatabase.Collection(collections.PayrollConfigurations)

	// Parse employee ID
	employeeID, err := primitive.ObjectIDFromHex(req.EmployeeID)
	if err != nil {
		return nil, errors.New("invalid employee ID")
	}

	// Deactivate existing salary structures for this employee
	_, err = salaryCollection.UpdateMany(
		ctx,
		bson.M{"employee_id": employeeID, "is_active": true},
		bson.M{"$set": bson.M{"is_active": false}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate old structures: %w", err)
	}

	// Get payroll configuration for company
	var config models.PayrollConfiguration
	err = configCollection.FindOne(ctx, bson.M{"company": companyID}).Decode(&config)
	if err != nil {
		// Use default configuration if not found
		config = models.PayrollConfiguration{
			DefaultBasicPercent:      50.0,
			DefaultHRAPercent:        50.0,
			DefaultStandardAllowance: 16.67,
			DefaultPerformanceBonus:  8.33,
			DefaultLTA:               8.33,
			PFEmployeePercent:        12.0,
			PFEmployerPercent:        12.0,
			ProfessionalTax:          200.0,
		}
	}

	// Calculate salary components using helper
	structure, err := helpers.CalculateSalaryComponents(req.MonthlyWage, &config)
	if err != nil {
		return nil, err
	}

	// Set additional fields
	structure.ID = primitive.NewObjectID()
	structure.EmployeeID = employeeID
	structure.Company = companyID
	structure.Currency = "INR"
	structure.IsActive = true

	if req.EffectiveFrom != "" {
		structure.EffectiveFrom = req.EffectiveFrom
	} else {
		structure.EffectiveFrom = helpers.FormatDate(time.Now())
	}

	// Calculate deductions
	pfEmployee, _, profTax := helpers.CalculateDeductions(structure.BasicSalary.Amount, &config)
	structure.TotalDeductions = pfEmployee + profTax
	structure.NetPay = helpers.CalculateNetPay(structure.TotalEarnings, structure.TotalDeductions)

	structure.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	structure.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Insert to database
	_, err = salaryCollection.InsertOne(ctx, structure)
	if err != nil {
		return nil, fmt.Errorf("failed to create salary structure: %w", err)
	}

	return structure, nil
}

// GetSalaryStructure retrieves active salary structure for an employee
func (s *SalaryService) GetSalaryStructure(employeeID primitive.ObjectID) (*models.SalaryStructure, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	salaryCollection := databases.MongoDBDatabase.Collection(collections.SalaryStructures)

	var structure models.SalaryStructure
	err := salaryCollection.FindOne(ctx, bson.M{
		"employee_id": employeeID,
		"is_active":   true,
	}).Decode(&structure)
	if err != nil {
		return nil, errors.New("salary structure not found")
	}

	return &structure, nil
}

// UpdateSalaryStructureRequest for updating salary
type UpdateSalaryStructureRequest struct {
	MonthlyWage   float64 `json:"monthly_wage" validate:"required"`
	EffectiveFrom string  `json:"effective_from"` // YYYY-MM-DD
}

// UpdateSalaryStructure updates salary by creating a new structure
func (s *SalaryService) UpdateSalaryStructure(employeeID primitive.ObjectID, req *UpdateSalaryStructureRequest, companyID primitive.ObjectID) (*models.SalaryStructure, error) {
	// Create new structure request
	createReq := &CreateSalaryStructureRequest{
		EmployeeID:    employeeID.Hex(),
		MonthlyWage:   req.MonthlyWage,
		EffectiveFrom: req.EffectiveFrom,
	}

	return s.CreateSalaryStructure(createReq, companyID)
}
