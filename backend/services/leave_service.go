package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveService struct{}

func NewLeaveService() *LeaveService {
	return &LeaveService{}
}

// ApplyLeaveRequest for creating leave application
type ApplyLeaveRequest struct {
	LeaveType string `json:"leave_type" validate:"required"`
	Reason    string `json:"reason" validate:"required"`
	StartDate string `json:"start_date" validate:"required"` // YYYY-MM-DD
	EndDate   string `json:"end_date" validate:"required"`   // YYYY-MM-DD
}

// ApplyLeave creates a new leave request
func (s *LeaveService) ApplyLeave(req *ApplyLeaveRequest, employeeID, companyID primitive.ObjectID) (*models.Leave, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)

	// Parse dates
	startDate, err := helpers.ParseDate(req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}
	endDate, err := helpers.ParseDate(req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format")
	}

	// Validate dates
	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Calculate days
	days := int(endDate.Sub(startDate).Hours()/24) + 1

	// Create leave
	leave := models.Leave{
		ID:         primitive.NewObjectID(),
		EmployeeID: employeeID,
		Company:    companyID,
		LeaveType:  models.LeaveType(req.LeaveType),
		Reason:     req.Reason,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Days:       days,
		Status:     models.LeavePending,
	}
	leave.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	leave.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err = leavesCollection.InsertOne(ctx, leave)
	if err != nil {
		return nil, fmt.Errorf("failed to apply leave: %w", err)
	}

	return &leave, nil
}

// LeaveWithUser represents a leave with populated user data
type LeaveWithUser struct {
	models.Leave `bson:",inline"`
	User         *models.User `bson:"user,omitempty" json:"user,omitempty"`
}

// ListLeaves retrieves leaves with filters and populates user data
func (s *LeaveService) ListLeaves(companyID primitive.ObjectID, filters map[string]interface{}, page, limit int64) ([]LeaveWithUser, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)

	// Add company filter
	filters["company"] = companyID

	// Convert employee_id filter if it's a string
	if empIDStr, ok := filters["employee_id"].(string); ok {
		if empID, err := primitive.ObjectIDFromHex(empIDStr); err == nil {
			filters["employee_id"] = empID
		}
	}

	skip := (page - 1) * limit

	// Build aggregation pipeline to join with users
	pipeline := []bson.M{
		{"$match": filters},
		{"$sort": bson.M{"created_at": -1}},
		{"$skip": skip},
		{"$limit": limit},
		{
			"$lookup": bson.M{
				"from":         collections.Users,
				"localField":   "employee_id",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$user",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	cursor, err := leavesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var leaves []LeaveWithUser
	if err = cursor.All(ctx, &leaves); err != nil {
		return nil, 0, err
	}

	total, err := leavesCollection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

// ApproveLeave approves a leave request and updates attendance
func (s *LeaveService) ApproveLeave(leaveID, approvedByID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)

	// Get leave details
	var leave models.Leave
	err := leavesCollection.FindOne(ctx, bson.M{"_id": leaveID}).Decode(&leave)
	if err != nil {
		return errors.New("leave not found")
	}

	if leave.Status != models.LeavePending {
		return errors.New("leave is not pending")
	}

	// Update leave status
	now := time.Now()
	result, err := leavesCollection.UpdateOne(
		ctx,
		bson.M{"_id": leaveID},
		bson.M{
			"$set": bson.M{
				"status":      models.LeaveApproved,
				"approved_by": approvedByID,
				"reviewed_at": helpers.FormatDateTime(now),
				"updated_at":  primitive.NewDateTimeFromTime(now),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("failed to approve leave")
	}

	// Create attendance records for leave days
	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)
	startDate, _ := helpers.ParseDate(leave.StartDate)
	endDate, _ := helpers.ParseDate(leave.EndDate)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		attendance := models.Attendance{
			ID:         primitive.NewObjectID(),
			EmployeeID: leave.EmployeeID,
			Company:    leave.Company,
			Date:       helpers.FormatDate(d),
			Status:     models.StatusOnLeave,
			Remarks:    "Approved leave: " + string(leave.LeaveType),
		}
		attendance.CreatedAt = primitive.NewDateTimeFromTime(now)
		attendance.UpdatedAt = primitive.NewDateTimeFromTime(now)

		attendanceCollection.InsertOne(ctx, attendance)
	}

	return nil
}

// RejectLeave rejects a leave request
func (s *LeaveService) RejectLeave(leaveID, rejectedByID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	leavesCollection := databases.MongoDBDatabase.Collection(collections.Leaves)

	// Get leave details
	var leave models.Leave
	err := leavesCollection.FindOne(ctx, bson.M{"_id": leaveID}).Decode(&leave)
	if err != nil {
		return errors.New("leave not found")
	}

	if leave.Status != models.LeavePending {
		return errors.New("leave is not pending")
	}

	// Update leave status
	result, err := leavesCollection.UpdateOne(
		ctx,
		bson.M{"_id": leaveID},
		bson.M{
			"$set": bson.M{
				"status":      models.LeaveRejected,
				"rejected_by": rejectedByID,
				"reviewed_at": helpers.FormatDateTime(time.Now()),
				"updated_at":  primitive.NewDateTimeFromTime(time.Now()),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("failed to reject leave")
	}

	return nil
}
