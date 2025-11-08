package constants

import "api.workzen.odoo/config"

var (
	// mongodb
	DatabaseMongodbURI    = config.GetConfig().GetString("databases.mongodb.uri")
	DatabaseMongodbDBName = config.GetConfig().GetString("databases.mongodb.db")
)
