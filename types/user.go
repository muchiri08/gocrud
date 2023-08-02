package types

import "errors"

type User struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func ValidateUser(u *User) (*User, error) {
	if u.Email == "" || u.Name == "" || u.Password == "" {
		return nil, errors.New("invalid user")
	}
	return u, nil
}
