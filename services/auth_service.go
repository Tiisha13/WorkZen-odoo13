// Package services provides business logic and service layer implementations for the WorkZen HRMS API,
// including authentication, user management, and various HR-related operations.
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/encryptions"
	"api.workzen.odoo/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// SignupRequest represents company signup request
type SignupRequest struct {
	CompanyName string `json:"company_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required"`
	Industry    string `json:"industry"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UserResponse represents user data with encrypted IDs for API responses
type UserResponse struct {
	ID               string              `json:"id,omitempty"`
	Username         string              `json:"username"`
	Email            string              `json:"email"`
	FirstName        string              `json:"first_name"`
	LastName         string              `json:"last_name"`
	Role             models.Role         `json:"role"`
	IsSuperAdmin     bool                `json:"is_super_admin,omitempty"`
	Designation      string              `json:"designation,omitempty"`
	DepartmentID     string              `json:"department_id,omitempty"`
	ManagerID        string              `json:"manager_id,omitempty"`
	EmployeeCode     string              `json:"employee_code,omitempty"`
	DateOfJoin       string              `json:"date_of_join,omitempty"`
	Status           models.UserStatus   `json:"status"`
	Phone            string              `json:"phone,omitempty"`
	Address          models.Address      `json:"address,omitempty"`
	ProfilePic       string              `json:"profile_pic,omitempty"`
	ResumeURL        string              `json:"resume_url,omitempty"`
	Company          string              `json:"company,omitempty"`
	BankDetails      *models.BankDetails `json:"bank_details,omitempty"`
	LastLogin        primitive.DateTime  `json:"last_login,omitempty"`
	EmailVerified    bool                `json:"email_verified"`
	TwoFactorEnabled bool                `json:"two_factor_enabled"`
	CreatedAt        primitive.DateTime  `json:"created_at,omitempty"`
	UpdatedAt        primitive.DateTime  `json:"updated_at,omitempty"`
}

// CompanyResponse represents company data with encrypted IDs for API responses
type CompanyResponse struct {
	ID         string             `json:"id,omitempty"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`
	Phone      string             `json:"phone,omitempty"`
	Industry   string             `json:"industry,omitempty"`
	Website    string             `json:"website,omitempty"`
	LogoURL    string             `json:"logo_url,omitempty"`
	Address    models.Address     `json:"address,omitempty"`
	OwnerID    string             `json:"owner_id,omitempty"`
	ApprovedBy string             `json:"approved_by,omitempty"`
	IsApproved bool               `json:"is_approved"`
	IsActive   bool               `json:"is_active"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty"`
	UpdatedAt  primitive.DateTime `json:"updated_at,omitempty"`
}

// LoginResponse represents the complete login response with user and company info
type LoginResponse struct {
	Token   string           `json:"token"`
	User    *UserResponse    `json:"user"`
	Company *CompanyResponse `json:"company,omitempty"`
}

// convertUserToResponse converts User model to UserResponse with encrypted IDs
func convertUserToResponse(user *models.User) (*UserResponse, error) {
	if user == nil {
		return nil, errors.New("user is nil")
	}

	response := &UserResponse{
		Username:         user.Username,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Role:             user.Role,
		IsSuperAdmin:     user.IsSuperAdmin,
		Designation:      user.Designation,
		EmployeeCode:     user.EmployeeCode,
		DateOfJoin:       user.DateOfJoin,
		Status:           user.Status,
		Phone:            user.Phone,
		Address:          user.Address,
		ProfilePic:       user.ProfilePic,
		ResumeURL:        user.ResumeURL,
		BankDetails:      user.BankDetails,
		LastLogin:        user.LastLogin,
		EmailVerified:    user.EmailVerified,
		TwoFactorEnabled: user.TwoFactorEnabled,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}

	// Encrypt user ID
	if !user.ID.IsZero() {
		encryptedID, err := encryptions.EncryptID(user.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt user ID: %w", err)
		}
		response.ID = encryptedID
	}

	// Encrypt company ID
	if !user.Company.IsZero() {
		encryptedCompanyID, err := encryptions.EncryptID(user.Company.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.Company = encryptedCompanyID
	}

	// Encrypt department ID
	if !user.DepartmentID.IsZero() {
		encryptedDeptID, err := encryptions.EncryptID(user.DepartmentID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt department ID: %w", err)
		}
		response.DepartmentID = encryptedDeptID
	}

	// Encrypt manager ID
	if !user.ManagerID.IsZero() {
		encryptedManagerID, err := encryptions.EncryptID(user.ManagerID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt manager ID: %w", err)
		}
		response.ManagerID = encryptedManagerID
	}

	return response, nil
}

// convertCompanyToResponse converts Company model to CompanyResponse with encrypted IDs
func convertCompanyToResponse(company *models.Company) (*CompanyResponse, error) {
	if company == nil {
		return nil, errors.New("company is nil")
	}

	response := &CompanyResponse{
		Name:       company.Name,
		Email:      company.Email,
		Phone:      company.Phone,
		Industry:   company.Industry,
		Website:    company.Website,
		LogoURL:    company.LogoURL,
		Address:    company.Address,
		IsApproved: company.IsApproved,
		IsActive:   company.IsActive,
		CreatedAt:  company.CreatedAt,
		UpdatedAt:  company.UpdatedAt,
	}

	// Encrypt company ID
	if !company.ID.IsZero() {
		encryptedID, err := encryptions.EncryptID(company.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt company ID: %w", err)
		}
		response.ID = encryptedID
	}

	// Encrypt owner ID
	if !company.OwnerID.IsZero() {
		encryptedOwnerID, err := encryptions.EncryptID(company.OwnerID.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt owner ID: %w", err)
		}
		response.OwnerID = encryptedOwnerID
	}

	// Encrypt approved by ID
	if !company.ApprovedBy.IsZero() {
		encryptedApprovedBy, err := encryptions.EncryptID(company.ApprovedBy.Hex())
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt approved by ID: %w", err)
		}
		response.ApprovedBy = encryptedApprovedBy
	}

	return response, nil
}

