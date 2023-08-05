package types

import (
	"errors"
)

type User struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(name, email, password, role string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
}

func ValidateUser(u *User) (*User, error) {
	if u.Email == "" || u.Name == "" || u.Password == "" || u.Role == "" {
		return nil, errors.New("invalid user")
	}
	return u, nil
}
