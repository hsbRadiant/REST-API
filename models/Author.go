package models

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

// 'Author' struct (Model)
type Author struct {
	ID   int    `json:"id" validate:"required" example:"1"`
	Name string `json:"name" validate:"required,min=2,max=30" example:"Harsimran"`
	// Book *Book  `json:"book"`
	// Created_at time.Time `json: "created_at, omitempty"`
	// Updated_at time.Time `json: "updated_at, omitempty"`
	// Deleted_at time.Time `json: "deleted_at"`
}

func (author *Author) Validate() error {
	err := validate.Struct(author)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		// for _, err := range err.(validator.ValidationErrors) {
		// 	fmt.Println(err)
		// }
		return err
	}
	return nil
}

// Example function to understand testing in Golang :
// func Example(i float64) (sqr float64) {
// 	sqr = i * i
// 	return
// }
