package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentController struct {
	service *services.DocumentService
}

func NewDocumentController() *DocumentController {
	return &DocumentController{
		service: services.NewDocumentService(),
	}
}

// UploadDocument handles file upload
func (dc *DocumentController) UploadDocument(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "File is required")
	}

	// Parse form data
	category := c.FormValue("category")
	description := c.FormValue("description")
	employeeID := c.FormValue("employee_id")

	if category == "" {
		return constants.HTTPErrors.BadRequest(c, "Category is required")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	uploadedBy, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	req := &services.UploadDocumentRequest{
		Category:    category,
		Description: description,
		EmployeeID:  employeeID,
	}

	document, err := dc.service.UploadDocument(file, req, companyID, uploadedBy)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Document uploaded successfully", document)
}

// ListDocuments retrieves documents with filters
func (dc *DocumentController) ListDocuments(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	category := c.Query("category")
	employeeID := c.Query("employee_id")

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	req := &services.ListDocumentsRequest{
		Category:   category,
		EmployeeID: employeeID,
		Page:       page,
		Limit:      limit,
	}

	documents, total, err := dc.service.ListDocuments(req, companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Documents retrieved successfully", fiber.Map{
		"documents": documents,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// DeleteDocument removes a document
func (dc *DocumentController) DeleteDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	documentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid document ID")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	err = dc.service.DeleteDocument(documentID, companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Document deleted successfully")
}
