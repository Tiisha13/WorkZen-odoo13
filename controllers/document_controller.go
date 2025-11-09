package controllers

import (
	"log"
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/helpers"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
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
	employeeIDStr := c.FormValue("employee_id")

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

	// Decrypt employee_id if provided (it's encrypted from frontend)
	var employeeID string
	if employeeIDStr != "" {
		empID, err := helpers.DecryptObjectID(employeeIDStr)
		if err != nil {
			return constants.HTTPErrors.BadRequest(c, "Invalid employee_id")
		}
		employeeID = empID.Hex()
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

	// Get current user
	user, err := middlewares.GetAuthUser(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	// For regular employees, filter to show only their own documents
	// HR and Admin can see all documents
	if !user.IsSuperAdmin && user.Role != "admin" && user.Role != "hr" {
		// Override employee_id to current user's ID
		userID, _ := middlewares.GetAuthUserID(c)
		employeeID = userID.Hex()
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

	// Convert to responses with encrypted IDs
	documentResponses := make([]services.DocumentResponse, 0, len(documents))
	for _, doc := range documents {
		docResp, err := services.ConvertDocumentToResponse(&doc)
		if err != nil {
			return constants.HTTPErrors.InternalServerError(c, "Failed to encrypt document IDs")
		}
		documentResponses = append(documentResponses, *docResp)
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Documents retrieved successfully", documentResponses, page, limit, total)
}

// DeleteDocument removes a document
func (dc *DocumentController) DeleteDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	documentID, err := helpers.DecryptObjectID(id)
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

// ViewDocument serves the document file for viewing (images, videos, audio, PDFs)
func (dc *DocumentController) ViewDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	documentID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid document ID")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	document, err := dc.service.GetDocumentByID(documentID, companyID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, "Document not found: "+err.Error())
	}

	// Check if file exists
	if document.FilePath == "" {
		return constants.HTTPErrors.NotFound(c, "Document file path is empty")
	}

	// Log for debugging
	log.Printf("📄 Attempting to serve file: %s", document.FilePath)

	// Send file for viewing in browser
	if err := c.SendFile(document.FilePath); err != nil {
		log.Printf("❌ Error serving file %s: %v", document.FilePath, err)
		return constants.HTTPErrors.InternalServerError(c, "Failed to send file: "+err.Error())
	}

	return nil
}

// DownloadDocument serves the document file for download
func (dc *DocumentController) DownloadDocument(c *fiber.Ctx) error {
	id := c.Params("id")
	documentID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid document ID")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	document, err := dc.service.GetDocumentByID(documentID, companyID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, "Document not found: "+err.Error())
	}

	// Check if file exists
	if document.FilePath == "" {
		return constants.HTTPErrors.NotFound(c, "Document file path is empty")
	}

	// Set headers for download
	c.Set("Content-Disposition", "attachment; filename=\""+document.FileName+"\"")

	if err := c.SendFile(document.FilePath); err != nil {
		return constants.HTTPErrors.InternalServerError(c, "Failed to send file: "+err.Error())
	}

	return nil
}
