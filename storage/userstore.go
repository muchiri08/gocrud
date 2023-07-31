package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/muchiri08/crud/types"
)

type UserStore interface {
	CreateUser() error
	GetAllUsers() ([]*types.User, error)
	DeleteUser(id int) (error, int64)
	UpdateUser(user *types.User) (int64, error)
}

func (s *PostgresStore) CreateUser(user *types.User) error {
	query := `INSERT INTO users(name, email, password) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteUser(id int) (error, int64) {
	query := `DELETE FROM users WHERE id = $1`
	res, err := s.db.Exec(query, id)
	row, err := res.RowsAffected()
	if err != nil {
		return err, 0
	}

	return nil, row
}

func (s *PostgresStore) UpdateUser(user *types.User) (int64, error) {
	query := `UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4`
	res, err := s.db.Exec(query, user.Name, user.Email, user.Password, user.Id)
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return row, nil
}

func (s *PostgresStore) GetAllUsers() ([]*types.User, error) {
	var users []*types.User
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user, err := mapRowToUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func mapRowToUser(rows *sql.Rows) (*types.User, error) {
	var user = new(types.User)
	err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
