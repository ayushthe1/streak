package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username" validate:"required,min=2,max=100"`
	// FirstName      string `json:"first_name" validate:"required,min=2,max=100"`
	// LastName       string `json:"last_name" validate:"required,min=2,max=100"`
	// Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"required,min=2"`
	LastActivity   int64  `json:"last_activity"`
	HashedPassword []byte
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	log.Println("Password is hashed and set in user field")
	user.HashedPassword = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
}
