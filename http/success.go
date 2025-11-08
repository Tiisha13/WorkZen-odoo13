package http

import "github.com/gofiber/fiber/v2"

type Success interface {
	// 2xx Success responses
	OK(c *fiber.Ctx, msg string, data interface{}) error
	Created(c *fiber.Ctx, msg string, data interface{}) error
	Accepted(c *fiber.Ctx, msg string, data interface{}) error
	NoContent(c *fiber.Ctx, msg string) error
	ResetContent(c *fiber.Ctx, msg string) error
	PartialContent(c *fiber.Ctx, msg string, data interface{}) error
	MultiStatus(c *fiber.Ctx, msg string, data interface{}) error
	AlreadyReported(c *fiber.Ctx, msg string, data interface{}) error
	IMUsed(c *fiber.Ctx, msg string, data interface{}) error
	EarlyHints(c *fiber.Ctx, msg string, data interface{}) error

	OKWithoutData(c *fiber.Ctx, msg string) error
	CreatedWithoutData(c *fiber.Ctx, msg string) error
	AcceptedWithoutData(c *fiber.Ctx, msg string) error
	PartialContentWithoutData(c *fiber.Ctx, msg string) error
	MultiStatusWithoutData(c *fiber.Ctx, msg string) error
}

func createSuccessResponse(c *fiber.Ctx, status int, msg string, defaultMsg string, data interface{}) error {
	if msg == "" {
		msg = defaultMsg
	}

	resp := fiber.Map{
		"message": msg,
		"success": true,
	}

	if data != nil {
		resp["data"] = data
	}

	return c.Status(status).JSON(resp)
}

type successImpl struct{}

func NewSuccess() Success {
	return &successImpl{}
}

func (s *successImpl) OK(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusOK, msg, "OK", data)
}

func (s *successImpl) Created(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusCreated, msg, "Created", data)
}

func (s *successImpl) Accepted(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusAccepted, msg, "Accepted", data)
}

func (s *successImpl) NoContent(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusNoContent, msg, "No Content", nil)
}

func (s *successImpl) ResetContent(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusResetContent, msg, "Reset Content", nil)
}

func (s *successImpl) PartialContent(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusPartialContent, msg, "Partial Content", data)
}
func (s *successImpl) MultiStatus(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusMultiStatus, msg, "Multi Status", data)
}
func (s *successImpl) AlreadyReported(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusAlreadyReported, msg, "Already Reported", data)
}
func (s *successImpl) IMUsed(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusIMUsed, msg, "IM Used", data)
}
func (s *successImpl) EarlyHints(c *fiber.Ctx, msg string, data interface{}) error {
	return createSuccessResponse(c, fiber.StatusEarlyHints, msg, "Early Hints", data)
}

func (s *successImpl) OKWithoutData(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusOK, msg, "OK", nil)
}
func (s *successImpl) CreatedWithoutData(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusCreated, msg, "Created", nil)
}

func (s *successImpl) AcceptedWithoutData(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusAccepted, msg, "Accepted", nil)
}

func (s *successImpl) PartialContentWithoutData(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusPartialContent, msg, "Partial Content", nil)
}

func (s *successImpl) MultiStatusWithoutData(c *fiber.Ctx, msg string) error {
	return createSuccessResponse(c, fiber.StatusMultiStatus, msg, "Multi Status", nil)
}
