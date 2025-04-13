package helpers

import (
	"vnuid-identity/entities"
	"vnuid-identity/models"
	"vnuid-identity/utils"
)

func AddLoginSession(user entities.User, device_id string, saved bool) (string, error) {
	token, err := utils.GenerateToken(user, device_id)
	if err != nil {
		return "", err
	}

	// Create login session
	if _, err := models.CreateSession(device_id, user.ID, saved); err != nil {
		return "", err
	}

	return token, nil
}
