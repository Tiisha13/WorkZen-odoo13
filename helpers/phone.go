package helpers

import (
	"fmt"

	"github.com/nyaruka/phonenumbers"
)

// ValidatePhoneNumber validates if a phone number is in correct format
func ValidatePhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse(phone, "IN")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

// ParsePhoneNumber cleans and returns a standardized phone number (10 digits)
func ParsePhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	num, err := phonenumbers.Parse(phone, "IN")
	if err != nil {
		return "", fmt.Errorf("invalid phone number format")
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number format")
	}

	// Format as national number (10 digits without country code)
	national := phonenumbers.Format(num, phonenumbers.NATIONAL)
	// Remove spaces from national format
	cleaned := ""
	for _, c := range national {
		if c >= '0' && c <= '9' {
			cleaned += string(c)
		}
	}

	return cleaned, nil
}

// NormalizePhoneNumber returns phone number with +91 prefix
func NormalizePhoneNumber(phone string) (string, error) {
	if phone == "" {
		return "", fmt.Errorf("phone number cannot be empty")
	}

	num, err := phonenumbers.Parse(phone, "IN")
	if err != nil {
		return "", fmt.Errorf("invalid phone number format")
	}

	if !phonenumbers.IsValidNumber(num) {
		return "", fmt.Errorf("invalid phone number format")
	}

	// Format as E164 (+91XXXXXXXXXX)
	return phonenumbers.Format(num, phonenumbers.E164), nil
}
