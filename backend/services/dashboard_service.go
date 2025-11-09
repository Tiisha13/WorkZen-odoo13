package services

import (
	"context"
	"fmt"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// DepartmentStats represents department-wise statistics
type DepartmentStats struct {
	Name    string `json:"name"`
	Count   int64  `json:"count"`
	Present int64  `json:"present"`
	Absent  int64  `json:"absent"`
	OnLeave int64  `json:"on_leave"`
}

// MonthlyAttendance represents monthly attendance trend
type MonthlyAttendance struct {
	Month   string `json:"month"`
	Present int64  `json:"present"`
	Absent  int64  `json:"absent"`
	OnLeave int64  `json:"on_leave"`
}

// LeaveTypeStats represents leave statistics by type
type LeaveTypeStats struct {
	Type     string `json:"type"`
	Pending  int64  `json:"pending"`
	Approved int64  `json:"approved"`
	Rejected int64  `json:"rejected"`
}

// AdminDashboardStats for company admins
type AdminDashboardStats struct {
	TotalEmployees       int64               `json:"total_employees"`
	ActiveEmployees      int64               `json:"active_employees"`
	InactiveEmployees    int64               `json:"inactive_employees"`
	PresentToday         int64               `json:"present_today"`
	AbsentToday          int64               `json:"absent_today"`
	OnLeaveToday         int64               `json:"on_leave_today"`
	PendingLeaves        int64               `json:"pending_leaves"`
	ApprovedLeaves       int64               `json:"approved_leaves"`
	RejectedLeaves       int64               `json:"rejected_leaves"`
	MissingBankAccounts  int64               `json:"missing_bank_accounts"`
	MissingManagers      int64               `json:"missing_managers"`
	TotalPayrollThisYear int64               `json:"total_payroll_this_year"`
	DepartmentStats      []DepartmentStats   `json:"department_stats"`
	MonthlyAttendance    []MonthlyAttendance `json:"monthly_attendance"`
	LeaveTypeStats       []LeaveTypeStats    `json:"leave_type_stats"`
	TotalDepartments     int64               `json:"total_departments"`
	AttendanceRate       float64             `json:"attendance_rate"`
}

// GetAdminDashboard retrieves dashboard stats for company admins
func (s *DashboardService) GetAdminDashboard(companyID primitive.ObjectID, userRole models.Role) (*AdminDashboardStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &AdminDashboardStats{}

	// Collections
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)
	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)
	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)

	// Define role hierarchy - who can see whom
	// Admin should not be counted as employee, they are managers
	var excludeRoles []models.Role
	switch userRole {
	case models.RoleAdmin:
		// Admin can see everyone except superadmin and other admins
		excludeRoles = []models.Role{models.RoleSuperAdmin, models.RoleAdmin}
	case models.RoleHR:
		// HR can see HR, Payroll, and Employees (not Admin or SuperAdmin)
		excludeRoles = []models.Role{models.RoleSuperAdmin, models.RoleAdmin}
	case models.RolePayroll:
		// Payroll can see Payroll and Employees (not Admin, HR, or SuperAdmin)
		excludeRoles = []models.Role{models.RoleSuperAdmin, models.RoleAdmin, models.RoleHR}
	case models.RoleEmployee:
		// Employees can only see themselves
		excludeRoles = []models.Role{models.RoleSuperAdmin, models.RoleAdmin, models.RoleHR, models.RolePayroll}
	default:
		// Default: exclude superadmin and admin
		excludeRoles = []models.Role{models.RoleSuperAdmin, models.RoleAdmin}
	}

	// Total employees (excluding admins and higher roles)
	total, err := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"role":    bson.M{"$nin": excludeRoles},
	})
	if err == nil {
		stats.TotalEmployees = total
	}

	// Active employees
	active, err := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.UserActive,
		"role":    bson.M{"$nin": excludeRoles},
	})
	if err == nil {
		stats.ActiveEmployees = active
	}

	// Inactive employees
	inactive, err := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.UserInactive,
		"role":    bson.M{"$nin": excludeRoles},
	})
	if err == nil {
		stats.InactiveEmployees = inactive
	}

	// Today's date for attendance
	today := time.Now().Format("2006-01-02")

	// Get all employees to filter attendance (exclude admin and higher roles)
	var employees []models.User
	empCursor, err := usersCollection.Find(ctx, bson.M{
		"company": companyID,
		"role":    bson.M{"$nin": excludeRoles},
	})
	if err == nil {
		defer empCursor.Close(ctx)
		empCursor.All(ctx, &employees)
	}

	// Build employee IDs list
	employeeIDs := make([]primitive.ObjectID, 0, len(employees))
	for _, emp := range employees {
		employeeIDs = append(employeeIDs, emp.ID)
	}

	// Present today (filter by employee IDs)
	present := int64(0)
	if len(employeeIDs) > 0 {
		present, err = attendanceCollection.CountDocuments(ctx, bson.M{
			"company":     companyID,
			"date":        today,
			"status":      models.StatusPresent,
			"employee_id": bson.M{"$in": employeeIDs},
		})
		if err == nil {
			stats.PresentToday = present
		}
	}

	// Absent today
	absent := int64(0)
	if len(employeeIDs) > 0 {
		absent, err = attendanceCollection.CountDocuments(ctx, bson.M{
			"company":     companyID,
			"date":        today,
			"status":      models.StatusAbsent,
			"employee_id": bson.M{"$in": employeeIDs},
		})
		if err == nil {
			stats.AbsentToday = absent
		}
	}

	// On leave today
	onLeave := int64(0)
	if len(employeeIDs) > 0 {
		onLeave, err = attendanceCollection.CountDocuments(ctx, bson.M{
			"company":     companyID,
			"date":        today,
			"status":      models.StatusOnLeave,
			"employee_id": bson.M{"$in": employeeIDs},
		})
		if err == nil {
			stats.OnLeaveToday = onLeave
		}
	}

	// Pending leaves (filter by employee IDs)
	pending := int64(0)
	if len(employeeIDs) > 0 {
		pending, err = leavesCollection.CountDocuments(ctx, bson.M{
			"company":     companyID,
			"status":      models.LeavePending,
			"employee_id": bson.M{"$in": employeeIDs},
		})
		if err == nil {
			stats.PendingLeaves = pending
		}
	}

	// Approved leaves
	approved := int64(0)
	if len(employeeIDs) > 0 {
		approved, err = leavesCollection.CountDocuments(ctx, bson.M{
			"company":     companyID,
			"status":      models.LeaveApproved,
			"employee_id": bson.M{"$in": employeeIDs},
		})
		if err == nil {
			stats.ApprovedLeaves = approved
		}
	}

	// Rejected leaves
	rejected := int64(0)
	if len(employeeIDs) > 0 {
		rejected, err = leavesCollection.CountDocuments(ctx, bson.M{
			"company":     companyID,
			"status":      models.LeaveRejected,
			"employee_id": bson.M{"$in": employeeIDs},
		})
		if err == nil {
			stats.RejectedLeaves = rejected
		}
	}

	// Missing bank accounts
	missingBank, err := usersCollection.CountDocuments(ctx, bson.M{
		"company":      companyID,
		"status":       models.UserActive,
		"bank_details": bson.M{"$exists": false},
		"role":         bson.M{"$nin": excludeRoles},
	})
	if err == nil {
		stats.MissingBankAccounts = missingBank
	}

	// Missing managers
	missingManager, err := usersCollection.CountDocuments(ctx, bson.M{
		"company":    companyID,
		"status":     models.UserActive,
		"manager_id": primitive.NilObjectID,
		"role":       bson.M{"$nin": excludeRoles},
	})
	if err == nil {
		stats.MissingManagers = missingManager
	}

	// Total payroll this year
	currentYear := time.Now().Format("2006")
	monthRegex := fmt.Sprintf("^%s-", currentYear)

	// Aggregate total payroll
	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"company": companyID,
			"month":   bson.M{"$regex": monthRegex},
		}},
		bson.M{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$net_pay"},
		}},
	}

	cursor, err := payrollCollection.Aggregate(ctx, pipeline)
	if err == nil {
		defer cursor.Close(ctx)
		var results []bson.M
		if err = cursor.All(ctx, &results); err == nil && len(results) > 0 {
			if total, ok := results[0]["total"].(int64); ok {
				stats.TotalPayrollThisYear = total
			}
		}
	}

	// Get department statistics
	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)
	deptCursor, err := departmentsCollection.Find(ctx, bson.M{"company": companyID})
	if err == nil {
		defer deptCursor.Close(ctx)
		var departments []models.Department
		if err = deptCursor.All(ctx, &departments); err == nil {
			stats.TotalDepartments = int64(len(departments))
			stats.DepartmentStats = make([]DepartmentStats, 0, len(departments))

			for _, dept := range departments {
				// Count employees in department (exclude admin and higher roles)
				empCount, _ := usersCollection.CountDocuments(ctx, bson.M{
					"company":       companyID,
					"department_id": dept.ID,
					"role":          bson.M{"$nin": excludeRoles},
				})

				// Get employee IDs in this department to filter attendance
				var deptEmployees []models.User
				deptCursor, _ := usersCollection.Find(ctx, bson.M{
					"company":       companyID,
					"department_id": dept.ID,
					"role":          bson.M{"$nin": excludeRoles},
				})
				if deptCursor != nil {
					deptCursor.All(ctx, &deptEmployees)
					deptCursor.Close(ctx)
				}

				// Get employee IDs
				deptEmployeeIDs := make([]primitive.ObjectID, 0, len(deptEmployees))
				for _, emp := range deptEmployees {
					deptEmployeeIDs = append(deptEmployeeIDs, emp.ID)
				}

				// Count present today in department
				presentCount := int64(0)
				if len(deptEmployeeIDs) > 0 {
					presentCount, _ = attendanceCollection.CountDocuments(ctx, bson.M{
						"company":     companyID,
						"date":        today,
						"status":      models.StatusPresent,
						"employee_id": bson.M{"$in": deptEmployeeIDs},
					})
				}

				// Count absent today in department
				absentCount := int64(0)
				if len(deptEmployeeIDs) > 0 {
					absentCount, _ = attendanceCollection.CountDocuments(ctx, bson.M{
						"company":     companyID,
						"date":        today,
						"status":      models.StatusAbsent,
						"employee_id": bson.M{"$in": deptEmployeeIDs},
					})
				}

				// Count on leave today in department
				leaveCount := int64(0)
				if len(deptEmployeeIDs) > 0 {
					leaveCount, _ = attendanceCollection.CountDocuments(ctx, bson.M{
						"company":     companyID,
						"date":        today,
						"status":      models.StatusOnLeave,
						"employee_id": bson.M{"$in": deptEmployeeIDs},
					})
				}

				stats.DepartmentStats = append(stats.DepartmentStats, DepartmentStats{
					Name:    dept.Name,
					Count:   empCount,
					Present: presentCount,
					Absent:  absentCount,
					OnLeave: leaveCount,
				})
			}
		}
	}

	// Get monthly attendance for last 6 months
	stats.MonthlyAttendance = make([]MonthlyAttendance, 0, 6)
	for i := 5; i >= 0; i-- {
		monthDate := time.Now().AddDate(0, -i, 0)
		monthStr := monthDate.Format("2006-01")
		monthName := monthDate.Format("Jan 2006")

		presentMonth := int64(0)
		absentMonth := int64(0)
		leaveMonth := int64(0)

		if len(employeeIDs) > 0 {
			presentMonth, _ = attendanceCollection.CountDocuments(ctx, bson.M{
				"company":     companyID,
				"date":        bson.M{"$regex": fmt.Sprintf("^%s", monthStr)},
				"status":      models.StatusPresent,
				"employee_id": bson.M{"$in": employeeIDs},
			})

			absentMonth, _ = attendanceCollection.CountDocuments(ctx, bson.M{
				"company":     companyID,
				"date":        bson.M{"$regex": fmt.Sprintf("^%s", monthStr)},
				"status":      models.StatusAbsent,
				"employee_id": bson.M{"$in": employeeIDs},
			})

			leaveMonth, _ = attendanceCollection.CountDocuments(ctx, bson.M{
				"company":     companyID,
				"date":        bson.M{"$regex": fmt.Sprintf("^%s", monthStr)},
				"status":      models.StatusOnLeave,
				"employee_id": bson.M{"$in": employeeIDs},
			})
		}

		stats.MonthlyAttendance = append(stats.MonthlyAttendance, MonthlyAttendance{
			Month:   monthName,
			Present: presentMonth,
			Absent:  absentMonth,
			OnLeave: leaveMonth,
		})
	}

	// Get leave type statistics
	leaveTypes := []models.LeaveType{
		models.LeaveSick,
		models.LeaveCasual,
		models.LeaveVacation,
	}

	stats.LeaveTypeStats = make([]LeaveTypeStats, 0, len(leaveTypes))
	for _, leaveType := range leaveTypes {
		pendingType := int64(0)
		approvedType := int64(0)
		rejectedType := int64(0)

		if len(employeeIDs) > 0 {
			pendingType, _ = leavesCollection.CountDocuments(ctx, bson.M{
				"company":     companyID,
				"leave_type":  leaveType,
				"status":      models.LeavePending,
				"employee_id": bson.M{"$in": employeeIDs},
			})

			approvedType, _ = leavesCollection.CountDocuments(ctx, bson.M{
				"company":     companyID,
				"leave_type":  leaveType,
				"status":      models.LeaveApproved,
				"employee_id": bson.M{"$in": employeeIDs},
			})

			rejectedType, _ = leavesCollection.CountDocuments(ctx, bson.M{
				"company":     companyID,
				"leave_type":  leaveType,
				"status":      models.LeaveRejected,
				"employee_id": bson.M{"$in": employeeIDs},
			})
		}

		stats.LeaveTypeStats = append(stats.LeaveTypeStats, LeaveTypeStats{
			Type:     string(leaveType),
			Pending:  pendingType,
			Approved: approvedType,
			Rejected: rejectedType,
		})
	}

	// Calculate attendance rate
	if stats.ActiveEmployees > 0 {
		stats.AttendanceRate = (float64(stats.PresentToday) / float64(stats.ActiveEmployees)) * 100
	}

	return stats, nil
}

