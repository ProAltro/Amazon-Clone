package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ProAltro/Amazon-Clone/entity"
)

var _ entity.CartService = (*CartService)(nil) //enforces that the service implements the interface

type CartService struct {
	db *DB
}

func NewCartService(db *DB) *CartService {
	return &CartService{db}
}

func (cs *CartService) AddProductToCart(ctx context.Context, pid int, quantity int) error {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	stock, err := getStock(tx, pid)
	if err != nil {
		return err
	}

	_, err = getCartItem(tx, uid, pid)

	if err == nil {
		if stock.Quantity+quantity > 10 {
			return fmt.Errorf("quantity exceeds stock: %w", entity.ErrConflict)
		}
		_, err = tx.Exec("UPDATE cart SET quantity=quantity+? WHERE user_id=? AND product_id=?", quantity, uid, pid)
	} else if errors.Is(err, entity.ErrNotFound) {
		if quantity > stock.Quantity {
			return fmt.Errorf("quantity exceeds stock: %w", entity.ErrConflict)
		}
		_, err = tx.Exec("INSERT INTO cart (user_id,product_id,quantity) VALUES (?,?,?)", uid, pid, quantity)
	} else {
		return err
	}

	if err != nil {
		return fmt.Errorf("error inserting cart item: %w", entity.ErrDB)
	}

	tx.Commit()
	return nil
}

func (cs *CartService) RemoveProductFromCart(ctx context.Context, pid int) error {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	_, err = getCartItem(tx, uid, pid)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM cart WHERE user_id=? AND product_id=?", uid, pid)
	if err != nil {
		return fmt.Errorf("error deleting cart item: %w", entity.ErrDB)
	}
	tx.Commit()
	return nil
}

func (cs *CartService) GetCart(ctx context.Context) (*entity.Cart, error) {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	cart, err := getCart(tx, uid)
	if err != nil {
		return nil, err
	}

	total, err := total(tx, uid)
	if err != nil {
		return nil, err
	}

	cart.Total = total
	tx.Commit()
	return cart, nil
}

func (cs *CartService) ModifyCart(ctx context.Context, pid int, quantity int) (cart *entity.Cart, err error) {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	_, err = getCartItem(tx, uid, pid)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("UPDATE cart SET quantity=? WHERE user_id=? AND product_id=?", quantity, uid, pid)
	if err != nil {
		return nil, fmt.Errorf("error modifying cart: %w", entity.ErrDB)
	}
	tx.Commit()
	return cs.GetCart(ctx)
}

func (cs *CartService) Total(ctx context.Context) (int, error) {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return 0, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	total, err := total(tx, uid)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (cs *CartService) ClearCart(ctx context.Context) error {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM cart WHERE user_id=?", uid)
	if err != nil {
		return fmt.Errorf("error clearing cart: %w", entity.ErrDB)
	}
	tx.Commit()
	return nil
}

func (cs *CartService) Checkout(ctx context.Context) error {
	uid := ctx.Value("uid").(int)
	tx, err := cs.db.BeginTx(nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	cart, err := getCart(tx, uid)
	if err != nil {
		return err
	}
	notEnoughStock := []entity.Stock{}
	for _, stock := range cart.Products {
		err = removeFromStockInInventory(tx, stock.Product.ID, stock.Quantity)
		if err == nil {
			continue
		} else if errors.Is(err, entity.ErrConflict) {
			notEnoughStock = append(notEnoughStock, stock)
		} else {
			return fmt.Errorf("error updating inventory: %w", entity.ErrDB)
		}
	}

	if len(notEnoughStock) > 0 {
		tx.Rollback()
		return fmt.Errorf("not enough stock of %v: %w", notEnoughStock, entity.ErrConflict)
	}
	cart.Total, err = total(tx, uid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = createOrder(tx, uid, cart.Products, cart.Total)
	if err != nil {
		tx.Rollback()
		return err
	}

	cs.ClearCart(ctx)

	tx.Commit()

	return nil
}

func getCartItem(tx *Tx, uid int, pid int) (*entity.Stock, error) {
	var stock entity.Stock
	row := tx.QueryRow("SELECT p.id,p.name,p.description,p.price,p.seller,c.quantity FROM products p JOIN cart c ON p.id=c.product_id WHERE c.user_id=? AND c.product_id=?", uid, pid)

	err := row.Scan(&stock.Product.ID, &stock.Product.Name, &stock.Product.Description, &stock.Product.Price, &stock.Product.Seller, &stock.Quantity)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("cart item does not exist: %w", entity.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning stock: %w", entity.ErrDB)
	}

	return &stock, nil
}

func getCart(tx *Tx, uid int) (*entity.Cart, error) {
	rows, err := tx.Query("SELECT p.id,p.name,p.description,p.price,p.seller,c.quantity FROM products p JOIN cart c ON p.id=c.product_id WHERE c.user_id=?", uid)
	if err != nil {
		return nil, fmt.Errorf("error getting cart: %w", entity.ErrDB)
	}
	defer rows.Close()

	var cart entity.Cart
	cart.UID = uid
	for rows.Next() {
		var stock entity.Stock
		err := rows.Scan(&stock.Product.ID, &stock.Product.Name, &stock.Product.Description, &stock.Product.Price, &stock.Product.Seller, &stock.Quantity)
		if err != nil {
			return nil, fmt.Errorf("error scanning cart: %w", entity.ErrDB)
		}
		cart.Products = append(cart.Products, stock)
	}

	return &cart, nil
}

func total(tx *Tx, uid int) (int, error) {

	rows, err := tx.Query("SELECT p.price,c.quantity FROM products p JOIN cart c ON p.id=c.product_id WHERE c.user_id=?", uid)
	if err != nil {
		return 0, fmt.Errorf("error getting cart: %w", entity.ErrDB)
	}
	defer rows.Close()

	var total int
	for rows.Next() {
		var price, quantity int
		err := rows.Scan(&price, &quantity)
		if err != nil {
			return 0, fmt.Errorf("error scanning cart: %w", entity.ErrDB)
		}
		total += price * quantity
	}

	return total, nil
}
