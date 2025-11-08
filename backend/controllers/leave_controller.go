package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveController struct {
	service *services.LeaveService
}

func NewLeaveController() *LeaveController {
	return &LeaveController{
		service: services.NewLeaveService(),
	}
}

// ApplyLeave handles leave application by employee
func (lc *LeaveController) ApplyLeave(c *fiber.Ctx) error {
	var req services.ApplyLeaveRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	employeeID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	leave, err := lc.service.ApplyLeave(&req, employeeID, companyID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Leave application submitted successfully", leave)
}

// ListLeaves retrieves leave records with filters
func (lc *LeaveController) ListLeaves(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	// Optional filters
	filters := make(map[string]interface{})
	if employeeIDStr := c.Query("employee_id"); employeeIDStr != "" {
		filters["employee_id"] = employeeIDStr
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	leaves, total, err := lc.service.ListLeaves(companyID, filters, page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Leaves retrieved successfully", fiber.Map{
		"leaves": leaves,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// ApproveLeave approves a leave request
func (lc *LeaveController) ApproveLeave(c *fiber.Ctx) error {
	id := c.Params("id")
	leaveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid leave ID")
	}

	approverID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	err = lc.service.ApproveLeave(leaveID, approverID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Leave approved successfully")
}

// RejectLeave rejects a leave request
func (lc *LeaveController) RejectLeave(c *fiber.Ctx) error {
	id := c.Params("id")
	leaveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid leave ID")
	}

	approverID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	err = lc.service.RejectLeave(leaveID, approverID)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Leave rejected successfully")
}
