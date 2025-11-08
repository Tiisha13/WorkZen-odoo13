# 🎯 WorkZen HRMS - Quick Completion Guide

## 🔧 Step 1: Fix Existing Code (5 minutes)

### Fix auth_service.go

Search and replace in `/services/auth_service.go`:

```go
// Find:
databases.MongoDatabase.Collection

// Replace with:
databases.MongoDBDatabase.Collection

// Find:
company.CreatedAt = time.Now()
company.UpdatedAt = time.Now()
adminUser.CreatedAt = time.Now()
adminUser.UpdatedAt = time.Now()

// Replace with:
company.CreatedAt = helpers.NowDateTime()
company.UpdatedAt = helpers.NowDateTime()
adminUser.CreatedAt = helpers.NowDateTime()
adminUser.UpdatedAt = helpers.NowDateTime()
```

### Fix company_service.go

Same replacements:

- `databases.MongoDatabase` → `databases.MongoDBDatabase`
- `time.Now()` → `helpers.NowDateTime()` for CreatedAt/UpdatedAt

---

## 🚀 Step 2: Create Remaining Services

### Template for Each Service

```go
package services

import (
	"context"
	"errors"
	"time"
	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type XService struct{}

func NewXService() *XService {
	return &XService{}
}

// Request/Response types
type XRequest struct {
	Field1 string `json:"field1" validate:"required"`
	Field2 int    `json:"field2"`
}

// Methods
func (s *XService) Create(req *XRequest) (*models.X, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := databases.MongoDBDatabase.Collection(collections.Xs)

	item := models.X{
		ID: primitive.NewObjectID(),
		Field1: req.Field1,
		Field2: req.Field2,
	}
	item.CreatedAt = helpers.NowDateTime()
	item.UpdatedAt = helpers.NowDateTime()

	_, err := collection.InsertOne(ctx, item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
```

---

## 📋 Quick Service Implementations

### User Service (`services/user_service.go`)

Key methods:

- `CreateUser(req *CreateUserRequest, companyID primitive.ObjectID) (*models.User, error)`
  - Generate username using `helpers.GenerateLoginID()`
  - Generate random password
  - Hash password
  - Assign to company
- `ListUsers(companyID primitive.ObjectID, filters map[string]interface{}, page, limit int64) ([]models.User, int64, error)`
  - Add companyID to filters
  - Paginate results
- `UpdateBankDetails(userID primitive.ObjectID, bankDetails *models.BankDetails) error`
  - Update user.BankDetails
- `UpdateStatus(userID primitive.ObjectID, status models.UserStatus) error`

---

### Attendance Service (`services/attendance_service.go`)

Key methods:

- `CheckIn(employeeID, companyID primitive.ObjectID) (*models.Attendance, error)`
  - Create attendance record with status=present
  - Check for existing check-in today
  - Store current time as check_in
- `CheckOut(employeeID primitive.ObjectID) error`
  - Find today's attendance
  - Update check_out time
  - Calculate work_hours using `helpers.CalculateWorkHours()`
- `GetMonthlyAttendance(employeeID primitive.ObjectID, month string) ([]models.Attendance, error)`
  - Filter by employeeID and month (YYYY-MM format)

---

### Leave Service (`services/leave_service.go`)

Key methods:

- `ApplyLeave(req *LeaveRequest, employeeID, companyID primitive.ObjectID) (*models.Leave, error)`
  - Create leave with status=pending
  - Calculate days between start and end date
- `ApproveLeave(leaveID, approvedByID primitive.ObjectID) error`
  - Update leave status to approved
  - Create attendance records with status=on_leave for each day
- `RejectLeave(leaveID, rejectedByID primitive.ObjectID) error`
  - Update status to rejected

---

### Salary Service (`services/salary_service.go`)

Key methods:

- `CreateSalaryStructure(employeeID, companyID primitive.ObjectID, monthlyWage float64) (*models.SalaryStructure, error)`
  - Get PayrollConfiguration for company
  - Use `helpers.CalculateSalaryComponents(monthlyWage, config)`
  - Save to database
- `GetSalaryStructure(employeeID primitive.ObjectID) (*models.SalaryStructure, error)`
  - Find active salary structure
- `UpdateSalaryStructure(employeeID primitive.ObjectID, newWage float64) error`
  - Deactivate old structure
  - Create new structure with updated wage
  - Use `helpers.RecalculateStructure()`

---

### Payroll Service (`services/payroll_service.go`)

Key methods:

- `CreatePayrun(companyID, generatedByID primitive.ObjectID, month string) (*models.Payrun, error)`
  - Get all active employees in company
  - For each employee:
    - Get salary structure
    - Calculate PF and tax using `helpers.CalculateDeductions()`
    - Check if has bank account
    - Check if has manager
    - Create payroll record
  - Store warning counts in payrun
- `ListPayruns(companyID primitive.ObjectID, page, limit int64) ([]models.Payrun, int64, error)`
- `GetEmployeePayroll(employeeID primitive.ObjectID, month string) (*models.Payroll, error)`
- `MarkAsPaid(payrollID primitive.ObjectID) error`
  - Update status to paid
  - Set paid_at timestamp

---

## 🎮 Step 3: Create Controllers

### Template for Each Controller

