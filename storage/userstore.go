package storage

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/muchiri08/crud/types"
)

type UserStore interface {
	CreateUser() (*types.User, error)
	GetAllUsers() ([]*types.User, error)
	DeleteUser(id int) (error, int64)
	UpdateUser(user *types.User) (int64, error)
	GetUserByEmail(email string) (*types.User, error)
}

func (s *PostgresStore) CreateUser(user *types.User) (*types.User, error) {
	query := `INSERT INTO users(name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
	row := s.db.QueryRow(query, user.Name, user.Email, user.Password, user.Role)

	err := row.Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
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
	query := `UPDATE users SET name = $1, email = $2, password = $3, role = $4 WHERE id = $5`
	res, err := s.db.Exec(query, user.Name, user.Email, user.Password, user.Role, user.Id)
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

func (s *PostgresStore) GetUserByEmail(email string) (*types.User, error) {
	var user = new(types.User)

	row := s.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	return user, nil
}

func mapRowToUser(rows *sql.Rows) (*types.User, error) {
	var user = new(types.User)
	err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
