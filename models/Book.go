package models

import (
	"gopkg.in/go-playground/validator.v9"
)

// 'Book' struct (Model)
type Book struct {
	ID          int    `json:"id" validate:"required" example:"1"`
	Title       string `json:"title" validate:"required,min=2,max=100" example:"Heyy! Go"`
	Description string `json:"description" validate:"omitempty,max=500" example:"A book on Golang"`
	// ISBN        string  `json: "isbn"`
	Author Author `json:"author" validate:"required"`
	// Created_at time.Time `json:"created_at"` // CHECK?
	// Updated_at time.Time `json:"updated_at"`
	// Deleted_at time.Time `json: "deleted_at"`
}

func (book *Book) Validate() error {
	err := validate.Struct(book)
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
