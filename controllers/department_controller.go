package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DepartmentController struct {
	service *services.DepartmentService
}

func NewDepartmentController() *DepartmentController {
	return &DepartmentController{
		service: services.NewDepartmentService(),
	}
}

// CreateDepartment creates a new department
func (dc *DepartmentController) CreateDepartment(c *fiber.Ctx) error {
	var req services.CreateDepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	department, err := dc.service.CreateDepartment(&req, companyID, authUserID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Department created successfully", department)
}

// ListDepartments retrieves all departments in the company
func (dc *DepartmentController) ListDepartments(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "50"), 10, 64)

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	departments, total, err := dc.service.ListDepartments(companyID, page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	if total == 0 {
		return constants.HTTPSuccess.OKWithoutData(c, "No departments found")
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Departments retrieved successfully", departments, page, limit, total)
}

// GetDepartmentByID retrieves a department by ID
func (dc *DepartmentController) GetDepartmentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	departmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid department ID")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	department, err := dc.service.GetDepartmentByID(departmentID, companyID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Department retrieved successfully", department)
}

// UpdateDepartment updates a department
func (dc *DepartmentController) UpdateDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	departmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid department ID")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	var req services.UpdateDepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err = dc.service.UpdateDepartment(departmentID, companyID, authUserID, &req)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Department updated successfully")
}

// DeleteDepartment soft deletes a department
func (dc *DepartmentController) DeleteDepartment(c *fiber.Ctx) error {
	id := c.Params("id")
	departmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid department ID")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	err = dc.service.DeleteDepartment(departmentID, companyID, authUserID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Department deleted successfully")
}
