// Package seed provides database seeding functionality for initial data setup
package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/encryptions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// SeedDatabase seeds the database with initial SuperAdmin and related data
func SeedDatabase(db *mongo.Database) error {
	log.Println("🌱 Starting database seeding...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if SuperAdmin already exists
	usersCollection := db.Collection("users")
	var existingSuperAdmin models.User
	err := usersCollection.FindOne(ctx, bson.M{"is_super_admin": true}).Decode(&existingSuperAdmin)
	if err == nil {
		log.Println("✅ SuperAdmin already exists. Skipping seed.")
		log.Printf("   Email: %s\n", existingSuperAdmin.Email)
		return nil
	}

	// Create SuperAdmin user
	superAdminID := primitive.NewObjectID()
	now := time.Now()

	hashedPassword := encryptions.HashPassword("SuperAdmin@123") // Default password

	superAdmin := models.User{
		ID:               superAdminID,
		Username:         "superadmin",
		Email:            "superadmin@workzen.com",
		Password:         hashedPassword,
		FirstName:        "Super",
		LastName:         "Admin",
		Role:             models.RoleSuperAdmin,
		IsSuperAdmin:     true,
		Status:           models.UserActive,
		Phone:            "+919999999999",
		EmailVerified:    true,
		TwoFactorEnabled: false,
		TimeStamp: models.TimeStamp{
			CreatedAt: primitive.NewDateTimeFromTime(now),
			UpdatedAt: primitive.NewDateTimeFromTime(now),
			IsDeleted: false,
		},
	}

	_, err = usersCollection.InsertOne(ctx, superAdmin)
	if err != nil {
		return fmt.Errorf("failed to create SuperAdmin: %w", err)
	}

	log.Println("✅ SuperAdmin created successfully!")
	log.Println("   ================================")
	log.Println("   Email:    superadmin@workzen.com")
	log.Println("   Password: SuperAdmin@123")
	log.Println("   Role:     SuperAdmin")
	log.Println("   ================================")
	log.Println("   ⚠️  IMPORTANT: Change the password after first login!")
	log.Println("")

	// Seed Demo Company (Optional)
	if err := seedDemoCompany(db, superAdminID, ctx); err != nil {
		log.Printf("⚠️  Warning: Failed to create demo company: %v\n", err)
		// Not returning error as SuperAdmin is the critical part
	}

	// Seed Demo Departments (Optional)
	if err := seedDemoDepartments(db, ctx); err != nil {
		log.Printf("⚠️  Warning: Failed to create demo departments: %v\n", err)
	}

	log.Println("🎉 Database seeding completed successfully!")
	return nil
}

// seedDemoCompany creates a demo company for testing
func seedDemoCompany(db *mongo.Database, superAdminID primitive.ObjectID, ctx context.Context) error {
	companiesCollection := db.Collection("companies")

	// Check if demo company exists
	var existingCompany models.Company
	err := companiesCollection.FindOne(ctx, bson.M{"email": "demo@workzen.com"}).Decode(&existingCompany)
	if err == nil {
		log.Println("✅ Demo company already exists. Skipping.")
		return nil
	}

	now := time.Now()
	demoCompanyID := primitive.NewObjectID()

	demoCompany := models.Company{
		ID:         demoCompanyID,
		Name:       "WorkZen Demo Company",
		Email:      "demo@workzen.com",
		Phone:      "+919876543210",
		Industry:   "Technology",
		Website:    "https://demo.workzen.com",
		IsApproved: true,
		IsActive:   true,
		ApprovedBy: superAdminID,
		OwnerID:    primitive.NewObjectID(), // Will be set when admin user is created
		TimeStamp: models.TimeStamp{
			CreatedAt: primitive.NewDateTimeFromTime(now),
			UpdatedAt: primitive.NewDateTimeFromTime(now),
		},
	}

	_, err = companiesCollection.InsertOne(ctx, demoCompany)
	if err != nil {
		return err
	}

	log.Println("✅ Demo company created!")
	log.Println("   Company: WorkZen Demo Company")
	log.Println("   Email:   demo@workzen.com")

	// Create demo admin user for the company
	return seedDemoAdmin(db, demoCompanyID, ctx)
}

// seedDemoAdmin creates a demo admin user for the demo company
func seedDemoAdmin(db *mongo.Database, companyID primitive.ObjectID, ctx context.Context) error {
	usersCollection := db.Collection("users")

	// Check if demo admin exists
	var existingAdmin models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": "admin@workzen.com"}).Decode(&existingAdmin)
	if err == nil {
		log.Println("✅ Demo admin already exists. Skipping.")
		return nil
	}

	now := time.Now()
	demoAdminID := primitive.NewObjectID()
	hashedPassword := encryptions.HashPassword("Admin@123")

	demoAdmin := models.User{
		ID:            demoAdminID,
		Username:      "demoadmin",
		Email:         "admin@workzen.com",
		Password:      hashedPassword,
		FirstName:     "Demo",
		LastName:      "Admin",
		Role:          models.RoleAdmin,
		IsSuperAdmin:  false,
		Status:        models.UserActive,
		Phone:         "+919876543211",
		Company:       companyID,
		EmployeeCode:  "WZ001",
		Designation:   "Chief Executive Officer",
		DateOfJoin:    time.Now().Format("2006-01-02"),
		EmailVerified: true,
		TimeStamp: models.TimeStamp{
			CreatedAt: primitive.NewDateTimeFromTime(now),
			UpdatedAt: primitive.NewDateTimeFromTime(now),
			IsDeleted: false,
		},
	}

	_, err = usersCollection.InsertOne(ctx, demoAdmin)
	if err != nil {
		return err
	}

	// Update company owner
	companiesCollection := db.Collection("companies")
	_, err = companiesCollection.UpdateOne(
		ctx,
		bson.M{"_id": companyID},
		bson.M{"$set": bson.M{"owner_id": demoAdminID}},
	)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to update company owner: %v\n", err)
	}

	log.Println("✅ Demo admin created!")
	log.Println("   Email:    admin@workzen.com")
	log.Println("   Password: Admin@123")
	log.Println("   Role:     Admin")
	log.Println("")

	return nil
}

