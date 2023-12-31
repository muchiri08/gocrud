package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/muchiri08/crud/storage"
	"net/http"
)

var ErrorForbidden = errors.New("forbidden")

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
			if errors.Is(err, ErrorForbidden) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
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

	mux.HandleFunc("/", changeToHttpHandlerFunc(s.HandleLogin)).Methods("POST")
	mux.HandleFunc("/login", changeToHttpHandlerFunc(s.HandleLogin)).Methods("POST")

	users := mux.PathPrefix("/users").Subrouter()
	users.Use(jWTMiddleware)
	users.HandleFunc("/add-user", changeToHttpHandlerFunc(s.HandleCreateUsers)).Methods("POST")
	users.HandleFunc("", changeToHttpHandlerFunc(s.HandleGetUsers)).Methods("GET")
	users.HandleFunc("/delete/{id}", changeToHttpHandlerFunc(s.HandleDeleteUsers)).Methods("DELETE")
	users.HandleFunc("/update", changeToHttpHandlerFunc(s.HandleUpdateUser)).Methods("PUT")

	products := mux.PathPrefix("/products").Subrouter()
	products.Use(jWTMiddleware)
	products.HandleFunc("/add", changeToHttpHandlerFunc(s.CreateProduct)).Methods("POST")
	products.HandleFunc("", changeToHttpHandlerFunc(s.HandleGetAllProducts)).Methods("GET")
	products.HandleFunc("/delete/{id}", changeToHttpHandlerFunc(s.HandleDeleteProduct)).Methods("DELETE")
	products.HandleFunc("/{id}", changeToHttpHandlerFunc(s.HandleGetProductById)).Methods("GET")
	products.HandleFunc("/update", changeToHttpHandlerFunc(s.HandleUpdateProduct)).Methods("PUT")

	fmt.Printf("listening to %s...\n", host)
	panic(http.ListenAndServe(host, mux))

}
