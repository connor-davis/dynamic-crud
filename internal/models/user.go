package models

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	Base
	Name  string `json:"name" gorm:"type:text;not null;" validate:"gte=3;required"`
	Email string `json:"email" gorm:"type:text;unique;not null;" validate:"email;required"`
}

func (u *User) Validate() error {
	validate := validator.New()

	return validate.Struct(u)
}
