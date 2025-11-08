package constants

import "api.workzen.odoo/config"

var (
	// smtp
	SMTPHost      = config.GetConfig().GetString("smtp.host")
	SMTPPort      = config.GetConfig().GetString("smtp.port")
	SMTPUsername  = config.GetConfig().GetString("smtp.username")
	SMTPPassword  = config.GetConfig().GetString("smtp.password")
	SMTPFromEmail = config.GetConfig().GetString("smtp.from_email")
	SMTPFromName  = config.GetConfig().GetString("smtp.from_name")
)