// Signup handles company registration and admin user creation
func (s *AuthService) Signup(req *SignupRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Check if email already exists
	var existingCompany models.Company
	err := companiesCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingCompany)
	if err == nil {
		return errors.New("company with this email already exists")
	}

	// Generate company code (first 3 letters of company name, uppercase)
	companyCode := strings.ToUpper(strings.ReplaceAll(req.CompanyName, " ", ""))
	if len(companyCode) > 4 {
		companyCode = companyCode[:4]
	}

	// Create company
	company := models.Company{
		ID:         primitive.NewObjectID(),
		Name:       req.CompanyName,
		Email:      req.Email,
		Phone:      req.Phone,
		Industry:   req.Industry,
		IsApproved: false,
		IsActive:   true,
	}
	// Company is self-created during signup, so CreatedBy will be the admin user (set after user creation)
	company.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	company.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	company.IsDeleted = false

	_, err = companiesCollection.InsertOne(ctx, company)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	// Generate username for admin (CompanyCode + FirstName initials + LastName initials + Year + 0001)
	now := time.Now()
	username := helpers.GenerateLoginID(companyCode, req.FirstName, req.LastName, now, 1)

	// Hash password
	hashedPassword := encryptions.HashPassword(req.Password)

	// Generate email verification token
	verificationToken, err := helpers.GenerateVerificationToken()
	if err != nil {
		companiesCollection.DeleteOne(ctx, bson.M{"_id": company.ID})
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Create admin user (inactive and unverified until email is confirmed)
	adminUser := models.User{
		ID:                     primitive.NewObjectID(),
		Username:               username,
		Email:                  req.Email,
		Password:               hashedPassword,
		FirstName:              req.FirstName,
		LastName:               req.LastName,
		Role:                   models.RoleAdmin,
		IsSuperAdmin:           false,
		Status:                 models.UserInactive, // Inactive until approved
		EmailVerified:          false,
		EmailVerificationToken: verificationToken,
		TokenExpiry:            primitive.NewDateTimeFromTime(helpers.VerificationTokenExpiry()),
		Company:                company.ID,
		DateOfJoin:             now.Format("2006-01-02"),
		EmployeeCode:           username,
	}
	// User is self-registering, CreatedBy will be self after creation
	adminUser.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	adminUser.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	adminUser.IsDeleted = false
	adminUser.CreatedBy = adminUser.ID
	adminUser.UpdatedBy = adminUser.ID

	_, err = usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		// Rollback: delete the company if user creation fails
		companiesCollection.DeleteOne(ctx, bson.M{"_id": company.ID})
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// Send verification email (non-blocking, log error if fails)
	go func() {
		if err := helpers.SendVerificationEmail(req.Email, verificationToken, req.CompanyName); err != nil {
			fmt.Printf("Failed to send verification email to %s: %v\n", req.Email, err)
		}
	}()

	return nil
}

