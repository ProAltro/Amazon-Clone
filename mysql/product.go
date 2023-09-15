package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ProAltro/Amazon-Clone/entity"
)

var _ entity.ProductService = (*ProductService)(nil) //enforces that the service implements the interface

type ProductService struct {
	db *DB
}

func NewProductService(db *DB) *ProductService {
	return &ProductService{db}
}

func (ps *ProductService) CreateProduct(ctx context.Context, name string, description string, price int, seller string) (*entity.Product, error) {
	tx, err := ps.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO products (name,description,price,seller) VALUES (?,?,?,?)", name, description, price, seller)
	if err != nil {
		return nil, fmt.Errorf("error inserting product: %w", entity.ErrDB)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", entity.ErrDB)
	}
	product := &entity.Product{
		ID:          int(id),
		Name:        name,
		Description: description,
		Price:       price,
		Seller:      seller,
	}

	tx.Commit()
	return product, nil
}

func (ps *ProductService) GetProduct(ctx context.Context, id int) (*entity.Product, error) {
	tx, err := ps.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	product, err := getProduct(tx, id)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return product, nil
}

func (ps *ProductService) GetProducts(ctx context.Context, ids []int) ([]entity.Product, error) {
	tx, err := ps.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	if err != nil {
		return nil, err
	}
	products, err := getProducts(tx, ids)
	if err != nil {
		return nil, err
	}
	if len(products) != len(ids) {
		return nil, fmt.Errorf("some products do not exist: %w", entity.ErrNotFound)
	}
	tx.Commit()
	return products, nil
}

func (ps *ProductService) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	tx, err := ps.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	row, err := tx.Query("SELECT id,name,description,price,seller FROM products")
	if err != nil {
		return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
	}
	products := []entity.Product{}
	for row.Next() {
		product := entity.Product{}
		err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Seller)
		if err != nil {
			return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
		}
		products = append(products, product)
	}

	if err != nil {
		return nil, err
	}

	tx.Commit()
	return products, nil
}

func (ps *ProductService) DeleteProduct(ctx context.Context, id int) error {
	tx, err := ps.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	//check is product exists before deleting
	_, err = getProduct(tx, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error deleting product: %w", entity.ErrDB)
	}

	tx.Commit()
	return nil
}

func getProduct(tx *Tx, id int) (*entity.Product, error) {
	var product entity.Product
	res := tx.QueryRow("SELECT id,name,description,price,seller FROM products WHERE id = ?", id)
	err := res.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Seller)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("product does not exist: %w", entity.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
	}

	return &product, nil
}

func getProducts(tx *Tx, ids []int) ([]entity.Product, error) {
	rows, err := tx.Query("SELECT id,name,description,price,seller FROM products WHERE id IN (?)", ids)
	products := []entity.Product{}
	if err != nil {
		return nil, fmt.Errorf("error querying: %w", entity.ErrDB)
	}

	for rows.Next() {
		product := entity.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Seller)
		if err != nil {
			return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
		}
		products = append(products, product)
	}

	return products, nil
}
