package databases

import (
	"context"

	"api.workzen.odoo/constants"
)

func RedisFlush() error {
	info := RedisClient.Del(context.Background(), constants.DatabaseRedisPrefix+"*")
	return info.Err()
}
