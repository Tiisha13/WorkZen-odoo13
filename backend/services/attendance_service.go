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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttendanceService struct{}

func NewAttendanceService() *AttendanceService {
	return &AttendanceService{}
}

// CheckIn creates a new attendance record or resets existing one
func (s *AttendanceService) CheckIn(employeeID, companyID primitive.ObjectID) (*models.Attendance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)

	today := helpers.FormatDate(time.Now())
	now := time.Now()

	// Check if already checked in today
	var existingAttendance models.Attendance
	err := attendanceCollection.FindOne(ctx, bson.M{
		"employee_id": employeeID,
		"date":        today,
	}).Decode(&existingAttendance)

	if err == nil {
		// Record exists - allow re-check-in only if already checked out
		if existingAttendance.CheckOut == "" {
			return nil, errors.New("already checked in today, please check out first")
		}
		// Re-check-in: Reset the record
		return nil, errors.New("already completed attendance for today")
	}

	// Create new attendance record
	attendance := models.Attendance{
		ID:         primitive.NewObjectID(),
		EmployeeID: employeeID,
		Company:    companyID,
		Date:       today,
		CheckIn:    now.Format("15:04:05"),
		Status:     models.StatusPresent,
	}
	attendance.CreatedAt = primitive.NewDateTimeFromTime(now)
	attendance.UpdatedAt = primitive.NewDateTimeFromTime(now)

	_, err = attendanceCollection.InsertOne(ctx, attendance)
	if err != nil {
		return nil, fmt.Errorf("failed to check in: %w", err)
	}

	return &attendance, nil
}

// CheckOut updates attendance with check-out time
func (s *AttendanceService) CheckOut(employeeID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)

	today := helpers.FormatDate(time.Now())
	now := time.Now()

	// Find today's attendance
	var attendance models.Attendance
	err := attendanceCollection.FindOne(ctx, bson.M{
		"employee_id": employeeID,
		"date":        today,
	}).Decode(&attendance)
	if err != nil {
		return errors.New("no check-in found for today")
	}

	// Calculate work hours
	checkOutTime := now.Format("15:04:05")
	workHours, err := helpers.CalculateWorkHours(attendance.CheckIn, checkOutTime)
	if err != nil {
		workHours = 0
	}

	// Update attendance
	result, err := attendanceCollection.UpdateOne(
		ctx,
		bson.M{"_id": attendance.ID},
		bson.M{
			"$set": bson.M{
				"check_out":  checkOutTime,
				"work_hours": workHours,
				"updated_at": primitive.NewDateTimeFromTime(now),
			},
		},
	)
	if err != nil || result.MatchedCount == 0 {
		return errors.New("failed to check out")
	}

	return nil
}

// GetMyAttendance retrieves attendance for current month
func (s *AttendanceService) GetMyAttendance(employeeID primitive.ObjectID, month string) ([]models.Attendance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)

	// If month not provided, use current month
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	// Find all attendance for the month
	filter := bson.M{
		"employee_id": employeeID,
		"date": bson.M{
			"$regex": "^" + month,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, err := attendanceCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attendances []models.Attendance
	if err = cursor.All(ctx, &attendances); err != nil {
		return nil, err
	}

	return attendances, nil
}

// ListAttendance retrieves attendance with filters
func (s *AttendanceService) ListAttendance(companyID primitive.ObjectID, filters map[string]interface{}, page, limit int64) ([]models.Attendance, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)

	// Add company filter
	filters["company"] = companyID

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := attendanceCollection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var attendances []models.Attendance
	if err = cursor.All(ctx, &attendances); err != nil {
		return nil, 0, err
	}

	total, err := attendanceCollection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// ResetAttendance deletes today's attendance to allow re-check-in
func (s *AttendanceService) ResetAttendance(employeeID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)

	today := helpers.FormatDate(time.Now())

	result, err := attendanceCollection.DeleteOne(ctx, bson.M{
		"employee_id": employeeID,
		"date":        today,
	})
	if err != nil {
		return errors.New("failed to reset attendance")
	}
	if result.DeletedCount == 0 {
		return errors.New("no attendance found for today")
	}

	return nil
}

// GetAttendanceSummary returns summary statistics
func (s *AttendanceService) GetAttendanceSummary(companyID primitive.ObjectID, date string) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attendanceCollection := databases.MongoDBDatabase.Collection(collections.Attendances)

	if date == "" {
		date = helpers.FormatDate(time.Now())
	}

	presentCount, _ := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date":    date,
		"status":  models.StatusPresent,
	})

	onLeaveCount, _ := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date":    date,
		"status":  models.StatusOnLeave,
	})

	absentCount, _ := attendanceCollection.CountDocuments(ctx, bson.M{
		"company": companyID,
		"date":    date,
		"status":  models.StatusAbsent,
	})

	return map[string]int64{
		"present":  presentCount,
		"on_leave": onLeaveCount,
		"absent":   absentCount,
	}, nil
}