```go
package controllers

import (
	"api.workzen.odoo/constants"
	"api.workzen.odoo/middlewares"
	"api.workzen.odoo/services"
	"github.com/gofiber/fiber/v2"
)

type XController struct {
	service *services.XService
}

func NewXController() *XController {
	return &XController{
		service: services.NewXService(),
	}
}

func (ctrl *XController) Create(c *fiber.Ctx) error {
	var req services.XRequest
	if err := c.BodyParser(&req); err != nil {
		return constants.HTTPErrors.BadRequest(c, "Invalid request body")
	}

	userID, _ := middlewares.GetAuthUserID(c)

	result, err := ctrl.service.Create(&req, userID)
	if err != nil {
		return constants.HTTPErrors.BadRequest(c, err.Error())
	}

	return constants.HTTPSuccess.Created(c, result, "Created successfully")
}

func (ctrl *XController) List(c *fiber.Ctx) error {
	// Parse query params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	companyID, _ := middlewares.GetAuthCompanyID(c)

	items, total, err := ctrl.service.List(companyID, int64(page), int64(limit))
	if err != nil {
		return constants.HTTPErrors.InternalServerError(c, err.Error())
	}

	return constants.HTTPSuccess.OK(c, fiber.Map{
		"items": items,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}, "Retrieved successfully")
}
```

---

## 🛣️ Step 4: Activate Routes

In `/routers/routes.go`, uncomment all route groups and initialize controllers:

```go
// Initialize ALL Controllers
authController := controllers.NewAuthController()
companyController := controllers.NewCompanyController()
userController := controllers.NewUserController()
attendanceController := controllers.NewAttendanceController()
leaveController := controllers.NewLeaveController()
salaryController := controllers.NewSalaryController()
payrollController := controllers.NewPayrollController()
documentController := controllers.NewDocumentController()
dashboardController := controllers.NewDashboardController()

// Then uncomment all the route groups below
```

---

## ⚡ Priority Order

Build in this order for fastest results:

1. **User Service & Controller** (Most critical)
2. **Attendance Service & Controller** (Daily operations)
3. **Salary Service & Controller** (Core feature)
4. **Payroll Service & Controller** (Core feature)
5. **Leave Service & Controller** (HR workflow)
6. **Company Controller** (Already has service)
7. **Document Service & Controller** (File uploads)
8. **Dashboard Service & Controller** (Analytics)

---

## 🧪 Testing Each Module

```bash
# 1. Signup
curl -X POST http://localhost:3000/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Test Corp",
    "email": "admin@test.com",
    "phone": "1234567890",
    "industry": "IT",
    "first_name": "John",
    "last_name": "Doe",
    "password": "password123"
  }'

# 2. Login (after SuperAdmin approves)
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "TESTJODO20250001",
    "password": "password123"
  }'

# 3. Get Profile (use token from login)
curl -X GET http://localhost:3000/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 📦 File Checklist

Create these files:

### Services (7 files)

- [✅] `services/auth_service.go` (fix needed)
- [✅] `services/company_service.go` (fix needed)
- [ ] `services/user_service.go`
- [ ] `services/attendance_service.go`
- [ ] `services/leave_service.go`
- [ ] `services/salary_service.go`
- [ ] `services/payroll_service.go`
- [ ] `services/document_service.go`
- [ ] `services/dashboard_service.go`

### Controllers (9 files)

- [✅] `controllers/auth_controller.go`
- [ ] `controllers/company_controller.go`
- [ ] `controllers/user_controller.go`
- [ ] `controllers/attendance_controller.go`
- [ ] `controllers/leave_controller.go`
- [ ] `controllers/salary_controller.go`
- [ ] `controllers/payroll_controller.go`
- [ ] `controllers/document_controller.go`
- [ ] `controllers/dashboard_controller.go`

---

## 🎁 Bonus: Dashboard Aggregation Example

```go
// GetAdminDashboard in dashboard_service.go
func (s *DashboardService) GetAdminDashboard(companyID primitive.ObjectID) (fiber.Map, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Total employees
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	totalEmployees, _ := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status": models.UserActive,
	})

	// Attendance today
	today := helpers.FormatDate(time.Now())
	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)
	presentToday, _ := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date": today,
		"status": models.StatusPresent,
	})

	// Pending leaves
	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)
	pendingLeaves, _ := leavesCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status": models.LeavePending,
	})

	// Missing bank accounts
	missingBank, _ := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status": models.UserActive,
		"bank_details": bson.M{"$exists": false},
	})

	return fiber.Map{
		"total_employees": totalEmployees,
		"present_today": presentToday,
		"pending_leaves": pendingLeaves,
		"missing_bank_accounts": missingBank,
	}, nil
}
```

---

## ✅ Final Checklist

Before deployment:

- [ ] Fix auth_service.go database references
- [ ] Fix company_service.go database references
- [ ] Create all remaining services
- [ ] Create all remaining controllers
- [ ] Uncomment all routes
- [ ] Test signup flow
- [ ] Test login flow
- [ ] Test authenticated endpoints
- [ ] Test RBAC permissions
- [ ] Add error logging
- [ ] Add request validation
- [ ] Set up MongoDB indexes
- [ ] Configure CORS properly
- [ ] Set up environment variables
- [ ] Add API documentation
- [ ] Write unit tests

---

## 🎯 Estimated Time to Complete

| Task                            | Time         |
| ------------------------------- | ------------ |
| Fix existing code               | 10 min       |
| User service + controller       | 45 min       |
| Attendance service + controller | 30 min       |
| Leave service + controller      | 30 min       |
| Salary service + controller     | 45 min       |
| Payroll service + controller    | 1 hour       |
| Document service + controller   | 45 min       |
| Dashboard service + controller  | 1 hour       |
| Testing & debugging             | 2 hours      |
| **TOTAL**                       | **~7 hours** |

---

**You're 65% done! The foundation is rock solid. Keep going! 🚀**
