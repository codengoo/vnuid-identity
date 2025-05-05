package models

import (
	"fmt"
	"vnuid-identity/databases"
	"vnuid-identity/entities"

	"github.com/google/uuid"
)

func CountActiveNFC(uid string) int64 {
	var count int64
	databases.DB.Model(&entities.NFC{}).Where("active = ? AND user_id = ?", true, uid).Count(&count)
	return count
}

func AddNFC(uid string) (string, error) {
	nfc := entities.NFC{
		ID:     uuid.New().String(),
		UserId: uid,
		Active: true,
		Status: "pending_release",
	}

	result := databases.DB.Create(&nfc)

	if result.Error != nil {
		return "", result.Error
	} else {
		return nfc.ID, nil
	}
}

func AddManyNFC(users []entities.User) ([]entities.NFC, error) {
	var nfcs []entities.NFC
	for _, user := range users {
		nfc := entities.NFC{
			UserId: user.ID,
			ID:     uuid.New().String(),
			Active: true,
			Status: "pending_release",
		}

		nfcs = append(nfcs, nfc)
	}

	if result := databases.DB.CreateInBatches(&nfcs, 50); result.Error != nil {
		return []entities.NFC{}, fmt.Errorf("failed to create users: %v", result.Error)
	}

	return nfcs, nil
}

func SetActiveNFC(id string, active bool) error {
	result := databases.DB.Model(&entities.NFC{}).Where("id = ?", id).Update("active", active)
	return result.Error
}
