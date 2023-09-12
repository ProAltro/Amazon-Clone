package mysql

import (
	"fmt"

	"github.com/ProAltro/Amazon-Clone/entity"
)

var _ entity.OrderService = (*OrderService)(nil)

type OrderService struct {
	db *DB
}

func NewOrderService(db *DB) *OrderService {
	return &OrderService{db: db}
}

func (service *OrderService) CreateOrder(order *entity.Order) (*entity.Order, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = createOrder(tx, order)
	if err != nil {
		return nil, err
	}

	return order, tx.Commit()
}

func (service *OrderService) FindAllOrders() ([]entity.Order, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	orders, err := getAllOrders(tx)
	if err != nil {
		return nil, err
	}
	return orders, tx.Commit()
}

func (service *OrderService) FindOrderById(id int) (*entity.Order, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	order, err := getOrderById(tx, id)
	if err != nil {
		return nil, err
	}
	return order, tx.Commit()
}

func (service *OrderService) FindOrdersByUser(user int) ([]entity.Order, error) {

	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	orders, err := getOrdersByUser(tx, user)
	if err != nil {
		return nil, err
	}
	return orders, tx.Commit()
}

func createOrder(tx *Tx, order *entity.Order) error {
	stmt, err := tx.Prepare("INSERT INTO orders(user, total_price, date) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(order.User, order.TotalPrice, order.Date)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	order.Id = int(id)

	return nil
}

func getAllOrders(tx *Tx) ([]entity.Order, error) {
	rows, err := tx.Query("SELECT id, user, total_price, date FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []entity.Order{}

	for rows.Next() {
		var order entity.Order
		var date entity.NullTime
		err := rows.Scan(&order.Id, &order.User, &order.TotalPrice, &date)
		order.Date = date.Time.String()
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func getOrderById(tx *Tx, id int) (*entity.Order, error) {
	row := tx.QueryRow("SELECT id, user, total_price, date FROM orders WHERE id=?", id)
	order := entity.Order{}
	var date entity.NullTime
	err := row.Scan(&order.Id, &order.User, &order.TotalPrice, &date)
	order.Date = date.Time.String()
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func getOrdersByUser(tx *Tx, user int) ([]entity.Order, error) {
	rows, err := tx.Query("SELECT id, user, total_price, date FROM orders WHERE user=?", user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []entity.Order{}

	for rows.Next() {
		var order entity.Order
		var date entity.NullTime
		err := rows.Scan(&order.Id, &order.User, &order.TotalPrice, &date)
		order.Date = date.Time.String()
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
