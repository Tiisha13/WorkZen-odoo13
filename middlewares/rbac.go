package middlewares

import (
	"api.workzen.odoo/constants"
	"api.workzen.odoo/databases/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RequireRole creates middleware that checks if user has one of the specified roles
func RequireRole(roles ...models.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		// Check if user has any of the required roles
		for _, role := range roles {
			if user.Role == role {
				return c.Next()
			}
		}

		return constants.HTTPErrors.Forbidden(c, "Insufficient permissions")
	}
}

// RequireSuperAdmin checks if user is a SuperAdmin
func RequireSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		if !user.IsSuperAdmin {
			return constants.HTTPErrors.Forbidden(c, "SuperAdmin access required")
		}

		return c.Next()
	}
}

// RequireCompanyAdmin checks if user is Admin or SuperAdmin
func RequireCompanyAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		if !user.IsSuperAdmin && user.Role != models.RoleAdmin {
			return constants.HTTPErrors.Forbidden(c, "Admin access required")
		}

		return c.Next()
	}
}

// RequireHROrAdmin checks if user is HR, Admin, or SuperAdmin
func RequireHROrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		if !user.IsSuperAdmin && user.Role != models.RoleAdmin && user.Role != models.RoleHR {
			return constants.HTTPErrors.Forbidden(c, "HR or Admin access required")
		}

		return c.Next()
	}
}

// RequirePayrollOrAdmin checks if user is Payroll Officer, Admin, or SuperAdmin
func RequirePayrollOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		if !user.IsSuperAdmin && user.Role != models.RoleAdmin && user.Role != models.RolePayroll {
			return constants.HTTPErrors.Forbidden(c, "Payroll Officer or Admin access required")
		}

		return c.Next()
	}
}

// CompanyScopeMiddleware ensures user can only access resources from their own company
func CompanyScopeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		// SuperAdmin can access all companies
		if user.IsSuperAdmin {
			return c.Next()
		}

		// Store company scope for filtering queries
		c.Locals("companyScope", user.Company)

		return c.Next()
	}
}

// CanAccessEmployee checks if the authenticated user can access a specific employee's data
func CanAccessEmployee(targetEmployeeID primitive.ObjectID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		// SuperAdmin and Admin can access all employees in their scope
		if user.IsSuperAdmin || user.Role == models.RoleAdmin {
			return c.Next()
		}

		// HR can access all employees in the company
		if user.Role == models.RoleHR {
			return c.Next()
		}

		// Payroll can access all employees for payroll purposes
		if user.Role == models.RolePayroll {
			return c.Next()
		}

		// Employee can only access their own data
		if user.ID == targetEmployeeID {
			return c.Next()
		}

		return constants.HTTPErrors.Forbidden(c, "You can only access your own data")
	}
}

// CanModifySalaryInfo checks if user can view/edit salary information
func CanModifySalaryInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetAuthUser(c)
		if err != nil {
			return constants.HTTPErrors.Unauthorized(c, "Authentication required")
		}

		// Only Admin and Payroll can modify salary info
		if !user.IsSuperAdmin && user.Role != models.RoleAdmin && user.Role != models.RolePayroll {
			return constants.HTTPErrors.Forbidden(c, "Only Admin or Payroll Officer can modify salary information")
		}

		return c.Next()
	}
}
