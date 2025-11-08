package constants

import "api.workzen.odoo/config"

var (
	// mongodb
	DatabaseMongodbURI    = config.GetConfig().GetString("databases.mongodb.uri")
	DatabaseMongodbDBName = config.GetConfig().GetString("databases.mongodb.db")

	// redis
	DatabaseRedisHost     = config.GetConfig().GetString("databases.redis.host")
	DatabaseRedisPort     = config.GetConfig().GetString("databases.redis.port")
	DatabaseRedisPassword = config.GetConfig().GetString("databases.redis.password")
	DatabaseRedisDB       = config.GetConfig().GetInt("databases.redis.db")
	DatabaseRedisPrefix   = config.GetConfig().GetString("databases.redis.prefix")
)
