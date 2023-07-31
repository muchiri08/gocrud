package api

import (
	"encoding/json"
	"github.com/muchiri08/crud/types"
	"net/http"
)

func (s *ApiServer) HandleCreateUsers(w http.ResponseWriter, r *http.Request) error {
	user := new(types.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return err
	}

	validatedUser, err := types.ValidateUser(user)
	if err != nil {
		return err
	}
	user = validatedUser
	if err := s.Store.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *ApiServer) HandleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := s.Store.GetAllUsers()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, users)
}
