package databases

import (
	"context"

	"api.workzen.odoo/constants"

	"github.com/go-redis/redis/v8"
)

// RedisClient is the Redis client
var RedisClient *redis.Client

// ConnectRedis connects to Redis
func ConnectRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     constants.DatabaseRedisHost + ":" + constants.DatabaseRedisPort,
		Password: constants.DatabaseRedisPassword,
		DB:       constants.DatabaseRedisDB,
	})

	err := RedisClient.Ping(context.Background()).Err()

	return err
}

// DisconnectRedis disconnects from Redis
func DisconnectRedis() error {
	if err := RedisClient.Close(); err != nil {
		return err
	}

	return nil
}