// CompanyStats represents company-wise statistics for SuperAdmin
type CompanyStats struct {
	Name            string `json:"name"`
	TotalEmployees  int64  `json:"count"`
	ActiveEmployees int64  `json:"present"`
	PresentToday    int64  `json:"absent"` // Reusing field names for chart compatibility
	OnLeave         int64  `json:"on_leave"`
}

// SuperAdminDashboardStats for platform-level stats
type SuperAdminDashboardStats struct {
	TotalCompanies        int64               `json:"total_companies"`
	ActiveCompanies       int64               `json:"active_companies"`
	PendingApprovals      int64               `json:"pending_approvals"`
	TotalEmployees        int64               `json:"total_employees"`
	ActiveEmployees       int64               `json:"active_employees"`
	InactiveEmployees     int64               `json:"inactive_employees"`
	TotalPayrollProcessed int64               `json:"total_payroll_processed"`
	TotalPayrunsGenerated int64               `json:"total_payruns_generated"`
	PlatformRevenue       float64             `json:"platform_revenue"`
	PresentToday          int64               `json:"present_today"`
	AbsentToday           int64               `json:"absent_today"`
	OnLeaveToday          int64               `json:"on_leave_today"`
	PendingLeaves         int64               `json:"pending_leaves"`
	ApprovedLeaves        int64               `json:"approved_leaves"`
	RejectedLeaves        int64               `json:"rejected_leaves"`
	TotalDepartments      int64               `json:"total_departments"`
	AttendanceRate        float64             `json:"attendance_rate"`
	DepartmentStats       []CompanyStats      `json:"department_stats"` // Company stats displayed as departments
	MonthlyAttendance     []MonthlyAttendance `json:"monthly_attendance"`
	LeaveTypeStats        []LeaveTypeStats    `json:"leave_type_stats"`
}

