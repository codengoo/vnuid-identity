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

// Check if user first logged in or saved this device
func CheckSession(device_id string, uid string) bool {
	var count int64
	databases.DB.Model(&entities.Session{}).Where("user_id = ?", uid).Count(&count)

	// First time logged in
	if count == 0 {
		return true
	}

	var session entities.Session
	result := databases.DB.Where("device_id = ? AND user_id = ?", device_id, uid).Order("created_at DESC").First(&session)

	if result.Error != nil {
		return false
	}

	return session.SavedDevice
}
