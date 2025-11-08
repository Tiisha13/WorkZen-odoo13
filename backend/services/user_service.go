package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
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

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// CreateUserRequest for creating a new user/employee
type CreateUserRequest struct {
	FirstName    string              `json:"first_name" validate:"required"`
	LastName     string              `json:"last_name" validate:"required"`
	Email        string              `json:"email" validate:"required,email"`
	Phone        string              `json:"phone"`
	Role         models.Role         `json:"role" validate:"required"`
	Designation  string              `json:"designation"`
	DepartmentID *primitive.ObjectID `json:"department_id"`
	ManagerID    *primitive.ObjectID `json:"manager_id"`
	DateOfJoin   string              `json:"date_of_join"` // YYYY-MM-DD
}

// CreateUser creates a new user within a company
func (s *UserService) CreateUser(req *CreateUserRequest, companyID primitive.ObjectID) (*models.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Check if email already exists
	count, err := usersCollection.CountDocuments(ctx, bson.M{"email": req.Email, "company": companyID})
	if err != nil {
		return nil, "", err
	}
	if count > 0 {
		return nil, "", errors.New("user with this email already exists in this company")
	}

	// Get company to generate login ID
	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)
	var company models.Company
	err = companiesCollection.FindOne(ctx, bson.M{"_id": companyID}).Decode(&company)
	if err != nil {
		return nil, "", errors.New("company not found")
	}

	// Generate company code
	companyCode := company.Name
	if len(companyCode) > 4 {
		companyCode = companyCode[:4]
	}

	// Count existing users to get serial number
	totalUsers, _ := usersCollection.CountDocuments(ctx, bson.M{"company": companyID})
	serial := int(totalUsers) + 1

	// Parse join date
	joinDate := time.Now()
	if req.DateOfJoin != "" {
		parsedDate, err := helpers.ParseDate(req.DateOfJoin)
		if err == nil {
			joinDate = parsedDate
		}
	}

	// Generate login ID
	username := helpers.GenerateLoginID(companyCode, req.FirstName, req.LastName, joinDate, serial)

	// Generate random temporary password (8 characters)
	tempPassword := generateRandomPassword(8)
	hashedPassword := encryptions.HashPassword(tempPassword)

	// Create user
	user := models.User{
		ID:          primitive.NewObjectID(),
		Username:    username,
		Email:       req.Email,
		Password:    hashedPassword,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        req.Role,
		Designation: req.Designation,
		DateOfJoin:  joinDate.Format("2006-01-02"),
		Status:      models.UserActive,
		Phone:       req.Phone,
		Company:     companyID,
	}

	if req.DepartmentID != nil {
		user.DepartmentID = *req.DepartmentID
	}
	if req.ManagerID != nil {
		user.ManagerID = *req.ManagerID
	}

	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Remove password from response
	user.Password = ""

	return &user, tempPassword, nil
}

// ListUsers retrieves users with pagination and filters
func (s *UserService) ListUsers(companyID primitive.ObjectID, filters map[string]interface{}, page, limit int64) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Add company filter
	filters["company"] = companyID

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := usersCollection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Remove passwords
	for i := range users {
		users[i].Password = ""
	}

	total, err := usersCollection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Remove password
	user.Password = ""

	return &user, nil
}

// UpdateUserStatus updates user active/inactive status
func (s *UserService) UpdateUserStatus(userID primitive.ObjectID, status models.UserStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdateBankDetailsRequest for updating bank information
type UpdateBankDetailsRequest struct {
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	IFSCCode      string `json:"ifsc_code"`
	BranchName    string `json:"branch_name"`
	PANNo         string `json:"pan_no"`
	UANNo         string `json:"uan_no"`
}

// UpdateBankDetails updates user's bank information
func (s *UserService) UpdateBankDetails(userID primitive.ObjectID, req *UpdateBankDetailsRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	bankDetails := models.BankDetails{
		AccountNumber: req.AccountNumber,
		BankName:      req.BankName,
		IFSCCode:      req.IFSCCode,
		BranchName:    req.BranchName,
		PANNo:         req.PANNo,
		UANNo:         req.UANNo,
	}

	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"bank_details": bankDetails,
				"updated_at":   primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DeleteUser soft deletes a user (marks as inactive)
func (s *UserService) DeleteUser(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"status":     models.UserInactive,
				"deleted_at": primitive.NewDateTimeFromTime(time.Now()),
				"updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// generateRandomPassword generates a random password
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