// GetSuperAdminDashboard retrieves platform-wide statistics
func (s *DashboardService) GetSuperAdminDashboard() (*SuperAdminDashboardStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &SuperAdminDashboardStats{}

	// Collections
	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)
	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)
	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)
	payrunsCollection := databases.MongoDBDatabase.Collection(collections.Payruns)
	departmentsCollection := databases.MongoDBDatabase.Collection(collections.Departments)

	today := time.Now().Format("2006-01-02")

	// Total companies
	total, err := companiesCollection.CountDocuments(ctx, bson.M{})
	if err == nil {
		stats.TotalCompanies = total
	}

	// Active companies
	active, err := companiesCollection.CountDocuments(ctx, bson.M{
		"is_approved": true,
	})
	if err == nil {
		stats.ActiveCompanies = active
	}

	// Pending approvals
	pending, err := companiesCollection.CountDocuments(ctx, bson.M{
		"is_approved": false,
	})
	if err == nil {
		stats.PendingApprovals = pending
	}

	// Total employees across platform (excluding SuperAdmins and Admins)
	totalEmployees, err := usersCollection.CountDocuments(ctx, bson.M{
		"role": bson.M{"$nin": []models.Role{models.RoleSuperAdmin, models.RoleAdmin}},
	})
	if err == nil {
		stats.TotalEmployees = totalEmployees
	}

	// Active employees
	activeEmployees, err := usersCollection.CountDocuments(ctx, bson.M{
		"role":      bson.M{"$nin": []models.Role{models.RoleSuperAdmin, models.RoleAdmin}},
		"is_active": true,
	})
	if err == nil {
		stats.ActiveEmployees = activeEmployees
	}

	// Inactive employees
	inactiveEmployees, err := usersCollection.CountDocuments(ctx, bson.M{
		"role":      bson.M{"$nin": []models.Role{models.RoleSuperAdmin, models.RoleAdmin}},
		"is_active": false,
	})
	if err == nil {
		stats.InactiveEmployees = inactiveEmployees
	}

	// Total departments
	totalDepts, err := departmentsCollection.CountDocuments(ctx, bson.M{})
	if err == nil {
		stats.TotalDepartments = totalDepts
	}

	// Present today
	presentToday, err := attendanceCollection.CountDocuments(ctx, bson.M{
		"date":   today,
		"status": models.StatusPresent,
	})
	if err == nil {
		stats.PresentToday = presentToday
	}

	// Absent today
	absentToday, err := attendanceCollection.CountDocuments(ctx, bson.M{
		"date":   today,
		"status": models.StatusAbsent,
	})
	if err == nil {
		stats.AbsentToday = absentToday
	}

	// On leave today
	onLeaveToday, err := leavesCollection.CountDocuments(ctx, bson.M{
		"start_date": bson.M{"$lte": today},
		"end_date":   bson.M{"$gte": today},
		"status":     models.LeaveApproved,
	})
	if err == nil {
		stats.OnLeaveToday = onLeaveToday
	}

	// Calculate attendance rate
	totalAttendanceRecordsToday := stats.PresentToday + stats.AbsentToday + stats.OnLeaveToday
	if totalAttendanceRecordsToday > 0 {
		stats.AttendanceRate = (float64(stats.PresentToday) / float64(totalAttendanceRecordsToday)) * 100
	}

	// Pending leaves
	pendingLeaves, err := leavesCollection.CountDocuments(ctx, bson.M{
		"status": models.LeavePending,
	})
	if err == nil {
		stats.PendingLeaves = pendingLeaves
	}

	// Approved leaves
	approvedLeaves, err := leavesCollection.CountDocuments(ctx, bson.M{
		"status": models.LeaveApproved,
	})
	if err == nil {
		stats.ApprovedLeaves = approvedLeaves
	}

	// Rejected leaves
	rejectedLeaves, err := leavesCollection.CountDocuments(ctx, bson.M{
		"status": models.LeaveRejected,
	})
	if err == nil {
		stats.RejectedLeaves = rejectedLeaves
	}

	// Total payroll processed
	processedPayroll, err := payrollCollection.CountDocuments(ctx, bson.M{})
	if err == nil {
		stats.TotalPayrollProcessed = processedPayroll
	}

	// Total payruns generated
	totalPayruns, err := payrunsCollection.CountDocuments(ctx, bson.M{})
	if err == nil {
		stats.TotalPayrunsGenerated = totalPayruns
	}

	// Platform revenue calculation (example: could be based on per-employee fees)
	// For now, just set to 0 - implement based on business model
	stats.PlatformRevenue = 0.0

	// Company-wise statistics (top companies by employee count)
	companyStatsPipeline := []bson.M{
		{
			"$match": bson.M{
				"is_approved": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         collections.Users,
				"localField":   "_id",
				"foreignField": "company",
				"as":           "employees",
			},
		},
		{
			"$addFields": bson.M{
				"total_employees": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": "$employees",
							"as":    "emp",
							"cond": bson.M{
								"$not": bson.M{
									"$in": []interface{}{
										"$$emp.role",
										[]models.Role{models.RoleSuperAdmin, models.RoleAdmin},
									},
								},
							},
						},
					},
				},
				"active_employees": bson.M{
					"$size": bson.M{
						"$filter": bson.M{
							"input": "$employees",
							"as":    "emp",
							"cond": bson.M{
								"$and": []bson.M{
									{
										"$eq": []interface{}{"$$emp.is_active", true},
									},
									{
										"$not": bson.M{
											"$in": []interface{}{
												"$$emp.role",
												[]models.Role{models.RoleSuperAdmin, models.RoleAdmin},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"$sort": bson.M{"total_employees": -1},
		},
		{
			"$limit": 10,
		},
		{
			"$project": bson.M{
				"name":             1,
				"total_employees":  1,
				"active_employees": 1,
			},
		},
	}

	companyCursor, err := companiesCollection.Aggregate(ctx, companyStatsPipeline)
	if err == nil {
		defer companyCursor.Close(ctx)
		for companyCursor.Next(ctx) {
			var result struct {
				Name            string `bson:"name"`
				TotalEmployees  int64  `bson:"total_employees"`
				ActiveEmployees int64  `bson:"active_employees"`
			}
			if err := companyCursor.Decode(&result); err == nil {
				stats.DepartmentStats = append(stats.DepartmentStats, CompanyStats{
					Name:            result.Name,
					TotalEmployees:  result.TotalEmployees,
					ActiveEmployees: result.ActiveEmployees,
					PresentToday:    0, // Will be calculated below
					OnLeave:         0, // Will be calculated below
				})
			}
		}
	}

	// Get attendance for today for each company
	for i := range stats.DepartmentStats {
		companyName := stats.DepartmentStats[i].Name

		// Find company ID by name
		var company models.Company
		err := companiesCollection.FindOne(ctx, bson.M{"name": companyName}).Decode(&company)
		if err == nil {
			// Present today for this company
			presentCount, err := attendanceCollection.CountDocuments(ctx, bson.M{
				"company": company.ID,
				"date":    today,
				"status":  models.StatusPresent,
			})
			if err == nil {
				stats.DepartmentStats[i].PresentToday = presentCount
			}

			// On leave today for this company
			onLeaveCount, err := leavesCollection.CountDocuments(ctx, bson.M{
				"company":    company.ID,
				"start_date": bson.M{"$lte": today},
				"end_date":   bson.M{"$gte": today},
				"status":     models.LeaveApproved,
			})
			if err == nil {
				stats.DepartmentStats[i].OnLeave = onLeaveCount
			}
		}
	}

	// Monthly attendance trends (last 6 months)
	now := time.Now()
	for i := 5; i >= 0; i-- {
		monthStart := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := monthStart.AddDate(0, 1, -1)

		monthStartStr := monthStart.Format("2006-01-02")
		monthEndStr := monthEnd.Format("2006-01-02")

		// Present count for the month
		presentCount, _ := attendanceCollection.CountDocuments(ctx, bson.M{
			"date": bson.M{
				"$gte": monthStartStr,
				"$lte": monthEndStr,
			},
			"status": models.StatusPresent,
		})

		// Absent count for the month
		absentCount, _ := attendanceCollection.CountDocuments(ctx, bson.M{
			"date": bson.M{
				"$gte": monthStartStr,
				"$lte": monthEndStr,
			},
			"status": models.StatusAbsent,
		})

		// On leave count for the month
		onLeaveCount, _ := leavesCollection.CountDocuments(ctx, bson.M{
			"start_date": bson.M{"$lte": monthEndStr},
			"end_date":   bson.M{"$gte": monthStartStr},
			"status":     models.LeaveApproved,
		})

		stats.MonthlyAttendance = append(stats.MonthlyAttendance, MonthlyAttendance{
			Month:   monthStart.Format("Jan 2006"),
			Present: presentCount,
			Absent:  absentCount,
			OnLeave: onLeaveCount,
		})
	}

	// Leave type statistics
	leaveTypes := []models.LeaveType{
		models.LeaveSick,
		models.LeaveCasual,
		models.LeaveVacation,
	}

	for _, leaveType := range leaveTypes {
		pending, _ := leavesCollection.CountDocuments(ctx, bson.M{
			"leave_type": leaveType,
			"status":     models.LeavePending,
		})

		approved, _ := leavesCollection.CountDocuments(ctx, bson.M{
			"leave_type": leaveType,
			"status":     models.LeaveApproved,
		})

		rejected, _ := leavesCollection.CountDocuments(ctx, bson.M{
			"leave_type": leaveType,
			"status":     models.LeaveRejected,
		})

		stats.LeaveTypeStats = append(stats.LeaveTypeStats, LeaveTypeStats{
			Type:     string(leaveType),
			Pending:  pending,
			Approved: approved,
			Rejected: rejected,
		})
	}

	return stats, nil
}
