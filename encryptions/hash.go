package encryptions

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"fmt"

	"api.workzen.odoo/constants"
)

// Hash hashes a string using the specified algorithm
func Hash(text string, bits int) (string, error) {
	var hash []byte
	var err error

	switch bits {
	case 128:
		hasher := sha256.New()
		hasher.Write([]byte(text))
		hash = hasher.Sum(nil)
		hash = hash[:16]
	case 256:
		hasher := sha256.New()
		hasher.Write([]byte(text))
		hash = hasher.Sum(nil)
	case 512:
		hasher := sha512.New()
		hasher.Write([]byte(text))
		hash = hasher.Sum(nil)
	case 1024:
		hasher := sha512.New512_256()
		hasher.Write([]byte(text))
		hash = hasher.Sum(nil)
	default:
		err = fmt.Errorf("invalid bits: %d", bits)
	}

	if bits == 128 || bits == 512 {
		return constants.EncryptionHashPrefix + "_" + base32.HexEncoding.EncodeToString(hash), err
	} else {
		return constants.EncryptionHashPrefix + "_" + base64.URLEncoding.EncodeToString(hash), err
	}
}

// Hash128 hashes a string using SHA-256 and returns the first 128 bits
func Hash128(text string) (string, error) {
	return Hash(text, 128)
}

// Hash256 hashes a string using SHA-256
func Hash256(text string) (string, error) {
	return Hash(text, 256)
}

// Hash512 hashes a string using SHA-512
func Hash512(text string) (string, error) {
	return Hash(text, 512)
}

// Hash1024 hashes a string using SHA-512/256
func Hash1024(text string) (string, error) {
	return Hash(text, 1024)
}

// HashWithSalt hashes a string with a salt using the specified algorithm
func HashWithSalt(text string, salt string, bits int) (string, error) {
	return Hash(text+salt, bits)
}

func Hash128WithSalt(text string, salt string) (string, error) {
	return HashWithSalt(text, salt, 128)
}

func Hash256WithSalt(text string, salt string) (string, error) {
	return HashWithSalt(text, salt, 256)
}

func Hash512WithSalt(text string, salt string) (string, error) {
	return HashWithSalt(text, salt, 512)
}

func Hash1024WithSalt(text string, salt string) (string, error) {
	return HashWithSalt(text, salt, 1024)
}
