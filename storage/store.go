package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/muchiri08/crud/types"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connString := "user=root dbname=gocrud password=KiNuThiaPro$2 sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) RunMigrationScript() error {
	log.Println("running migration...")

	script, err := os.ReadFile("storage/migration.sql")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(string(script))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (s *PostgresStore) InitAdmin() {
	log.Println("initializing admin...")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345"), 12)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	admin := types.NewUser("admin", "admin@gmail.com", string(passwordHash), "admin")
	s.CreateUser(admin)

}
