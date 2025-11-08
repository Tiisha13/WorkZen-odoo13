package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PayrollController struct {
	service *services.PayrollService
}

func NewPayrollController() *PayrollController {
	return &PayrollController{
		service: services.NewPayrollService(),
	}
}

// CreateConfiguration creates or updates payroll configuration
func (pc *PayrollController) CreateConfiguration(c *fiber.Ctx) error {
	var req services.CreatePayrollConfigurationRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	config, err := pc.service.CreateConfiguration(&req, companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Payroll configuration saved successfully", config)
}

// GetConfiguration retrieves payroll configuration
func (pc *PayrollController) GetConfiguration(c *fiber.Ctx) error {
	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	config, err := pc.service.GetConfiguration(companyID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Payroll configuration retrieved successfully", config)
}

// CreatePayrun generates monthly payroll for all employees
func (pc *PayrollController) CreatePayrun(c *fiber.Ctx) error {
	var req services.CreatePayrunRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	generatedBy, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	payrun, err := pc.service.CreatePayrun(&req, companyID, generatedBy)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Payrun generated successfully", payrun)
}

// ListPayruns retrieves all payruns for company
func (pc *PayrollController) ListPayruns(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	payruns, total, err := pc.service.ListPayruns(companyID, page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Payruns retrieved successfully", payruns, page, limit, total)
}

// GetEmployeePayroll retrieves payroll for specific employee and month
func (pc *PayrollController) GetEmployeePayroll(c *fiber.Ctx) error {
	employeeIDStr := c.Params("employee_id")
	employeeID, err := primitive.ObjectIDFromHex(employeeIDStr)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid employee ID")
	}

	month := c.Query("month") // YYYY-MM format
	if month == "" {
		return constants.HTTPErrors.BadRequest(c, "Month parameter is required")
	}

	payroll, err := pc.service.GetEmployeePayroll(employeeID, month)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Payroll retrieved successfully", payroll)
}

// MarkAsPaid marks a payroll record as paid
func (pc *PayrollController) MarkAsPaid(c *fiber.Ctx) error {
	id := c.Params("id")
	payrollID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid payroll ID")
	}

	err = pc.service.MarkAsPaid(payrollID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Payroll marked as paid successfully")
}
