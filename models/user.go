package models

import (
	"fmt"
	"vnuid-identity/databases"
	"vnuid-identity/entities"
	"vnuid-identity/utils"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
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
		result = databases.DB.Model(&entities.User{}).Where("id = ? OR profile_id = ?", id, id).First(&user)
	} else {
		result = databases.DB.Model(&entities.User{}).Where("email = ? OR sid = ? OR gid = ?", id, id, id).First(&user)
	}

	if result.Error != nil {
		return entities.User{}, result.Error
	}
	return user, nil
}

func GetMe(id string) (entities.Profile, error) {
	var user entities.Profile
	result := databases.DB.Model(&entities.Profile{}).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return entities.Profile{}, result.Error
	}
	return user, nil
}

func RemoveUsers(input []string) error {
	result := databases.DB.Where("id IN ? OR profile_id IN ?", input, input).Delete(&entities.User{})
	if result.Error != nil {
		return fmt.Errorf("could not delete records: %v", result.Error)
	}

	return nil
}

func AddUser(input entities.User) (entities.User, error) {
	password, err := utils.GeneratePassword(input.Password)
	if err != nil {
		return entities.User{}, fmt.Errorf("failed to generate password")
	}

	user := entities.User{
		Type:      input.Type,
		Email:     input.Email,
		Sid:       input.Sid,
		Gid:       input.Gid,
		ID:        uuid.New().String(),
		Password:  password,
		ProfileId: input.ProfileId,
	}

	result := databases.DB.Create(&user)

	if result.Error != nil {
		// find profile and delete
		databases.DB.Where("id = ?", user.ProfileId).Delete(&entities.Profile{})
		return entities.User{}, fmt.Errorf("failed to create user: %v", result.Error)
	}

	return user, nil
}

func AddUserInfo(input entities.Profile) (entities.Profile, error) {
	input.ID = uuid.New().String()

	result := databases.DB.Create(&input)
	if result.Error != nil {
		return entities.Profile{}, fmt.Errorf("failed to create user: %v", result.Error)
	}

	return input, nil
}

func VerifyPassword(id string, password string) (bool, entities.User) {
	user, err := GetUser(id)

	if err != nil {
		return false, entities.User{}
	}

	isValid := utils.VerifyPassword(user.Password, password)

	if isValid {
		return true, user
	} else {
		return false, entities.User{}
	}
}

func SetAuthenticator(id string, authenticator string) error {
	user, err := GetUser(id)

	if err != nil {
		return err
	}

	user.Authenticator = authenticator
	result := databases.DB.Save(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SetBiometric(id string) (string, error) {
	user, err := GetUser(id)

	if err != nil {
		return "", err
	}

	user.BiometricKey = uuid.New().String()
	result := databases.DB.Save(&user)

	if result.Error != nil {
		return "", result.Error
	}

	return user.BiometricKey, nil
}

func VerifyAuthenticator(id string, code string) (bool, entities.User) {
	user, err := GetUser(id)

	if err != nil {
		return false, entities.User{}
	}

	isValid := totp.Validate(code, user.Authenticator)

	if isValid {
		return true, user
	} else {
		return false, entities.User{}
	}
}

func VerifyNFC(id string, code string) (bool, entities.User) {
	user, err := GetUser(id)

	if err != nil {
		return false, entities.User{}
	}

	result := databases.DB.Model(&entities.NFC{}).Where("id = ? AND user_id = ? AND active = ?", code, user.ID, true).First(&entities.NFC{})

	if result.Error != nil {
		return false, entities.User{}
	} else {
		return true, user
	}
}

func VerifyBioCode(id string, code string) (bool, entities.User) {
	user, err := GetUser(id)
	if err != nil {
		return false, entities.User{}
	}

	bioCode := user.BiometricKey
	if bioCode != code {
		return false, entities.User{}
	}

	return true, user
}
