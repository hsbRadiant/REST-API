package models

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// MODELS (for the application) :-
// 'User' struct (Model)
type User struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=30" example:"Harsimran"`
	Email    string `json:"email" validate:"required,email" example:"hshs@gmail.com"`
	Password string `json:"password" validate:"required,min=10" example:"hars123456"`
}

func (user *User) Validate() error {
	err := validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			// fmt.Println("ERROR ERROR")
			return err
		}
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err)
		}
		return err
	}
	return nil
}
