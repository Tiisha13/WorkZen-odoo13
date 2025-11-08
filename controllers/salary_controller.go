package controllers

import (
	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SalaryController struct {
	service *services.SalaryService
}

func NewSalaryController() *SalaryController {
	return &SalaryController{
		service: services.NewSalaryService(),
	}
}

// CreateSalaryStructure creates a new salary structure for employee
func (sc *SalaryController) CreateSalaryStructure(c *fiber.Ctx) error {
	var req services.CreateSalaryStructureRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	salary, err := sc.service.CreateSalaryStructure(&req, companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Salary structure created successfully", salary)
}

// GetSalaryStructure retrieves active salary structure for employee
func (sc *SalaryController) GetSalaryStructure(c *fiber.Ctx) error {
	employeeIDStr := c.Params("employee_id")
	employeeID, err := primitive.ObjectIDFromHex(employeeIDStr)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid employee ID")
	}

	salary, err := sc.service.GetSalaryStructure(employeeID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Salary structure retrieved successfully", salary)
}

// UpdateSalaryStructure updates salary structure (creates new version)
func (sc *SalaryController) UpdateSalaryStructure(c *fiber.Ctx) error {
	employeeIDStr := c.Params("employee_id")
	employeeID, err := primitive.ObjectIDFromHex(employeeIDStr)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid employee ID")
	}

	var req services.UpdateSalaryStructureRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	salary, err := sc.service.UpdateSalaryStructure(employeeID, &req, companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Salary structure updated successfully", salary)
}
