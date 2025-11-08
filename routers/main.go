// Package routers provides HTTP routing configuration and initialization for the accounts service.
package routers

import (
	"api.workzen.odoo/constants"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Init initializes the router
func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		Prefork:       isProduction(),
		CaseSensitive: true,
		ServerHeader:  "Accounts Server: WorkFlecks",
		AppName:       "Accounts Server: WorkFlecks",
	})

	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(etag.New(etag.Config{Weak: true}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(helmet.New())

	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${time} ${status} - ${method} ${latency} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Kolkata",
	}))

	// api := app.Group("/api/v1")

	// 404 Route
	app.Use(func(c *fiber.Ctx) error {
		return constants.HTTPErrors.NotFound(c, "The requested resource was not found")
	})

	return app
}

func isProduction() bool {
	return constants.ServerMode == "production"
}
