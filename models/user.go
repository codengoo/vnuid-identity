package models

import (
	"fmt"
	"vnuid-identity/databases"
	"vnuid-identity/entities"
	"vnuid-identity/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func isUUID(text string) bool {
	_, err := uuid.Parse(text)
	return err == nil
}

func GetUser(id string) (entities.User, error) {
	var user entities.User
	var result *gorm.DB
	if isUUID(id) {
		result = databases.DB.Model(&entities.User{}).Where("id = ?", id).First(&user)
	} else {
		result = databases.DB.Model(&entities.User{}).Where("email = ? OR s_id = ? OR g_id = ?", id, id, id).First(&user)
	}

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
	type Result struct {
		User entities.User
		Err  error
	}
	ch := make(chan Result, len(input))

	// Create routines
	for _, u := range input {
		go func(data entities.User) {
			pass, err := utils.GeneratePassword()
			data.ID = uuid.New().String()
			data.Password = pass
			ch <- Result{User: data, Err: err}
		}(u)
	}

	var users []entities.User
	for range input {
		res := <-ch
		if res.Err != nil {
			return fmt.Errorf("failed to generate password")
		}
		users = append(users, res.User)
	}
	if result := databases.DB.CreateInBatches(&users, 50); result.Error != nil {
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
