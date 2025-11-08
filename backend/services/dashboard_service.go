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

// AdminDashboardStats for company admins
type AdminDashboardStats struct {
	TotalEmployees       int64 `json:"total_employees"`
	ActiveEmployees      int64 `json:"active_employees"`
	InactiveEmployees    int64 `json:"inactive_employees"`
	PresentToday         int64 `json:"present_today"`
	AbsentToday          int64 `json:"absent_today"`
	OnLeaveToday         int64 `json:"on_leave_today"`
	PendingLeaves        int64 `json:"pending_leaves"`
	ApprovedLeaves       int64 `json:"approved_leaves"`
	RejectedLeaves       int64 `json:"rejected_leaves"`
	MissingBankAccounts  int64 `json:"missing_bank_accounts"`
	MissingManagers      int64 `json:"missing_managers"`
	TotalPayrollThisYear int64 `json:"total_payroll_this_year"`
}

// GetAdminDashboard retrieves dashboard stats for company admins
func (s *DashboardService) GetAdminDashboard(companyID primitive.ObjectID) (*AdminDashboardStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &AdminDashboardStats{}

	// Collections
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)
	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)
	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)

	// Total employees
	total, err := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"role":    bson.M{"$ne": models.RoleSuperAdmin},
	})
	if err == nil {
		stats.TotalEmployees = total
	}

	// Active employees
	active, err := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.UserActive,
		"role":    bson.M{"$ne": models.RoleSuperAdmin},
	})
	if err == nil {
		stats.ActiveEmployees = active
	}

	// Inactive employees
	inactive, err := usersCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.UserInactive,
		"role":    bson.M{"$ne": models.RoleSuperAdmin},
	})
	if err == nil {
		stats.InactiveEmployees = inactive
	}

	// Today's date for attendance
	today := time.Now().Format("2006-01-02")

	// Present today
	present, err := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date":    today,
		"status":  models.StatusPresent,
	})
	if err == nil {
		stats.PresentToday = present
	}

	// Absent today
	absent, err := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date":    today,
		"status":  models.StatusAbsent,
	})
	if err == nil {
		stats.AbsentToday = absent
	}

	// On leave today
	onLeave, err := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date":    today,
		"status":  models.StatusOnLeave,
	})
	if err == nil {
		stats.OnLeaveToday = onLeave
	}

	// Pending leaves
	pending, err := leavesCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.LeavePending,
	})
	if err == nil {
		stats.PendingLeaves = pending
	}

	// Approved leaves
	approved, err := leavesCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.LeaveApproved,
	})
	if err == nil {
		stats.ApprovedLeaves = approved
	}

	// Rejected leaves
	rejected, err := leavesCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"status":  models.LeaveRejected,
	})
	if err == nil {
		stats.RejectedLeaves = rejected
	}

	// Missing bank accounts
	missingBank, err := usersCollection.CountDocuments(ctx, bson.M{
		"company":      companyID,
		"status":       models.UserActive,
		"bank_details": bson.M{"$exists": false},
		"role":         bson.M{"$ne": models.RoleSuperAdmin},
	})
	if err == nil {
		stats.MissingBankAccounts = missingBank
	}

	// Missing managers
	missingManager, err := usersCollection.CountDocuments(ctx, bson.M{
		"company":    companyID,
		"status":     models.UserActive,
		"manager_id": primitive.NilObjectID,
		"role":       bson.M{"$ne": models.RoleSuperAdmin},
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

	return stats, nil
}

// SuperAdminDashboardStats for platform-level stats
type SuperAdminDashboardStats struct {
	TotalCompanies        int64   `json:"total_companies"`
	ActiveCompanies       int64   `json:"active_companies"`
	PendingApprovals      int64   `json:"pending_approvals"`
	TotalEmployees        int64   `json:"total_employees"`
	TotalPayrollProcessed int64   `json:"total_payroll_processed"`
	TotalPayrunsGenerated int64   `json:"total_payruns_generated"`
	PlatformRevenue       float64 `json:"platform_revenue"` // Could be calculated based on pricing model
}

// GetSuperAdminDashboard retrieves platform-wide statistics
func (s *DashboardService) GetSuperAdminDashboard() (*SuperAdminDashboardStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stats := &SuperAdminDashboardStats{}

	// Collections
	companiesCollection := databases.MongoDBDatabase.Collection(collections.Companies)
	usersCollection := databases.MongoDBDatabase.Collection(collections.Users)
	payrollCollection := databases.MongoDBDatabase.Collection(collections.Payrolls)
	payrunsCollection := databases.MongoDBDatabase.Collection(collections.Payruns)

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

	// Total employees across platform
	totalEmployees, err := usersCollection.CountDocuments(ctx, bson.M{
		"role": bson.M{"$ne": models.RoleSuperAdmin},
	})
	if err == nil {
		stats.TotalEmployees = totalEmployees
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

	return stats, nil
}
