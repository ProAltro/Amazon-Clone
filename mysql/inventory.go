package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/ProAltro/Amazon-Clone/entity"
)

var _ entity.InventoryService = (*InventoryService)(nil) //enforces that the service implements the interface

type InventoryService struct {
	db *DB
}

func NewInventoryService(db *DB) *InventoryService {
	return &InventoryService{db}
}

func (is *InventoryService) AddStockToInventory(ctx context.Context, id int, quantity int) error {
	tx, err := is.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	//check if product exists
	_, err = getProduct(tx, id)
	if err != nil {
		return err
	}

	_, err = getStock(tx, id)
	if err == nil {
		return fmt.Errorf("stock already exists: %w", entity.ErrConflict)
	} else if !errors.Is(err, entity.ErrNotFound) {
		return fmt.Errorf("error getting stock: %w", entity.ErrDB)
	}

	_, err = tx.Exec("INSERT INTO inventory (product_id,quantity) VALUES (?,?)", id, quantity)
	if err != nil {
		return fmt.Errorf("error inserting stock: %w", entity.ErrDB)
	}

	return nil
}

func (is *InventoryService) UpdateStockInInventory(ctx context.Context, id int, quantity int) error {
	tx, err := is.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	err = updateStockInInventory(tx, id, quantity)
	if err != nil {
		return err
	}

	return nil
}

func (is *InventoryService) GetStockFromInventory(ctx context.Context, id int) (*entity.Stock, error) {
	tx, err := is.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	stock, err := getStock(tx, id)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (is *InventoryService) GetAllStocksFromInventory(ctx context.Context) ([]*entity.Stock, error) {
	tx, err := is.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	stocks, err := getAllStocks(tx)
	if err != nil {
		return nil, err
	}

	return stocks, nil
}

func (is *InventoryService) RemoveStockFromInventory(ctx context.Context, id int) error {
	tx, err := is.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	_, err = getStock(tx, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM inventory WHERE product_id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting stock: %w", entity.ErrDB)
	}

	return nil
}

func getStock(tx *Tx, id int) (*entity.Stock, error) {
	var stock entity.Stock
	//joiniung product and inventory tables to get list of products and their quantities
	row := tx.QueryRow("SELECT p.id,p.name,p.description,p.price,.pseller,i.quantity FROM products p JOIN inventory i ON p.id=i.product_id WHERE p.id=?", id)
	if row.Err() != nil {
		return nil, fmt.Errorf("product does not exist: %w", entity.ErrNotFound)
	}
	err := row.Scan(&stock.Product.ID, &stock.Product.Name, &stock.Product.Description, &stock.Product.Price, &stock.Product.Seller, &stock.Quantity)

	if err != nil {
		return nil, fmt.Errorf("error scanning row: %w", entity.ErrDB)
	}
	return &stock, nil
}

func getAllStocks(tx *Tx) ([]*entity.Stock, error) {
	stocks := []*entity.Stock{}
	rows, err := tx.Query("SELECT p.id,p.name,p.description,p.price,.pseller,i.quantity FROM products p JOIN inventory i ON p.id=i.product_id")
	if err != nil {
		return nil, fmt.Errorf("error scanning row: %w", entity.ErrDB)
	}
	for rows.Next() {
		stock := entity.Stock{}
		err := rows.Scan(&stock.Product.ID, &stock.Product.Name, &stock.Product.Description, &stock.Product.Price, &stock.Product.Seller, &stock.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", entity.ErrDB)
		}
		stocks = append(stocks, &stock)
	}
	return stocks, nil
}

func updateStockInInventory(tx *Tx, id int, quantity int) error {
	if quantity < 0 {
		return fmt.Errorf("quantity cannot be negative: %w", entity.ErrBadData)
	}

	_, err := getStock(tx, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE inventory SET quantity=? WHERE product_id=?", quantity, id)
	if err != nil {
		return fmt.Errorf("error updating stock: %w", entity.ErrDB)
	}

	return nil

}

func removeFromStockInInventory(tx *Tx, id int, quantity int) error {
	stock, err := getStock(tx, id)
	if err != nil {
		return err
	}
	err = updateStockInInventory(tx, id, stock.Quantity-quantity)
	if errors.Is(err, entity.ErrBadData) {
		return fmt.Errorf("not enough stock: %w", entity.ErrConflict)
	} else if err != nil {
		return err
	}

	return nil
}
