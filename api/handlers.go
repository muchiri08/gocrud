package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/muchiri08/crud/types"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

const passwordCost = 12

func (s *ApiServer) HandleCreateUsers(w http.ResponseWriter, r *http.Request) error {
	user := new(types.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return err
	}
	//encrypting the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), passwordCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)

	validatedUser, err := types.ValidateUser(user)
	if err != nil {
		return err
	}
	user = validatedUser
	user, err = s.Store.CreateUser(user)
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusCreated, user)

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

func (s *ApiServer) HandleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	user := new(types.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		return err
	}

	user, err := types.ValidateUser(user)
	if err != nil {
		return err
	}

	row, err := s.Store.UpdateUser(user)
	if err != nil {
		return err
	}
	if row > 0 {
		msg := fmt.Sprintf("user with id %d updated successfully", user.Id)
		writeJSON(w, http.StatusOK, msg)
	} else {
		writeJSON(w, http.StatusBadRequest, "invalid user")
	}

	return nil
}

func (s *ApiServer) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	product := new(types.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		return err
	}

	product, err := s.Store.CreateProduct(product)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, product)

	return nil
}

func (s *ApiServer) HandleGetAllProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := s.Store.GetAllProducts()
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, products)

	return nil
}

func (s *ApiServer) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}

	row, err := s.Store.DeleteProduct(id)
	if err != nil {
		return err
	}
	if row > 0 {
		msg := fmt.Sprintf("product with id %d deleted successfully", id)
		writeJSON(w, http.StatusOK, msg)
	} else {
		writeJSON(w, http.StatusNotFound, "invalid product")
	}

	return nil
}

func (s *ApiServer) HandleGetProductById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	product, err := s.Store.GetProductById(id)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, product)

	return nil
}

func (s *ApiServer) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) error {
	product := new(types.Product)
	if err := json.NewDecoder(r.Body).Decode(product); err != nil {
		return err
	}

	row, err := s.Store.UpdateProduct(product)
	if err != nil {
		return err
	}

	if row > 0 {
		msg := fmt.Sprintf("product with id %d updated successfully", product.Id)
		writeJSON(w, http.StatusOK, msg)
	} else {
		writeJSON(w, http.StatusBadRequest, "invalid product id")
	}

	return nil
}

func (s *ApiServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	var userReq = new(types.UserRequest)
	err := json.NewDecoder(r.Body).Decode(userReq)
	if err != nil {
		return err
	}

	user, err := s.Store.GetUserByEmail(userReq.Email)
	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, err.Error())
	}

	ok, err := comparePassword(user.Password, userReq.Password)
	if !ok && err == nil {
		return writeJSON(w, http.StatusUnauthorized, "invalid credentials")
	}

	return writeJSON(w, http.StatusOK, "successfully authenticated")
}

func comparePassword(hashedPass, plainPass string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainPass)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
