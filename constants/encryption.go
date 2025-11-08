package constants

import "api.workzen.odoo/config"

var (
	// encryption JWT
	EncryptionJWTSecret = config.GetConfig().GetString("encryption.jwt.secret")
	EncryptionJWTExpire = config.GetConfig().GetInt64("encryption.jwt.expire") // in hours

	// encryption AES
	EncryptionAESKey = config.GetConfig().GetString("encryption.aes.key")

	// encryption AES ID
	EncryptionAESIDKey = config.GetConfig().GetString("encryption.aes.ids.key")

	// encryption AES with rounds
	EncryptionAESRounds = config.GetConfig().GetInt("encryption.aes.rounds")

	// IV for AES CBC mode
	EncryptionAESIV = []byte(config.GetConfig().GetString("encryption.aes.iv"))

	// encryption HASH
	EncryptionHashSalt   = config.GetConfig().GetString("encryption.hash.salt")
	EncryptionHashSecret = config.GetConfig().GetString("encryption.hash.secret")
	EncryptionHashPrefix = config.GetConfig().GetString("encryption.hash.prefix")

	// encryption PASSWORD
	EncryptionPasswordSalt   = config.GetConfig().GetString("encryption.password.salt")
	EncryptionPasswordSecret = config.GetConfig().GetString("encryption.password.secret")
	EncryptionPasswordPrefix = config.GetConfig().GetString("encryption.password.prefix")
	EncryptionPasswordRounds = config.GetConfig().GetInt("encryption.password.rounds")
)
