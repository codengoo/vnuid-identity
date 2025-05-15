package helpers

import (
	"vnuid-identity/entities"
	"vnuid-identity/models"
	"vnuid-identity/utils"
)

func AddLoginSession(user entities.User, session entities.Session) (string, error) {
	token, err := utils.GenerateToken(user, session.DeviceId)
	if err != nil {
		return "", err
	}

	// Create login session
	session.UserId = user.ID
	if _, err := models.CreateSession(session); err != nil {
		return "", err
	}

	return token, nil
}
