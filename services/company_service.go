package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CompanyService struct{}

func NewCompanyService() *CompanyService {
	return &CompanyService{}
}

// CreateCompanyRequest for creating a new company
type CreateCompanyRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Industry string `json:"industry"`
	Website  string `json:"website"`
}

// CreateCompany creates a new company (SuperAdmin only)
func (s *CompanyService) CreateCompany(req *CreateCompanyRequest) (*CompanyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)

	// Check if email already exists
	count, err := companiesCollection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("company with this email already exists")
	}

	company := models.Company{
		ID:         primitive.NewObjectID(),
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Industry:   req.Industry,
		Website:    req.Website,
		IsApproved: true, // Auto-approved when created by SuperAdmin
		IsActive:   true,
	}
	company.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	company.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = companiesCollection.InsertOne(ctx, company)
	if err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	// Convert to response with encrypted IDs
	companyResponse, err := convertCompanyToResponse(&company)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare company response: %w", err)
	}

	return companyResponse, nil
}

// ListCompanies retrieves all companies with pagination
func (s *CompanyService) ListCompanies(page, limit int64) ([]CompanyResponse, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := companiesCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var companies []models.Company
	if err = cursor.All(ctx, &companies); err != nil {
		return nil, 0, err
	}

	// Convert to response with encrypted IDs
	companyResponses := make([]CompanyResponse, 0, len(companies))
	for i := range companies {
		companyResp, err := convertCompanyToResponse(&companies[i])
		if err != nil {
			return nil, 0, fmt.Errorf("failed to prepare company response: %w", err)
		}
		companyResponses = append(companyResponses, *companyResp)
	}

	total, err := companiesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return companyResponses, total, nil
}

// GetCompanyByID retrieves a company by ID
func (s *CompanyService) GetCompanyByID(companyID primitive.ObjectID) (*CompanyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)

	var company models.Company
	err := companiesCollection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
	if err != nil {
		return nil, errors.New("company not found")
	}

	// Convert to response with encrypted IDs
	companyResponse, err := convertCompanyToResponse(&company)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare company response: %w", err)
	}

	return companyResponse, nil
}

// ApproveCompany approves a pending company signup
func (s *CompanyService) ApproveCompany(companyID, approvedByID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Update company
	result, err := companiesCollection.UpdateOne(
		ctx,
		bson.M{"_id": companyID},
		bson.M{
			"$set": bson.M{
				"is_approved": true,
				"approved_by": approvedByID,
				"updated_at":  primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("company not found")
	}

	// Activate admin user for this company
	_, err = usersCollection.UpdateMany(
		ctx,
		bson.M{
			"company": companyID,
			"role":    models.RoleAdmin,
		},
		bson.M{
			"$set": bson.M{
				"status":     models.UserActive,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to activate admin user: %w", err)
	}

	return nil
}

// DeactivateCompany deactivates a company and all its users
func (s *CompanyService) DeactivateCompany(companyID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Deactivate company
	result, err := companiesCollection.UpdateOne(
		ctx,
		bson.M{"_id": companyID},
		bson.M{
			"$set": bson.M{
				"is_active":  false,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("company not found")
	}

	// Deactivate all users of this company
	_, err = usersCollection.UpdateMany(
		ctx,
		bson.M{"company": companyID},
		bson.M{
			"$set": bson.M{
				"status":     models.UserInactive,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to deactivate users: %w", err)
	}

	return nil
}
