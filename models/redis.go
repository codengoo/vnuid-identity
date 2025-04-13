package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"vnuid-identity/databases"

	"github.com/google/uuid"
)

type Login2FaAcceptConfig struct {
	SaveDevice bool   `json:"save_device"`
	DeviceID   string `json:"device_id"`
	UID        string `json:"uid"`
}

type LoginByCode2FaConfig struct {
	Code     int    `json:"code" validate:"required"`
	Session  string `json:"session" validate:"required"`
	DeviceID string `json:"device_id"`
}

var LOGIN_KEY = "qr:login:"
var LOGIN_CODE_KEY = "qr:login:code:"
var LOGIN_ACCEPT_KEY = "qr:login:accept:"

func SetJsonValue(key string, value interface{}) error {
	bgctx := context.Background()
	jsonSession, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = databases.RD.Set(bgctx, key, jsonSession, 60*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetJsonValue(key string, value interface{}) error {
	bgctx := context.Background()
	val, err := databases.RD.Get(bgctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(val), value)
	if err != nil {
		return err
	}

	return nil
}

func SetLoginSession(uid string) (string, error) {
	bgctx := context.Background()
	session := uuid.New().String()

	key := fmt.Sprintf("%s%s", LOGIN_KEY, session)
	if err := databases.RD.Set(bgctx, key, uid, 60*time.Second).Err(); err != nil {
		return "", err
	}

	return session, nil
}

func SetLoginCodeSession(uid string, value LoginByCode2FaConfig) error {
	key := fmt.Sprintf("%s%s", LOGIN_CODE_KEY, uid)
	return SetJsonValue(key, value)
}

func SetLoginAcceptSession(session string, value Login2FaAcceptConfig) error {
	var LOGIN_ACCEPT_KEY = "qr:login:accept:"
	key := fmt.Sprintf("%s%s", LOGIN_ACCEPT_KEY, session)

	return SetJsonValue(key, value)
}

func GetLoginSession(session string) (string, error) {
	key := fmt.Sprintf("%s%s", LOGIN_KEY, session)

	uid, err := databases.RD.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	} else {
		return uid, nil
	}
}

func GetLoginAcceptSession(session string) (Login2FaAcceptConfig, error) {
	key := fmt.Sprintf("%s%s", LOGIN_ACCEPT_KEY, session)

	var config Login2FaAcceptConfig
	if err := GetJsonValue(key, &config); err != nil {
		return Login2FaAcceptConfig{}, err
	} else {
		return config, nil
	}
}

func GetLoginCodeSession(uid string) (LoginByCode2FaConfig, error) {
	key := fmt.Sprintf("%s%s", LOGIN_CODE_KEY, uid)

	var config LoginByCode2FaConfig
	if err := GetJsonValue(key, &config); err != nil {
		return LoginByCode2FaConfig{}, err
	} else {
		return config, nil
	}
}