// seedDemoDepartments creates demo departments for the demo company
func seedDemoDepartments(db *mongo.Database, ctx context.Context) error {
	departmentsCollection := db.Collection("departments")

	// Check if departments exist
	count, err := departmentsCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("✅ Departments already exist. Skipping.")
		return nil
	}

	now := time.Now()
	departments := []interface{}{
		models.Department{
			ID:          primitive.NewObjectID(),
			Name:        "Engineering",
			Description: "Software development and technical operations",
			TimeStamp: models.TimeStamp{
				CreatedAt: primitive.NewDateTimeFromTime(now),
				UpdatedAt: primitive.NewDateTimeFromTime(now),
			},
		},
		models.Department{
			ID:          primitive.NewObjectID(),
			Name:        "Human Resources",
			Description: "Employee management and recruitment",
			TimeStamp: models.TimeStamp{
				CreatedAt: primitive.NewDateTimeFromTime(now),
				UpdatedAt: primitive.NewDateTimeFromTime(now),
			},
		},
		models.Department{
			ID:          primitive.NewObjectID(),
			Name:        "Finance",
			Description: "Financial planning and payroll management",
			TimeStamp: models.TimeStamp{
				CreatedAt: primitive.NewDateTimeFromTime(now),
				UpdatedAt: primitive.NewDateTimeFromTime(now),
			},
		},
		models.Department{
			ID:          primitive.NewObjectID(),
			Name:        "Sales",
			Description: "Business development and sales operations",
			TimeStamp: models.TimeStamp{
				CreatedAt: primitive.NewDateTimeFromTime(now),
				UpdatedAt: primitive.NewDateTimeFromTime(now),
			},
		},
	}

	_, err = departmentsCollection.InsertMany(ctx, departments)
	if err != nil {
		return err
	}

	log.Println("✅ Demo departments created!")
	log.Println("   - Engineering")
	log.Println("   - Human Resources")
	log.Println("   - Finance")
	log.Println("   - Sales")
	log.Println("")

	return nil
}
