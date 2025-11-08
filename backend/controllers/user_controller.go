package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		service: services.NewUserService(),
	}
}

// CreateUser creates a new employee
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var req services.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	user, password, err := uc.service.CreateUser(&req, companyID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "User created successfully", fiber.Map{
		"user":     user,
		"password": password,
	})
}

// ListUsers retrieves all users in company
func (uc *UserController) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	filters := make(map[string]interface{})
	users, total, err := uc.service.ListUsers(companyID, filters, page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Users retrieved successfully", users, page, limit, total)
}

// GetUserByID retrieves a user by ID
func (uc *UserController) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	user, err := uc.service.GetUserByID(userID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "User retrieved successfully", user)
}

// UpdateUserStatus updates user status
func (uc *UserController) UpdateUserStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err = uc.service.UpdateUserStatus(userID, models.UserStatus(req.Status))
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "User status updated successfully")
}

// UpdateBankDetails updates user bank details
func (uc *UserController) UpdateBankDetails(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	var req services.UpdateBankDetailsRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err = uc.service.UpdateBankDetails(userID, &req)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Bank details updated successfully")
}

// DeleteUser soft deletes a user
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid user ID")
	}

	err = uc.service.DeleteUser(userID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "User deleted successfully")
}
