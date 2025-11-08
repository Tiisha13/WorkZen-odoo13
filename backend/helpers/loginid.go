package helpers

import (
	"fmt"
	"strings"
	"time"
)

// GenerateLoginID creates a unique username following the pattern:
// [CompanyCode][First2LettersFirstName][First2LettersLastName][YearOfJoining][SerialNumber]
func GenerateLoginID(companyCode, firstName, lastName string, joinDate time.Time, serial int) string {
	// Clean and extract initials
	fInitial := "XX"
	lInitial := "XX"

	if len(firstName) >= 2 {
		fInitial = strings.ToUpper(firstName[:2])
	} else if len(firstName) == 1 {
		fInitial = strings.ToUpper(firstName + "X")
	}

	if len(lastName) >= 2 {
		lInitial = strings.ToUpper(lastName[:2])
	} else if len(lastName) == 1 {
		lInitial = strings.ToUpper(lastName + "X")
	}

	// Format year and serial number
	year := joinDate.Format("2006")
	serialStr := fmt.Sprintf("%04d", serial)

	// Build final username
	loginID := fmt.Sprintf("%s%s%s%s%s", strings.ToUpper(companyCode), fInitial, lInitial, year, serialStr)
	return loginID
}
