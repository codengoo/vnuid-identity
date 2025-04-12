package models

import (
	"vnuid-identity/databases"
	"vnuid-identity/entities"

	"github.com/google/uuid"
)

func CreateSession(device_id string, uid string, saved bool) (string, error) {
	session := entities.Session{
		DeviceId:    device_id,
		LoginMethod: "google",
		SavedDevice: saved,
		UserId:      uid,
		ID:          uuid.New().String(),
	}

	result := databases.DB.Create(&session)

	if result.Error != nil {
		return "", result.Error
	} else {
		return session.ID, nil
	}
}

func CheckSession(device_id string, uid string) bool {
	var count int64
	databases.DB.Model(&entities.User{}).Where("uid = ?", uid).Count(&count)
	if count == 0 {
		return true
	}

	result := databases.DB.Where("device_id = ? AND user_id = ? AND active = ?", device_id, uid, true).First(&entities.Session{})

	return result.Error != nil
}
