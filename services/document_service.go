package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentService struct{}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

// UploadDocumentRequest for file upload
type UploadDocumentRequest struct {
	Category    string `json:"category" validate:"required"`
	Description string `json:"description"`
	EmployeeID  string `json:"employee_id"` // Optional - for employee-specific documents
}

// UploadDocument handles file upload and document creation
func (s *DocumentService) UploadDocument(file *multipart.FileHeader, req *UploadDocumentRequest, companyID, uploadedByID primitive.ObjectID) (*models.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	documentCollection := databases.MongoDBDatabase.Collection(collections.Documents)

	// Validate category
	validCategories := []models.DocumentCategory{
		models.DocumentCategoryResume,
		models.DocumentCategoryIDProof,
		models.DocumentCategoryPayslip,
		models.DocumentCategoryPolicy,
		models.DocumentCategoryReport,
		models.DocumentCategoryOther,
	}
	isValid := false
	categoryEnum := models.DocumentCategory(req.Category)
	for _, cat := range validCategories {
		if categoryEnum == cat {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid document category")
	}

	// Parse employee ID if provided
	var employeeID primitive.ObjectID
	if req.EmployeeID != "" {
		empID, err := primitive.ObjectIDFromHex(req.EmployeeID)
		if err != nil {
			return nil, errors.New("invalid employee_id format")
		}
		employeeID = empID
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueID := helpers.GetNewUUID()
	newFilename := fmt.Sprintf("%s%s", uniqueID, ext)

	// Build path: /assets/uploads/{companyID}/{category}/{YYYY}/{MM}/
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")
	uploadPath := filepath.Join("assets", "uploads", companyID.Hex(), req.Category, year, month)

	// Create directory structure
	err := os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(uploadPath, newFilename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create document record
	document := models.Document{
		ID:          primitive.NewObjectID(),
		Company:     companyID,
		EmployeeID:  employeeID,
		Category:    categoryEnum,
		FileName:    file.Filename,
		FilePath:    filePath,
		FileType:    file.Header.Get("Content-Type"),
		Size:        file.Size,
		Description: req.Description,
		UploadedBy:  uploadedByID,
		IsPrivate:   req.Category == string(models.DocumentCategoryPayslip),
	}
	document.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	document.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = documentCollection.InsertOne(ctx, document)
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save document record: %w", err)
	}

	return &document, nil
}

// ListDocumentsRequest for filtering documents
type ListDocumentsRequest struct {
	Category   string `json:"category"`
	EmployeeID string `json:"employee_id"`
	Page       int64  `json:"page"`
	Limit      int64  `json:"limit"`
}

// ListDocuments retrieves documents with filtering
func (s *DocumentService) ListDocuments(req *ListDocumentsRequest, companyID primitive.ObjectID) ([]models.Document, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	documentCollection := databases.MongoDBDatabase.Collection(collections.Documents)

	// Build filter
	filter := bson.M{"company": companyID}
	if req.Category != "" {
		filter["category"] = req.Category
	}
	if req.EmployeeID != "" {
		empID, err := primitive.ObjectIDFromHex(req.EmployeeID)
		if err == nil {
			filter["employee_id"] = empID
		}
	}

	// Pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := documentCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var documents []models.Document
	if err = cursor.All(ctx, &documents); err != nil {
		return nil, 0, err
	}

	total, err := documentCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// DeleteDocument removes document record and file
func (s *DocumentService) DeleteDocument(documentID, companyID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	documentCollection := databases.MongoDBDatabase.Collection(collections.Documents)

	// Find document
	var document models.Document
	err := documentCollection.FindOne(ctx, bson.M{
		"_id":     documentID,
		"company": companyID,
	}).Decode(&document)
	if err != nil {
		return errors.New("document not found")
	}

	// Delete file
	if document.FilePath != "" {
		err = os.Remove(document.FilePath)
		if err != nil && !os.IsNotExist(err) {
			// Log error but continue with database deletion
			fmt.Printf("Warning: Failed to delete file %s: %v\n", document.FilePath, err)
		}
	}

	// Delete from database
	result, err := documentCollection.DeleteOne(ctx, bson.M{"_id": documentID})
	if err != nil || result.DeletedCount == 0 {
		return errors.New("failed to delete document record")
	}

	return nil
}
