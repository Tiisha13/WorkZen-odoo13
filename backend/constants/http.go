// Package constants provides global constants and shared instances for the application.
package constants

import "api.workzen.odoo/http"

var (
	HTTPSuccess = http.NewSuccess()
	HTTPErrors  = http.NewHTTPErrors()
)
