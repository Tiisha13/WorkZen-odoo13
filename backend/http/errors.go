// Package http provides HTTP utilities and error handling for the accounts service.
package http

import "github.com/gofiber/fiber/v2"

type Error interface {
	// 4xx errors
	InvalidBody(c *fiber.Ctx, msg string) error
	InvalidParams(c *fiber.Ctx, msg string) error
	InvalidQuery(c *fiber.Ctx, msg string) error
	BadRequest(c *fiber.Ctx, msg string) error
	Already(c *fiber.Ctx, msg string) error
	NotFound(c *fiber.Ctx, msg string) error
	Unauthorized(c *fiber.Ctx, msg string) error
	Forbidden(c *fiber.Ctx, msg string) error
	Conflict(c *fiber.Ctx, msg string) error

	// 5xx errors
	InternalServerError(c *fiber.Ctx, msg string) error
	NotImplemented(c *fiber.Ctx, msg string) error
	ServiceUnavailable(c *fiber.Ctx, msg string) error
	GatewayTimeout(c *fiber.Ctx, msg string) error

	// Custom errors
	Custom(c *fiber.Ctx, status int, msg string) error
}

// Helper function to create error responses
func createErrorResponse(c *fiber.Ctx, status int, msg string, defaultMsg string) error {
	if msg == "" {
		msg = defaultMsg
	}
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"success": false,
	})
}

type httpErrors struct{}

func NewHTTPErrors() Error {
	return &httpErrors{}
}

// 4xx errors
func (e *httpErrors) InvalidBody(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusBadRequest, msg, "Invalid Body")
}

func (e *httpErrors) InvalidParams(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusBadRequest, msg, "Invalid Params")
}

func (e *httpErrors) InvalidQuery(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusBadRequest, msg, "Invalid Query")
}

func (e *httpErrors) BadRequest(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusBadRequest, msg, "Bad Request")
}

func (e *httpErrors) Already(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusAlreadyReported, msg, "Already Exists")
}

func (e *httpErrors) NotFound(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusNotFound, msg, "Not Found")
}

func (e *httpErrors) Unauthorized(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusUnauthorized, msg, "Unauthorized")
}

func (e *httpErrors) Forbidden(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusForbidden, msg, "Forbidden")
}

func (e *httpErrors) Conflict(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusConflict, msg, "Conflict")
}

// 5xx errors
func (e *httpErrors) InternalServerError(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusInternalServerError, msg, "Internal Server Error")
}

func (e *httpErrors) NotImplemented(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusNotImplemented, msg, "Not Implemented")
}

func (e *httpErrors) ServiceUnavailable(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusServiceUnavailable, msg, "Service Unavailable")
}

func (e *httpErrors) GatewayTimeout(c *fiber.Ctx, msg string) error {
	return createErrorResponse(c, fiber.StatusGatewayTimeout, msg, "Gateway Timeout")
}

// Custom errors
func (e *httpErrors) Custom(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{
		"message": msg,
		"success": false,
	})
}