// Login authenticates user and generates JWT token
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)

	// Find user by username or email (exclude soft-deleted users)
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"username": req.Username},
					{"email": req.Username},
				},
			},
			{
				"$or": []bson.M{
					{"is_deleted": false},
					{"is_deleted": bson.M{"$exists": false}},
				},
			},
		},
	}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if !encryptions.ComparePassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Check if email is verified (except for SuperAdmin)
	if !user.IsSuperAdmin && !user.EmailVerified {
		return nil, errors.New("please verify your email address before logging in")
	}

	// Check if user is active
	if user.Status != models.UserActive {
		return nil, errors.New("user account is inactive")
	}

	var company models.Company
	var companyResponse *CompanyResponse

	// If user belongs to a company, check if company is approved and active
	if !user.IsSuperAdmin {
		err = companiesCollection.FindOne(ctx, helpers.AddNotDeletedFilter(bson.M{"_id": user.Company})).Decode(&company)
		if err != nil {
			return nil, errors.New("company not found")
		}

		if !company.IsApproved {
			return nil, errors.New("company registration is pending approval")
		}

		if !company.IsActive {
			return nil, errors.New("company account is inactive")
		}

		// Convert company to response with encrypted IDs
		companyResponse, err = convertCompanyToResponse(&company)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare company response: %w", err)
		}
	}

	// Generate JWT token with encrypted IDs
	expireTime := time.Now().Add(helpers.JWTExpireDuration)

	// Encrypt user ID for JWT
	encryptedUserID, err := encryptions.EncryptID(user.ID.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt user ID for JWT: %w", err)
	}

	// Encrypt company ID for JWT
	encryptedCompanyID, err := encryptions.EncryptID(user.Company.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt company ID for JWT: %w", err)
	}

	payload := map[string]any{
		"id":             encryptedUserID,
		"username":       user.Username,
		"role":           user.Role,
		"company":        encryptedCompanyID,
		"is_super_admin": user.IsSuperAdmin,
	}

	token, err := helpers.GenerateJWT(payload, expireTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"last_login": primitive.NewDateTimeFromTime(time.Now())}},
	)

	// Remove password from response
	user.Password = ""

	// Convert user to response with encrypted IDs
	userResponse, err := convertUserToResponse(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare user response: %w", err)
	}

	// Prepare login response
	loginResponse := &LoginResponse{
		Token:   token,
		User:    userResponse,
		Company: companyResponse,
	}

	return loginResponse, nil
}

// GetUserProfile retrieves user profile by ID
func (s *AuthService) GetUserProfile(userID primitive.ObjectID) (*UserResponse, error) {
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

	// Convert user to response with encrypted IDs
	userResponse, err := convertUserToResponse(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare user response: %w", err)
	}

	return userResponse, nil
}

// VerifyEmail verifies a user's email using the verification token
func (s *AuthService) VerifyEmail(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)

	// Find user by verification token
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{
		"email_verification_token": token,
	}).Decode(&user)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	// Check if token has expired
	if time.Now().After(user.TokenExpiry.Time()) {
		return errors.New("verification token has expired")
	}

	// Check if email is already verified
	if user.EmailVerified {
		return errors.New("email is already verified")
	}

	// Update user to mark email as verified
	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"email_verified":       true,
				"timestamp.updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
			"$unset": bson.M{
				"email_verification_token": "",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	// Get company details for welcome email
	var company models.Company
	companiesCollection.FindOne(ctx, bson.M{"_id": user.Company}).Decode(&company)

	// Send welcome email (non-blocking)
	go func() {
		if err := helpers.SendWelcomeEmail(user.Email, user.FirstName, company.Name, user.Username); err != nil {
			fmt.Printf("Failed to send welcome email to %s: %v\n", user.Email, err)
		}
	}()

	return nil
}

// ResendVerificationEmail resends the verification email
func (s *AuthService) ResendVerificationEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)

	// Find user by email
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if already verified
	if user.EmailVerified {
		return errors.New("email is already verified")
	}

	// Generate new verification token
	verificationToken, err := helpers.GenerateVerificationToken()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Update user with new token
	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{
			"$set": bson.M{
				"email_verification_token": verificationToken,
				"token_expiry":             primitive.NewDateTimeFromTime(helpers.VerificationTokenExpiry()),
				"timestamp.updated_at":     primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update verification token: %w", err)
	}

	// Get company details
	var company models.Company
	companiesCollection.FindOne(ctx, bson.M{"_id": user.Company}).Decode(&company)

	// Send verification email (non-blocking)
	go func() {
		if err := helpers.SendVerificationEmail(user.Email, verificationToken, company.Name); err != nil {
			fmt.Printf("Failed to send verification email to %s: %v\n", user.Email, err)
		}
	}()

	return nil
}

// ChangePassword updates user password
func (s *AuthService) ChangePassword(userID primitive.ObjectID, req *ChangePasswordRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Find user
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !encryptions.ComparePassword(req.OldPassword, user.Password) {
		return errors.New("old password is incorrect")
	}

	// Hash new password
	hashedPassword := encryptions.HashPassword(req.NewPassword)

	// Update password
	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"password":             hashedPassword,
				"timestamp.updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// ResetPasswordByAdmin allows admin to reset employee password
func (s *AuthService) ResetPasswordByAdmin(adminID, targetUserID primitive.ObjectID, newPassword string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)

	// Verify admin has permission (handled by middleware, but double-check)
	var admin models.User
	err := usersCollection.FindOne(ctx, bson.M{"_id": adminID}).Decode(&admin)
	if err != nil {
		return "", errors.New("admin not found")
	}

	if admin.Role != models.RoleAdmin && !admin.IsSuperAdmin {
		return "", errors.New("only admins can reset passwords")
	}

	// Hash new password
	hashedPassword := encryptions.HashPassword(newPassword)

	// Update target user password
	_, err = usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": targetUserID},
		bson.M{
			"$set": bson.M{
				"password":             hashedPassword,
				"timestamp.updated_at": primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to reset password: %w", err)
	}

	return newPassword, nil
}
