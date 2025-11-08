// Package middlewares provides HTTP middleware functions for the WorkZen HRMS API,
// including authentication, authorization, and request processing.
package middlewares

import (
	"context"
	"strings"
	"time"

	"api.workzen.odoo/constants"
	"api.workzen.odoo/databases"
	"api.workzen.odoo/databases/collections"
	"api.workzen.odoo/databases/models"
	"api.workzen.odoo/encryptions"
	"api.workzen.odoo/helpers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthMiddleware verifies JWT token and extracts user information
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return constants.HTTPErrors.Unauthorized(c, "Authorization header required")
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return constants.HTTPErrors.Unauthorized(c, "Invalid authorization format")
		}

		token := parts[1]

		// Verify JWT
		valid, claims, err := helpers.VerifyJWT(token)
		if err != nil || !valid {
			return constants.HTTPErrors.Unauthorized(c, "Invalid or expired token")
		}

		// Extract encrypted user ID from claims and decrypt it
		encryptedUserID, ok := claims["id"].(string)
		if !ok {
			return constants.HTTPErrors.Unauthorized(c, "Invalid token claims")
		}

		// Decrypt user ID
		userIDStr, err := encryptions.DecryptID(encryptedUserID)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Invalid user ID in token")
		}

		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Invalid user ID in token")
		} 
		
		// Fetch user from database
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userCollection := databases.GetMongoDBCollection(collections.Users)
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "User not found")
		}

		// Check if user is active
		if user.Status != models.UserActive {
			return constants.HTTPErrors.Forbidden(c, "User account is inactive")
		}

		// Store user in context
		c.Locals("user", user)
		c.Locals("userID", user.ID)
		c.Locals("companyID", user.Company)
		c.Locals("role", user.Role)
		c.Locals("isSuperAdmin", user.IsSuperAdmin)

		return c.Next()
	}
}

// GetAuthUser retrieves the authenticated user from context
func GetAuthUser(c *fiber.Ctx) (*models.User, error) {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "User not found in context")
	}
	return &user, nil
}

// GetAuthUserID retrieves the authenticated user ID from context
func GetAuthUserID(c *fiber.Ctx) (primitive.ObjectID, error) {
	userID, ok := c.Locals("userID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}
	return userID, nil
}

// GetAuthCompanyID retrieves the authenticated user's company ID from context
func GetAuthCompanyID(c *fiber.Ctx) (primitive.ObjectID, error) {
	companyID, ok := c.Locals("companyID").(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusUnauthorized, "Company ID not found in context")
	}
	return companyID, nil
}
