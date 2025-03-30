package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(data interface{}) []string {
	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		var msgs []string

		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("%s failed on %s", err.Field(), err.Tag())
			msgs = append(msgs, msg)
		}

		return msgs
	}

	return nil
}
