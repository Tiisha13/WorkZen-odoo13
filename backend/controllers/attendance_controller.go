package controllers

import (
	"strconv"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
)

type AttendanceController struct {
	service *services.AttendanceService
}

func NewAttendanceController() *AttendanceController {
	return &AttendanceController{
		service: services.NewAttendanceService(),
	}
}

// CheckIn handles employee check-in
func (ac *AttendanceController) CheckIn(c *fiber.Ctx) error {
	userID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	attendance, err := ac.service.CheckIn(userID, companyID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, "Check-in successful", attendance)
}

// CheckOut handles employee check-out
func (ac *AttendanceController) CheckOut(c *fiber.Ctx) error {
	userID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	err = ac.service.CheckOut(userID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, err.Error())
	}

	return constants.HTTPSuccess.OKWithoutData(c, "Check-out successful")
}

// GetMyAttendance retrieves attendance for logged-in user
func (ac *AttendanceController) GetMyAttendance(c *fiber.Ctx) error {
	userID, err := middlewares.GetAuthUserID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	month := c.Query("month") // YYYY-MM format

	attendances, err := ac.service.GetMyAttendance(userID, month)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Attendance retrieved successfully", attendances)
}

// ListAttendance retrieves all attendance records (HR/Admin)
func (ac *AttendanceController) ListAttendance(c *fiber.Ctx) error {
	page, _ := strconv.ParseInt(c.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.Query("limit", "10"), 10, 64)

	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	// Optional filters
	filters := make(map[string]interface{})
	if employeeID := c.Query("employee_id"); employeeID != "" {
		filters["employee_id"] = employeeID
	}
	if date := c.Query("date"); date != "" {
		filters["date"] = date
	}

	attendances, total, err := ac.service.ListAttendance(companyID, filters, page, limit)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OkWithPagination(c, "Attendance list retrieved successfully", attendances, page, limit, total)
} // GetAttendanceSummary retrieves attendance summary statistics
func (ac *AttendanceController) GetAttendanceSummary(c *fiber.Ctx) error {
	companyID, err := middlewares.GetAuthCompanyID(c)
	if err != nil {
		return constants.HTTPErrors.Unauthorized(c, err.Error())
	}

	date := c.Query("date") // Optional, defaults to today

	summary, err := ac.service.GetAttendanceSummary(companyID, date)
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, "Attendance summary retrieved successfully", summary)
}
