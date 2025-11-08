package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/encryptions"
	"api.workzen.odoo/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DepartmentService struct{}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{}
}

// DepartmentResponse with encrypted IDs
type DepartmentResponse struct {
	ID          string `json:"id"`
	Company     string `json:"company"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	HeadID      string `json:"head_id,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// convertDepartmentToResponse converts Department model to response with encrypted IDs
func convertDepartmentToResponse(dept *models.Department) (*DepartmentResponse, error) {
	encryptedID, err := encryptions.EncryptID(dept.ID.Hex())
	if err != nil {
		return nil, err
	}

	encryptedCompanyID, err := encryptions.EncryptID(dept.Company.Hex())
	if err != nil {
		return nil, err
	}

	response := &DepartmentResponse{
		ID:          encryptedID,
		Company:     encryptedCompanyID,
		Name:        dept.Name,
		Description: dept.Description,
		CreatedAt:   dept.CreatedAt.Time().Unix(),
		UpdatedAt:   dept.UpdatedAt.Time().Unix(),
	}

	if !dept.HeadID.IsZero() {
		encryptedHeadID, err := encryptions.EncryptID(dept.HeadID.Hex())
		if err != nil {
			return nil, err
		}
		response.HeadID = encryptedHeadID
	}

	return response, nil
}

// CreateDepartmentRequest for creating a new department
type CreateDepartmentRequest struct {
	Name        string              `json:"name" validate:"required"`
	Description string              `json:"description"`
	HeadID      *primitive.ObjectID `json:"head_id"`
}

// CreateDepartment creates a new department within a company
func (s *DepartmentService) CreateDepartment(req *CreateDepartmentRequest, companyID, authUserID primitive.ObjectID) (*DepartmentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)

	// Check if department with same name already exists in this company
	count, err := departmentsCollection.CountDocuments(ctx, helpers.AddNotDeletedFilter(bson.M{
		"name":    req.Name,
		"company": companyID,
	}))
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("department with this name already exists in your company")
	}

	department := models.Department{
		ID:          primitive.NewObjectID(),
		Company:     companyID,
		Name:        req.Name,
		Description: req.Description,
	}

	if req.HeadID != nil {
		department.HeadID = *req.HeadID
	}

	// Set timestamps with creator information
	department.CreatedAt, department.CreatedBy = helpers.SetCreatedTimestamp(authUserID)
	department.UpdatedAt, department.UpdatedBy = helpers.SetUpdatedTimestamp(authUserID)
	department.IsDeleted = false

	_, err = departmentsCollection.InsertOne(ctx, department)
	if err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	// Convert to response with encrypted IDs
	return convertDepartmentToResponse(&department)
}

// ListDepartments retrieves all departments for a company with pagination
func (s *DepartmentService) ListDepartments(companyID primitive.ObjectID, page, limit int64) ([]DepartmentResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)

	// Filter by company and exclude soft-deleted records
	query := helpers.AddNotDeletedFilter(bson.M{
		"company": companyID,
	})

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := departmentsCollection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var departments []models.Department
	if err = cursor.All(ctx, &departments); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := departmentsCollection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response with encrypted IDs
	var responses []DepartmentResponse
	for _, dept := range departments {
		response, err := convertDepartmentToResponse(&dept)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

// GetDepartmentByID retrieves a department by ID (company-scoped)
func (s *DepartmentService) GetDepartmentByID(departmentID, companyID primitive.ObjectID) (*DepartmentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)

	var department models.Department
	err := departmentsCollection.FindOne(ctx, helpers.AddNotDeletedFilter(bson.M{
		"_id":     departmentID,
		"company": companyID,
	})).Decode(&department)
	if err != nil {
		return nil, errors.New("department not found")
	}

	return convertDepartmentToResponse(&department)
}

// UpdateDepartmentRequest for updating a department
type UpdateDepartmentRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	HeadID      *primitive.ObjectID `json:"head_id"`
}

// UpdateDepartment updates a department (company-scoped)
func (s *DepartmentService) UpdateDepartment(departmentID, companyID, authUserID primitive.ObjectID, req *UpdateDepartmentRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)

	// Check if department with same name already exists (excluding current department)
	if req.Name != "" {
		count, err := departmentsCollection.CountDocuments(ctx, helpers.AddNotDeletedFilter(bson.M{
			"name":    req.Name,
			"company": companyID,
			"_id":     bson.M{"$ne": departmentID},
		}))
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("department with this name already exists in your company")
		}
	}

	updateFields := bson.M{}
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Description != "" {
		updateFields["description"] = req.Description
	}
	if req.HeadID != nil {
		updateFields["head_id"] = *req.HeadID
	}

	// Set update timestamp
	updatedAt, updatedBy := helpers.SetUpdatedTimestamp(authUserID)
	updateFields["updated_at"] = updatedAt
	updateFields["updated_by"] = updatedBy

	result, err := departmentsCollection.UpdateOne(
		ctx,
		helpers.AddNotDeletedFilter(bson.M{
			"_id":     departmentID,
			"company": companyID,
		}),
		bson.M{"$set": updateFields},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("department not found or you don't have permission")
	}

	return nil
}

// DeleteDepartment soft deletes a department (company-scoped)
func (s *DepartmentService) DeleteDepartment(departmentID, companyID, authUserID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)

	// Check if any employees are assigned to this department
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	count, err := usersCollection.CountDocuments(ctx, helpers.AddNotDeletedFilter(bson.M{
		"department_id": departmentID,
		"company":       companyID,
	}))
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete department: %d employees are assigned to it", count)
	}

	deletedAt, deletedBy := helpers.SetDeletedTimestamp(authUserID)

	result, err := departmentsCollection.UpdateOne(
		ctx,
		helpers.AddNotDeletedFilter(bson.M{
			"_id":     departmentID,
			"company": companyID,
		}),
		bson.M{
			"$set": bson.M{
				"is_deleted": true,
				"deleted_at": &deletedAt,
				"deleted_by": &deletedBy,
				"updated_at": deletedAt,
				"updated_by": deletedBy,
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("department not found or already deleted")
	}

	return nil
}
