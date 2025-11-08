package helpers

import "github.com/google/uuid"

func GetNewUUID() string {
	return uuid.New().String()
}

func VerifyUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
