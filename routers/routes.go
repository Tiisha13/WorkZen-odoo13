package routers

import (
	"api.workzen.odoo/controllers"
	"api.workzen.odoo/middlewares"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(app *fiber.App) {
	// Initialize Controllers
	authController := controllers.NewAuthController()
	// companyController := controllers.NewCompanyController()
	// userController := controllers.NewUserController()
	// attendanceController := controllers.NewAttendanceController()
	// leaveController := controllers.NewLeaveController()
	// salaryController := controllers.NewSalaryController()
	// payrollController := controllers.NewPayrollController()
	// documentController := controllers.NewDocumentController()
	// dashboardController := controllers.NewDashboardController()

	// API v1 Routes
	api := app.Group("/api/v1")

	// Health Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "WorkZen HRMS",
		})
	})

	// ==================== AUTH ROUTES ====================
	auth := api.Group("/auth")
	auth.Post("/signup", authController.Signup)
	auth.Post("/login", authController.Login)
	auth.Get("/me", middlewares.AuthMiddleware(), authController.GetMe)
	auth.Post("/change-password", middlewares.AuthMiddleware(), authController.ChangePassword)

	// ==================== COMPANY ROUTES ====================
	// companies := api.Group("/companies")
	// companies.Use(middlewares.AuthMiddleware())
	// companies.Post("/", middlewares.RequireSuperAdmin(), companyController.CreateCompany)
	// companies.Get("/", middlewares.RequireSuperAdmin(), companyController.ListCompanies)
	// companies.Get("/:id", companyController.GetCompany)
	// companies.Patch("/:id/approve", middlewares.RequireSuperAdmin(), companyController.ApproveCompany)
	// companies.Patch("/:id/deactivate", middlewares.RequireSuperAdmin(), companyController.DeactivateCompany)

	// ==================== USER ROUTES ====================
	// users := api.Group("/users")
	// users.Use(middlewares.AuthMiddleware())
	// users.Post("/", middlewares.RequireHROrAdmin(), userController.CreateUser)
	// users.Get("/", middlewares.RequireHROrAdmin(), userController.ListUsers)
	// users.Get("/:id", userController.GetUser)
	// users.Patch("/:id/status", middlewares.RequireCompanyAdmin(), userController.UpdateUserStatus)
	// users.Patch("/:id/bank", userController.UpdateBankDetails)
	// users.Delete("/:id", middlewares.RequireCompanyAdmin(), userController.DeleteUser)

	// ==================== ATTENDANCE ROUTES ====================
	// attendance := api.Group("/attendance")
	// attendance.Use(middlewares.AuthMiddleware())
	// attendance.Post("/check-in", attendanceController.CheckIn)
	// attendance.Post("/check-out", attendanceController.CheckOut)
	// attendance.Get("/me", attendanceController.GetMyAttendance)
	// attendance.Get("/", middlewares.RequireHROrAdmin(), attendanceController.ListAttendance)

	// ==================== LEAVE ROUTES ====================
	// leaves := api.Group("/leaves")
	// leaves.Use(middlewares.AuthMiddleware())
	// leaves.Post("/", leaveController.ApplyLeave)
	// leaves.Get("/", middlewares.RequireHROrAdmin(), leaveController.ListLeaves)
	// leaves.Patch("/:id/approve", middlewares.RequireHROrAdmin(), leaveController.ApproveLeave)
	// leaves.Patch("/:id/reject", middlewares.RequireHROrAdmin(), leaveController.RejectLeave)

	// ==================== SALARY STRUCTURE ROUTES ====================
	// salary := api.Group("/salary-structure")
	// salary.Use(middlewares.AuthMiddleware())
	// salary.Post("/", middlewares.CanModifySalaryInfo(), salaryController.CreateSalaryStructure)
	// salary.Get("/:employeeId", salaryController.GetSalaryStructure)
	// salary.Patch("/:employeeId", middlewares.CanModifySalaryInfo(), salaryController.UpdateSalaryStructure)

	// ==================== PAYROLL CONFIGURATION ROUTES ====================
	// payrollConfig := api.Group("/payroll/configuration")
	// payrollConfig.Use(middlewares.AuthMiddleware())
	// payrollConfig.Post("/", middlewares.RequireCompanyAdmin(), payrollController.CreateConfiguration)
	// payrollConfig.Get("/", middlewares.RequirePayrollOrAdmin(), payrollController.GetConfiguration)
	// payrollConfig.Patch("/", middlewares.RequireCompanyAdmin(), payrollController.UpdateConfiguration)

	// ==================== PAYROLL & PAYRUN ROUTES ====================
	// payruns := api.Group("/payruns")
	// payruns.Use(middlewares.AuthMiddleware())
	// payruns.Post("/", middlewares.RequirePayrollOrAdmin(), payrollController.CreatePayrun)
	// payruns.Get("/", middlewares.RequirePayrollOrAdmin(), payrollController.ListPayruns)

	// payrolls := api.Group("/payrolls")
	// payrolls.Use(middlewares.AuthMiddleware())
	// payrolls.Get("/:employeeId", payrollController.GetEmployeePayroll)
	// payrolls.Patch("/:id/mark-paid", middlewares.RequirePayrollOrAdmin(), payrollController.MarkAsPaid)

	// ==================== DOCUMENT ROUTES ====================
	// documents := api.Group("/documents")
	// documents.Use(middlewares.AuthMiddleware())
	// documents.Post("/upload", documentController.UploadDocument)
	// documents.Get("/", middlewares.RequireHROrAdmin(), documentController.ListDocuments)
	// documents.Delete("/:id", middlewares.RequireCompanyAdmin(), documentController.DeleteDocument)

	// ==================== DASHBOARD ROUTES ====================
	// dashboard := api.Group("/dashboard")
	// dashboard.Use(middlewares.AuthMiddleware())
	// dashboard.Get("/admin", middlewares.RequireCompanyAdmin(), dashboardController.GetAdminDashboard)
	// dashboard.Get("/superadmin", middlewares.RequireSuperAdmin(), dashboardController.GetSuperAdminDashboard)
}
