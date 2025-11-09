package helpers

import (
	"errors"

	"api.workzen.odoo/databases/models"
)

// CalculateSalaryComponents computes the actual amounts for each salary component based on monthly wage
func CalculateSalaryComponents(monthlyWage float64, config *models.PayrollConfiguration) (*models.SalaryStructure, error) {
	if monthlyWage <= 0 {
		return nil, errors.New("monthly wage must be greater than zero")
	}

	if config == nil {
		// Use default percentages if no configuration provided
		config = &models.PayrollConfiguration{
			DefaultBasicPercent:      40.0, // 40% of monthly wage
			DefaultHRAPercent:        40.0, // 40% of Basic (16% of monthly wage)
			DefaultStandardAllowance: 15.0, // 15% of monthly wage
			DefaultPerformanceBonus:  10.0, // 10% of monthly wage
			DefaultLTA:               10.0, // 10% of monthly wage
		}
	}

	structure := &models.SalaryStructure{
		MonthlyWage: monthlyWage,
		YearlyWage:  monthlyWage * 12,
		WageType:    models.WageTypeFixed,
		Currency:    "INR",
		IsActive:    true,
	}

	// Calculate Basic Salary (percentage of monthly wage)
	basicAmount := (config.DefaultBasicPercent / 100) * monthlyWage
	structure.BasicSalary = models.SalaryComponent{
		Name:   "Basic Salary",
		Type:   models.ComponentTypePercentage,
		Value:  config.DefaultBasicPercent,
		Amount: basicAmount,
	}

	// Calculate HRA (percentage of Basic Salary)
	hraAmount := (config.DefaultHRAPercent / 100) * basicAmount
	structure.HouseRentAllowance = models.SalaryComponent{
		Name:   "House Rent Allowance",
		Type:   models.ComponentTypePercentage,
		Value:  config.DefaultHRAPercent,
		Amount: hraAmount,
	}

	// Calculate Standard Allowance (percentage of monthly wage)
	standardAllowanceAmount := (config.DefaultStandardAllowance / 100) * monthlyWage
	structure.StandardAllowance = models.SalaryComponent{
		Name:   "Standard Allowance",
		Type:   models.ComponentTypePercentage,
		Value:  config.DefaultStandardAllowance,
		Amount: standardAllowanceAmount,
	}

	// Calculate Performance Bonus (percentage of monthly wage)
	performanceBonusAmount := (config.DefaultPerformanceBonus / 100) * monthlyWage
	structure.PerformanceBonus = models.SalaryComponent{
		Name:   "Performance Bonus",
		Type:   models.ComponentTypePercentage,
		Value:  config.DefaultPerformanceBonus,
		Amount: performanceBonusAmount,
	}

	// Calculate LTA (percentage of monthly wage)
	ltaAmount := (config.DefaultLTA / 100) * monthlyWage
	structure.LeaveTravelAllowance = models.SalaryComponent{
		Name:   "Leave Travel Allowance",
		Type:   models.ComponentTypePercentage,
		Value:  config.DefaultLTA,
		Amount: ltaAmount,
	}

	// Calculate Fixed Allowance (remaining amount to match monthly wage)
	totalComponentsAmount := basicAmount + hraAmount + standardAllowanceAmount +
		performanceBonusAmount + ltaAmount
	fixedAllowanceAmount := monthlyWage - totalComponentsAmount

	if fixedAllowanceAmount < 0 {
		return nil, errors.New("total component values exceed monthly wage")
	}

	structure.FixedAllowance = models.SalaryComponent{
		Name:   "Fixed Allowance",
		Type:   models.ComponentTypeFixed,
		Value:  fixedAllowanceAmount,
		Amount: fixedAllowanceAmount,
	}

	// Calculate total earnings (should equal monthly wage)
	structure.TotalEarnings = totalComponentsAmount + fixedAllowanceAmount

	return structure, nil
}

// CalculateDeductions computes PF and Professional Tax based on configuration
func CalculateDeductions(basicSalary float64, config *models.PayrollConfiguration) (pfEmployee, pfEmployer, profTax float64) {
	if config == nil {
		// Default values
		config = &models.PayrollConfiguration{
			PFEmployeePercent: 12.0,
			PFEmployerPercent: 12.0,
			ProfessionalTax:   200.0,
		}
	}

	// Calculate PF (on basic salary)
	pfEmployee = (config.PFEmployeePercent / 100) * basicSalary
	pfEmployer = (config.PFEmployerPercent / 100) * basicSalary
	profTax = config.ProfessionalTax

	return pfEmployee, pfEmployer, profTax
}

// CalculateNetPay computes the final take-home salary after deductions
func CalculateNetPay(grossSalary, totalDeductions float64) float64 {
	netPay := grossSalary - totalDeductions
	if netPay < 0 {
		return 0
	}
	return netPay
}

// ValidateComponentTotal ensures that the sum of all components does not exceed monthly wage
func ValidateComponentTotal(components []models.SalaryComponent, monthlyWage float64) error {
	total := 0.0
	for _, comp := range components {
		total += comp.Amount
	}

	if total > monthlyWage {
		return errors.New("total of all salary components exceeds monthly wage")
	}

	return nil
}

// RecalculateStructure recomputes all salary components when wage or percentages change
func RecalculateStructure(structure *models.SalaryStructure, config *models.PayrollConfiguration) error {
	if structure == nil {
		return errors.New("salary structure cannot be nil")
	}

	newStructure, err := CalculateSalaryComponents(structure.MonthlyWage, config)
	if err != nil {
		return err
	}

	// Update the structure with recalculated values
	structure.BasicSalary = newStructure.BasicSalary
	structure.HouseRentAllowance = newStructure.HouseRentAllowance
	structure.StandardAllowance = newStructure.StandardAllowance
	structure.PerformanceBonus = newStructure.PerformanceBonus
	structure.LeaveTravelAllowance = newStructure.LeaveTravelAllowance
	structure.FixedAllowance = newStructure.FixedAllowance
	structure.TotalEarnings = newStructure.TotalEarnings
	structure.YearlyWage = newStructure.YearlyWage

	return nil
}
