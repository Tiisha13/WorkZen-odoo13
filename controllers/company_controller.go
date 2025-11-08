package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/encryptions"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyController struct {
	service *services.CompanyService
}

func NewCompanyController() *CompanyController {
	return &CompanyController{
		service: services.NewCompanyService(),
	}
}

// CreateCompany creates a new company (SuperAdmin only)
func (cc *CompanyController) CreateCompany(c *fiber.Ctx) error {
	var req services.CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	company, err := cc.service.CreateCompany(&req)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Company created successfully", company)
}

// ListCompanies retrieves all companies (SuperAdmin only)
func (cc *CompanyController) ListCompanies(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)

	companies, total, err := cc.service.ListCompanies(page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Companies retrieved successfully", companies, page, limit, total)
}

// GetCompanyByID retrieves a company by ID
func (cc *CompanyController) GetCompanyByID(c *fiber.Ctx) error {
	encryptedID := c.Params("id")

	// Decrypt the company ID
	decryptedID, err := encryptions.DecryptID(encryptedID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid company ID")
	}

	companyID, err := primitive.ObjectIDFromHex(decryptedID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid company ID")
	}

	company, err := cc.service.GetCompanyByID(companyID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Company retrieved successfully", company)
}

// ApproveCompany approves a company (SuperAdmin only)
func (cc *CompanyController) ApproveCompany(c *fiber.Ctx) error {
	encryptedID := c.Params("id")

	// Decrypt the company ID
	decryptedID, err := encryptions.DecryptID(encryptedID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid company ID")
	}

	companyID, err := primitive.ObjectIDFromHex(decryptedID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid company ID")
	}

	err = cc.service.ApproveCompany(companyID, primitive.NilObjectID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Company approved successfully")
}

// DeactivateCompany deactivates a company (SuperAdmin only)
func (cc *CompanyController) DeactivateCompany(c *fiber.Ctx) error {
	encryptedID := c.Params("id")

	// Decrypt the company ID
	decryptedID, err := encryptions.DecryptID(encryptedID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid company ID")
	}

	companyID, err := primitive.ObjectIDFromHex(decryptedID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid company ID")
	}

	err = cc.service.DeactivateCompany(companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Company deactivated successfully")
}
