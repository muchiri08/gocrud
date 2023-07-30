package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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
