package models

import (
	"fmt"
	"vnuid-identity/databases"
	"vnuid-identity/entities"
	"vnuid-identity/utils"

	"github.com/google/uuid"
)

func isUUID(text string) bool {
	_, err := uuid.Parse(text)
	return err == nil
}

func GetUser(id string) (entities.User, error) {
	var user entities.User
	result := databases.DB.Where("id = ? OR email = ? OR sid = ? OR gid = ?", id, id, id, id).First(&user)

	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user, nil
}

func RemoveUsers(input []string) error {
	var uuids []string
	var emails []string

	for _, item := range input {
		if isUUID(item) {
			uuids = append(uuids, item)
		} else {
			emails = append(emails, item)
		}
	}

	result := databases.DB.Where("id IN ? OR email IN ?", uuids, emails).Delete(&entities.User{})
	if result.Error != nil {
		return fmt.Errorf("could not delete records: %v", result.Error)
	}

	return nil
}

func AddManyUser(input []entities.User) error {
	var users []entities.User

	for _, data := range input {
		password, err := utils.GeneratePassword()
		if err != nil {
			return fmt.Errorf("failed to generate password")
		}

		user := entities.User{
			ID:            uuid.New().String(),
			Email:         data.Email,
			SID:           data.SID,
			GID:           data.GID,
			Name:          data.Name,
			OfficialClass: data.OfficialClass,
			Type:          data.Type,
			Password:      password,
		}
		users = append(users, user)
	}

	if result := databases.DB.Create(&users); result.Error != nil {
		return fmt.Errorf("failed to create users: %v", result.Error)
	}

	return nil
}

func AddUser(input entities.User) error {
	password, err := utils.GeneratePassword()
	if err != nil {
		return fmt.Errorf("failed to generate password")
	}

	user := entities.User{
		Type:          input.Type,
		Email:         input.Email,
		SID:           input.SID,
		GID:           input.GID,
		Name:          input.Name,
		OfficialClass: input.OfficialClass,
		ID:            uuid.New().String(),
		Password:      password,
	}

	result := databases.DB.Create(&user)

	if result.Error != nil {
		return fmt.Errorf("failed to create user: %v", result.Error)
	}

	return nil
}
