package storage

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/muchiri08/crud/types"
)

type ProductsStore interface {
	CreateProduct() error
	GetAllProducts() ([]*types.Product, error)
	DeleteProduct(id int) (int64, error)
	GetProductById(id int) (*types.Product, error)
	UpdateProduct(product *types.Product) error
}

func (s *PostgresStore) CreateProduct(product *types.Product) error {
	query := `INSERT INTO products(name, price) VALUES ($1, $2)`
	_, err := s.db.Exec(query, product.ProductName, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteProduct(id int) (int64, error) {
	res, err := s.db.Exec("DELETE FROM products WHERE id  = $1", id)
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return row, nil
}

func (s *PostgresStore) UpdateProduct(product *types.Product) error {
	query := "UPDATE products SET name = $1, price = $2 WHERE id = $3"
	_, err := s.db.Exec(query, product.ProductName, product.Price, product.Id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) GetProductById(id int) (*types.Product, error) {
	var product = new(types.Product)
	row := s.db.QueryRow("SELECT * FROM products WHERE id = $1", id)
	err := row.Scan(&product.Id, &product.ProductName, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid product id")
		}
		return nil, err
	}
	return product, nil
}

func (s *PostgresStore) GetAllProducts() ([]*types.Product, error) {
	var products []*types.Product
	rows, err := s.db.Query("SELECT * FROM products")
	for rows.Next() {
		product, err := mapRowToProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, err
}

func mapRowToProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(&product.Id, &product.ProductName, &product.Price)
	if err != nil {
		return nil, err
	}
	return product, nil
}
