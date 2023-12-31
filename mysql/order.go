package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ProAltro/Amazon-Clone/entity"
)

var _ entity.OrderService = (*OrderService)(nil) //enforces that the service implements the interface

type OrderService struct {
	db *DB
}

func NewOrderService(db *DB) *OrderService {
	return &OrderService{db}
}

func (os *OrderService) CreateOrder(ctx context.Context, products []entity.Stock, total int) (*entity.Order, error) {
	tx, err := os.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()
	err = createOrder(tx, ctx.Value("uid").(int), products, total)

	if err != nil {
		return nil, err
	}

	tx.Commit()
	return nil, nil
}

func (os *OrderService) GetOrder(ctx context.Context, id int) (*entity.Order, error) {
	tx, err := os.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	order, err := getOrder(tx, id)
	if err != nil {
		return nil, err
	}
	//check if order is of the user
	if ctx.Value("uid") != order.UID {
		return nil, fmt.Errorf("user not authorized: %w", entity.ErrUnauthorized)
	}

	tx.Commit()
	return order, nil
}

func (os *OrderService) GetOrders(ctx context.Context, ids []int) ([]entity.Order, error) {
	tx, err := os.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	orders, err := getOrders(tx, ids)
	filteredOrders := []entity.Order{}
	for i, order := range orders {
		if ctx.Value("uid") != order.UID {
			filteredOrders = append(filteredOrders, orders[i])
		}
	}
	if len(filteredOrders) == 0 {
		return nil, fmt.Errorf("no order with those ids: %w", entity.ErrNotFound)
	}
	orders = filteredOrders
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return orders, nil
}

func (os *OrderService) GetOrdersOfUser(ctx context.Context, uid int) ([]entity.Order, error) {

	tx, err := os.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()
	rows, err := tx.Query("SELECT id,total,products FROM orders WHERE user_id=?", uid)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %w", entity.ErrDB)
	}
	defer rows.Close()
	orders := []entity.Order{}
	for rows.Next() {
		order := entity.Order{UID: uid}
		json_products := ""
		err := rows.Scan(&order.ID, &order.Total, &json_products)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", entity.ErrDB)
		}
		products, err := entity.FromJSON([]byte(json_products))
		if err != nil {
			return nil, fmt.Errorf("error converting json to products: %w", entity.ErrInternal)
		}
		order.Products = products

		orders = append(orders, order)
	}

	tx.Commit()
	return orders, nil
}

func (os *OrderService) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	tx, err := os.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id,total,uid,products FROM orders")
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %w", entity.ErrDB)
	}
	defer rows.Close()
	orders := []entity.Order{}
	for rows.Next() {
		order := entity.Order{}
		err := rows.Scan(&order.ID, &order.Total, &order.UID, &order.Products)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", entity.ErrDB)
		}
		orders = append(orders, order)
	}

	tx.Commit()
	return orders, nil
}

func (os *OrderService) DeleteOrder(ctx context.Context, id int) error {
	tx, err := os.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()
	//get order
	order, err := getOrder(tx, id)
	if err != nil {
		return err
	}
	//check if order is of the user
	if ctx.Value("uid") != order.UID {
		return fmt.Errorf("user not authorized: %w", entity.ErrUnauthorized)
	}

	_, err = tx.Exec("DELETE FROM orders WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting order: %w", entity.ErrDB)
	}

	tx.Commit()
	return nil
}

func createOrder(tx *Tx, uid int, products []entity.Stock, total int) error {
	//insert into orders table
	//convert products to json
	//insert into orders table
	json_products, err := entity.ToJSON(products)
	if err != nil {
		return fmt.Errorf("error converting products to json: %w", entity.ErrInternal)
	}

	_, err = tx.Exec("INSERT INTO orders (total,user_id,products) VALUES (?,?,?)", total, uid, json_products)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("error inserting into orders table: %w", entity.ErrDB)
	}
	return nil
}

func getOrder(tx *Tx, id int) (*entity.Order, error) {
	result := tx.QueryRow("SELECT id,total,user_id,products FROM orders WHERE id=?", id)
	order := &entity.Order{}
	json_products := ""
	err := result.Scan(&order.ID, &order.Total, &order.UID, &json_products)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("order does not exist: %w", entity.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning order: %w", entity.ErrDB)
	}

	products, err := entity.FromJSON([]byte(json_products))
	if err != nil {
		return nil, fmt.Errorf("error converting json to products: %w", entity.ErrInternal)
	}
	order.Products = products

	return order, nil
}

func getOrders(tx *Tx, ids []int) ([]entity.Order, error) {
	rows, err := tx.Query("SELECT id,total,user_idid,products FROM orders WHERE id IN (?)", ids)
	if err != nil {
		return nil, fmt.Errorf("error getting orders: %w", entity.ErrDB)
	}
	defer rows.Close()
	orders := []entity.Order{}
	for rows.Next() {
		order := entity.Order{}
		json_products := ""
		err := rows.Scan(&order.ID, &order.Total, &order.UID, &json_products)
		if err != nil {
			return nil, fmt.Errorf("error scanning order: %w", entity.ErrDB)
		}
		products, err := entity.FromJSON([]byte(json_products))
		if err != nil {
			return nil, fmt.Errorf("error converting json to products: %w", entity.ErrInternal)
		}
		order.Products = products

		orders = append(orders, order)
	}
	return orders, nil
}
