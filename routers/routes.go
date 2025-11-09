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
	companyController := controllers.NewCompanyController()
	userController := controllers.NewUserController()
	departmentController := controllers.NewDepartmentController()
	attendanceController := controllers.NewAttendanceController()
	leaveController := controllers.NewLeaveController()
	salaryController := controllers.NewSalaryController()
	payrollController := controllers.NewPayrollController()
	documentController := controllers.NewDocumentController()
	dashboardController := controllers.NewDashboardController()

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
	auth.Get("/verify-email", authController.VerifyEmail)
	auth.Post("/resend-verification", authController.ResendVerificationEmail)
	auth.Get("/me", middlewares.AuthMiddleware(), authController.GetMe)
	auth.Post("/change-password", middlewares.AuthMiddleware(), authController.ChangePassword)

	// ==================== COMPANY ROUTES ====================
	companies := api.Group("/companies")
	companies.Use(middlewares.AuthMiddleware())
	companies.Post("/", middlewares.RequireSuperAdmin(), companyController.CreateCompany)
	companies.Get("/", middlewares.RequireSuperAdmin(), companyController.ListCompanies)
	companies.Get("/:id", companyController.GetCompanyByID)
	companies.Patch("/:id/approve", middlewares.RequireSuperAdmin(), companyController.ApproveCompany)
	companies.Patch("/:id/deactivate", middlewares.RequireSuperAdmin(), companyController.DeactivateCompany)

	// ==================== USER ROUTES ====================
	users := api.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	users.Post("/", middlewares.RequireHROrAdmin(), userController.CreateUser)
	users.Get("/", middlewares.RequireHROrAdmin(), userController.ListUsers)
	users.Get("/:id", userController.GetUserByID)
	users.Put("/:id", middlewares.RequireHROrAdmin(), userController.UpdateUser)
	users.Patch("/:id/status", middlewares.RequireCompanyAdmin(), userController.UpdateUserStatus)
	users.Patch("/:id/bank", userController.UpdateBankDetails)
	users.Delete("/:id", middlewares.RequireCompanyAdmin(), userController.DeleteUser)

	// ==================== DEPARTMENT ROUTES ====================
	departments := api.Group("/departments")
	departments.Use(middlewares.AuthMiddleware())
	departments.Post("/", middlewares.RequireHROrAdmin(), departmentController.CreateDepartment)
	departments.Get("/", departmentController.ListDepartments)
	departments.Get("/:id", departmentController.GetDepartmentByID)
	departments.Patch("/:id", middlewares.RequireHROrAdmin(), departmentController.UpdateDepartment)
	departments.Delete("/:id", middlewares.RequireCompanyAdmin(), departmentController.DeleteDepartment)

	// ==================== ATTENDANCE ROUTES ====================
	attendance := api.Group("/attendance")
	attendance.Use(middlewares.AuthMiddleware())
	attendance.Post("/check-in", attendanceController.CheckIn)
	attendance.Post("/check-out", attendanceController.CheckOut)
	attendance.Delete("/reset", attendanceController.ResetAttendance)
	attendance.Get("/me", attendanceController.GetMyAttendance)
	attendance.Get("/", middlewares.RequireHROrAdmin(), attendanceController.ListAttendance)
	attendance.Get("/summary", middlewares.RequireHROrAdmin(), attendanceController.GetAttendanceSummary)

	// ==================== LEAVE ROUTES ====================
	leaves := api.Group("/leaves")
	leaves.Use(middlewares.AuthMiddleware())
	leaves.Post("/", leaveController.ApplyLeave)
	leaves.Get("/", leaveController.ListLeaves) // All users can list (filtered by role in controller)
	leaves.Patch("/:id/approve", middlewares.RequireHROrAdmin(), leaveController.ApproveLeave)
	leaves.Patch("/:id/reject", middlewares.RequireHROrAdmin(), leaveController.RejectLeave)

	// ==================== SALARY STRUCTURE ROUTES ====================
	salary := api.Group("/salary-structure")
	salary.Use(middlewares.AuthMiddleware())
	salary.Post("/", middlewares.CanModifySalaryInfo(), salaryController.CreateSalaryStructure)
	salary.Get("/:employee_id", salaryController.GetSalaryStructure)
	salary.Patch("/:employee_id", middlewares.CanModifySalaryInfo(), salaryController.UpdateSalaryStructure)

	// ==================== PAYROLL CONFIGURATION ROUTES ====================
	payrollConfig := api.Group("/payroll/configuration")
	payrollConfig.Use(middlewares.AuthMiddleware())
	payrollConfig.Post("/", middlewares.RequireCompanyAdmin(), payrollController.CreateConfiguration)
	payrollConfig.Get("/", middlewares.RequirePayrollOrAdmin(), payrollController.GetConfiguration)

	// ==================== PAYROLL & PAYRUN ROUTES ====================
	payruns := api.Group("/payruns")
	payruns.Use(middlewares.AuthMiddleware())
	payruns.Post("/", middlewares.RequirePayrollOrAdmin(), payrollController.CreatePayrun)
	payruns.Get("/", payrollController.ListPayruns) // Allow all authenticated users with role filtering

	payrolls := api.Group("/payrolls")
	payrolls.Use(middlewares.AuthMiddleware())
	payrolls.Get("/:employee_id", payrollController.GetEmployeePayroll)
	payrolls.Patch("/:id/mark-paid", middlewares.RequirePayrollOrAdmin(), payrollController.MarkAsPaid)

	// ==================== DOCUMENT ROUTES ====================
	documents := api.Group("/documents")
	documents.Use(middlewares.AuthMiddleware())
	documents.Post("/", documentController.UploadDocument)              // Upload document
	documents.Get("/", documentController.ListDocuments)                // All users can list (filtered by role in controller)
	documents.Get("/:id/view", documentController.ViewDocument)         // View document (images, videos, PDFs)
	documents.Get("/:id/download", documentController.DownloadDocument) // Download document
	documents.Delete("/:id", middlewares.RequireCompanyAdmin(), documentController.DeleteDocument)

	// ==================== DASHBOARD ROUTES ====================
	dashboard := api.Group("/dashboard")
	dashboard.Use(middlewares.AuthMiddleware())
	dashboard.Get("/", dashboardController.GetDashboard) // General dashboard for all users
	dashboard.Get("/admin", middlewares.RequireCompanyAdmin(), dashboardController.GetAdminDashboard)
	dashboard.Get("/superadmin", middlewares.RequireSuperAdmin(), dashboardController.GetSuperAdminDashboard)
}
