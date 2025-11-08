package constants

import "api.workzen.odoo/config"

var (
	// server
	ServerPort  = config.GetConfig().GetString("server.port")
	ServerMode  = config.GetConfig().GetString("server.mode")
	FrontendURL = config.GetConfig().GetString("frontend_url")
)
