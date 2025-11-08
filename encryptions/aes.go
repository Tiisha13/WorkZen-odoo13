package encryptions

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"api.workzen.odoo/constants"
)

// Fixed IV for encryption/decryption (16 bytes)
var fixedIV = constants.EncryptionAESIV

// EncryptAES encrypts a string using AES encryption
func EncryptAES(text string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(text))
	stream := cipher.NewCFBEncrypter(block, fixedIV)
	stream.XORKeyStream(ciphertext, []byte(text))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts a string using AES encryption
func DecryptAES(text string, key string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < 1 {
		return "", fmt.Errorf("ciphertext too short")
	}

	stream := cipher.NewCFBDecrypter(block, fixedIV)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func EncryptAESWithRounds(text string, key string, rounds int) (string, error) {
	for i := 0; i < rounds; i++ {
		encrypted, err := EncryptAES(text, key)
		if err != nil {
			return "", err
		}
		text = encrypted
	}
	return text, nil
}

func DecryptAESWithRounds(text string, key string, rounds int) (string, error) {
	for i := 0; i < rounds; i++ {
		decrypted, err := DecryptAES(text, key)
		if err != nil {
			return "", err
		}
		text = decrypted
	}
	return text, nil
}
