// Package controllers provides HTTP request handlers for the WorkZen API.
package controllers

import (
	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Signup handles POST /api/v1/auth/signup
func (ctrl *AuthController) Signup(c *fiber.Ctx) error {
	var req services.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err := ctrl.authService.Signup(&req)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, err.Error())
	}

	return constants.HTTPSuccess.CreatedWithoutData(c, "Signup successful, awaiting approval")
}

// Login handles POST /api/v1/auth/login
func (ctrl *AuthController) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	loginResponse, err := ctrl.authService.Login(&req)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Login successful", fiber.Map{
		"token":   loginResponse.Token,
		"user":    loginResponse.User,
		"company": loginResponse.Company,
	})
}

// GetMe handles GET /api/v1/auth/me
func (ctrl *AuthController) GetMe(c *fiber.Ctx) error {
	userID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	user, err := ctrl.authService.GetUserProfile(userID)
	if err != nil {
		return constants.HTTPErrors.NotFound(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "User profile retrieved successfully", user)
}

// ChangePassword handles POST /api/v1/auth/change-password
func (ctrl *AuthController) ChangePassword(c *fiber.Ctx) error {
	userID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	var req services.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	err = ctrl.authService.ChangePassword(userID, &req)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Password changed successfully")
}
