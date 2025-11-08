package helpers

import "time"

// GenerateVerificationToken generates a UUID-based verification token
func GenerateVerificationToken() (string, error) {
	return GetNewUUID(), nil
}

// VerificationTokenExpiry returns the expiry time for verification tokens (24 hours)
func VerificationTokenExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}
