package models

import "time"

type UserRequest struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required,min=4"`
	FirstName string `json:"first_name" validate:"required,min=4"`
	LastName  string `json:"last_name" validate:"required,min=4"`
	Email     string `json:"email" validate:"required,email"`
	Age       int    `json:"age" validate:"required"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}
