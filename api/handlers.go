package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/muchiri08/crud/types"
	"net/http"
	"strconv"
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

func (s *ApiServer) HandleDeleteUsers(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}

	err, row := s.Store.DeleteUser(id)
	if err != nil {
		return err
	}

	if row > 0 {
		msg := fmt.Sprintf("user with id %d deleted successfully", id)
		writeJSON(w, http.StatusOK, msg)
	} else {
		writeJSON(w, http.StatusNotFound, "invalid user")
	}

	return nil
}
