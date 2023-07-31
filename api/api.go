package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/muchiri08/crud/storage"
	"net/http"
)

type ApiServer struct {
	Address Address
	Store   *storage.PostgresStore
}

type Address struct {
	Host string
	Port string
}

func NewApiServer(address Address, store *storage.PostgresStore) *ApiServer {
	return &ApiServer{
		Address: address,
		Store:   store,
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func changeToHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func (s *ApiServer) Run() {
	mux := mux.NewRouter()
	host := fmt.Sprintf("%s%s", s.Address.Host, s.Address.Port)

	mux.HandleFunc("/add-user", changeToHttpHandlerFunc(s.HandleCreateUsers)).Methods("POST")
	mux.HandleFunc("/users", changeToHttpHandlerFunc(s.HandleGetUsers)).Methods("GET")
	mux.HandleFunc("/delete/{id}", changeToHttpHandlerFunc(s.HandleDeleteUsers)).Methods("DELETE")
	mux.HandleFunc("/update", changeToHttpHandlerFunc(s.HandleUpdateUser)).Methods("PUT")

	fmt.Printf("listening to %s...\n", host)
	panic(http.ListenAndServe(host, mux))

}
