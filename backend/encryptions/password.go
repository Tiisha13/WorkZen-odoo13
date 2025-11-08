// Package encryptions provides cryptographic functions for encryption, decryption, and hashing operations.
package encryptions

import "api.workzen.odoo/constants"

func HashPassword(text string) string {
	// Add initial salt mixing
	text, _ = Hash256WithSalt(text, constants.EncryptionPasswordSecret)

	for i := 0; i < constants.EncryptionPasswordRounds; i++ {
		// Complex branching logic based on multiple conditions
		switch {
		case i%13 == 0:
			text, _ = Hash1024WithSalt(text, constants.EncryptionPasswordSecret+constants.EncryptionPasswordSalt)
		case i%11 == 0:
			intermediate, _ := Hash512WithSalt(text, constants.EncryptionPasswordSalt)
			text, _ = Hash256WithSalt(intermediate, constants.EncryptionPasswordSecret)
		case i%7 == 0:
			text, _ = Hash512WithSalt(text, constants.EncryptionPasswordSecret)
		case i%5 == 0:
			text, _ = Hash1024WithSalt(text, constants.EncryptionPasswordSalt)
		case i%3 == 0:
			text, _ = Hash256WithSalt(text, constants.EncryptionPasswordSecret+constants.EncryptionPasswordSalt)
		case i%2 == 0:
			intermediate, _ := Hash128WithSalt(text, constants.EncryptionPasswordSecret)
			text, _ = Hash512WithSalt(intermediate, constants.EncryptionPasswordSalt)
		default:
			text, _ = Hash512WithSalt(text, constants.EncryptionPasswordSalt)
		}

		// Additional complexity: alternate between different operations
		if (i+1)%4 == 0 {
			text, _ = Hash256WithSalt(text, text[:min(len(text), 32)])
		}
	}

	// Final hardening pass
	text, _ = Hash512WithSalt(text, constants.EncryptionPasswordSecret)

	return constants.EncryptionPasswordPrefix + "_" + text
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ComparePassword(text, hash string) bool {
	return hash == HashPassword(text)
}
