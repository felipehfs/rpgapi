package models

import "github.com/go-playground/validator"

// Character is litle boring
type Character struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty" validate:"required"`
	Attack  int    `json:"attack,omitempty" validate:"required,gte=0"`
	Defense int    `json:"defense,omitempty" validate:"required,gte=0"`
	Speed   int    `json:"speed,omitempty" validate:"required,gte=0"`
	Life    int    `json:"life,omitempty" validate:"required,gte=0"`
}

// IsValid check if the character is prepared to save
func (c *Character) IsValid() error {
	validate := validator.New()
	return validate.Struct(c)
}
