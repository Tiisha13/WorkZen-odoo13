package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/helpers"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		service: services.NewUserService(),
	}
}

// CreateUserRequestAPI represents the request structure for creating a user via API with encrypted IDs
type CreateUserRequestAPI struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone"`
	Role         string `json:"role" validate:"required"`
	Designation  string `json:"designation"`
	DepartmentID string `json:"department_id"`
	ManagerID    string `json:"manager_id"`
	DateOfJoin   string `json:"date_of_join"`
}

// CreateUser creates a new employee
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var reqAPI CreateUserRequestAPI
	if err := c.BodyParser(&reqAPI); err != nil {
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

	// Convert API request to service request
	req := services.CreateUserRequest{
		FirstName:   reqAPI.FirstName,
		LastName:    reqAPI.LastName,
		Email:       reqAPI.Email,
		Phone:       reqAPI.Phone,
		Role:        models.Role(reqAPI.Role),
		Designation: reqAPI.Designation,
		DateOfJoin:  reqAPI.DateOfJoin,
	}

	// Decrypt department_id if provided
	if reqAPI.DepartmentID != "" {
		deptID, err := helpers.DecryptObjectID(reqAPI.DepartmentID)
		if err != nil {
			return constants.HTTPErrors.BadRequest(c, "Invalid department_id")
		}
		req.DepartmentID = &deptID
	}

	// Decrypt manager_id if provided
	if reqAPI.ManagerID != "" {
		mgrID, err := helpers.DecryptObjectID(reqAPI.ManagerID)
		if err != nil {
			return constants.HTTPErrors.BadRequest(c, "Invalid manager_id")
		}
		req.ManagerID = &mgrID
	}

	user, password, err := uc.service.CreateUser(&req, companyID, authUserID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "User created successfully", fiber.Map{
		"user":     user,
		"password": password,
	})
}

// ListUsers retrieves all users in company (excluding own account)
func (uc *UserController) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	// Get authenticated user ID to exclude from results
	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	filters := make(map[string]interface{})
	// Exclude the authenticated user's own account
	filters["_id"] = map[string]interface{}{"$ne": authUserID}

	users, total, err := uc.service.ListUsers(companyID, filters, page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	if total == 0 {
		return constants.HTTPSuccess.OKWithoutData(c, "No users found")
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Users retrieved successfully", users, page, limit, total)
}

// GetUserByID retrieves a user by ID
func (uc *UserController) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	user, err := uc.service.GetUserByID(userID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "User retrieved successfully", user)
}

// UpdateUserRequestAPI for API with encrypted IDs
type UpdateUserRequestAPI struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Role         string `json:"role"`
	Designation  string `json:"designation"`
	DepartmentID string `json:"department_id"`
	Password     string `json:"password"`
}

// UpdateUser updates user details
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	var reqAPI UpdateUserRequestAPI
	if err := c.BodyParser(&reqAPI); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	// Convert API request to service request
	req := services.UpdateUserRequest{
		FirstName:   reqAPI.FirstName,
		LastName:    reqAPI.LastName,
		Email:       reqAPI.Email,
		Phone:       reqAPI.Phone,
		Role:        models.Role(reqAPI.Role),
		Designation: reqAPI.Designation,
		Password:    reqAPI.Password,
	}

	// Decrypt department_id if provided
	if reqAPI.DepartmentID != "" {
		deptID, err := helpers.DecryptObjectID(reqAPI.DepartmentID)
		if err != nil {
			return constants.HTTPErrors.BadRequest(c, "Invalid department_id")
		}
		req.DepartmentID = &deptID
	}

	user, err := uc.service.UpdateUser(userID, authUserID, &req)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "User updated successfully", user)
}

// UpdateUserStatus updates user status
func (uc *UserController) UpdateUserStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err = uc.service.UpdateUserStatus(userID, authUserID, models.UserStatus(req.Status))
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "User status updated successfully")
}

// UpdateBankDetails updates user bank details
func (uc *UserController) UpdateBankDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	var req services.UpdateBankDetailsRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err = uc.service.UpdateBankDetails(userID, authUserID, &req)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Bank details updated successfully")
}

// DeleteUser soft deletes a user
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := helpers.DecryptObjectID(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	authUserID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	err = uc.service.DeleteUser(userID, authUserID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "User deleted successfully")
}
