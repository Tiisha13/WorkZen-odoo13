package controllers

import (
	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
)

type DashboardController struct {
	service *services.DashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		service: services.NewDashboardService(),
	}
}

// GetDashboard retrieves dashboard statistics for any authenticated user
func (dc *DashboardController) GetDashboard(c *fiber.Ctx) error {
	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	user, err := middlewares.GetAuthUser(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	stats, err := dc.service.GetAdminDashboard(companyID, user.Role)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Dashboard statistics retrieved successfully", stats)
}

// GetAdminDashboard retrieves dashboard statistics for company admins
func (dc *DashboardController) GetAdminDashboard(c *fiber.Ctx) error {
	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	user, err := middlewares.GetAuthUser(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	stats, err := dc.service.GetAdminDashboard(companyID, user.Role)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Dashboard statistics retrieved successfully", stats)
}

// GetSuperAdminDashboard retrieves platform-wide statistics for super admin
func (dc *DashboardController) GetSuperAdminDashboard(c *fiber.Ctx) error {
	stats, err := dc.service.GetSuperAdminDashboard()
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Platform statistics retrieved successfully", stats)
}
