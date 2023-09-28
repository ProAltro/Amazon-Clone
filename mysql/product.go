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

func (ps *ProductService) CreateProduct(ctx context.Context, name string, description string, price int, seller string, images []string) (*entity.Product, error) {
	tx, err := ps.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	uri_images := []string{}
	for _, image := range images {
		uri_images = append(uri_images, "https://pramitpal.me/amazon/api/v1/images/"+image+".jpg")
	}
	json_uri_images, err := entity.Images_To_JSON(uri_images)
	if err != nil {
		return nil, fmt.Errorf("error converting images to json: %w", entity.ErrInternal)
	}
	result, err := tx.Exec("INSERT INTO products (name,description,price,seller,images) VALUES (?,?,?,?,?)", name, description, price, seller, json_uri_images)

	if err != nil {
		fmt.Println(err)
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

	row, err := tx.Query("SELECT id,name,description,price,seller,images FROM products")
	if err != nil {
		return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
	}
	products := []entity.Product{}
	for row.Next() {
		product := entity.Product{}
		images := ""
		err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Seller, &images)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
		}
		product.Images, err = entity.JSON_To_Image([]byte(images))
		if err != nil {
			return nil, fmt.Errorf("error converting json to images: %w", entity.ErrInternal)
		}
		products = append(products, product)
	}

	if err != nil {
		return nil, err
	}

	tx.Commit()
	return products, nil
}

func (ps *ProductService) UpdateProduct(ctx context.Context, update entity.Product) (*entity.Product, error) {
	tx, err := ps.db.BeginTx(nil)

	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	//check is product exists before updating
	prod, err := getProduct(tx, update.ID)
	if err != nil {
		return nil, fmt.Errorf("the product does not exist: %w", entity.ErrNotFound)
	}

	if update.Images != nil {
		uri_images := []string{}
		for _, image := range update.Images {
			uri_images = append(uri_images, "https://pramitpal.me/amazon/api/v1/images/"+image+".jpg")
		}
		update.Images = uri_images
	}
	prod.Update(update)
	images, err := entity.Images_To_JSON(prod.Images)
	_, err = tx.Exec("UPDATE products SET name=?,description=?,price=?,seller=?,images=? WHERE id=?", prod.Name, prod.Description, prod.Price, prod.Seller, images, prod.ID)

	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", entity.ErrDB)
	}

	tx.Commit()
	return prod, nil
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
		return fmt.Errorf("error deleting product: %w", entity.ErrDB)
	}

	tx.Commit()
	return nil
}

func getProduct(tx *Tx, id int) (*entity.Product, error) {
	var product entity.Product
	res := tx.QueryRow("SELECT id,name,description,price,seller,images FROM products WHERE id = ?", id)
	images := ""
	err := res.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Seller, &images)
	product.Images, err = entity.JSON_To_Image([]byte(images))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("product does not exist: %w", entity.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
	}

	return &product, nil
}

func getProducts(tx *Tx, ids []int) ([]entity.Product, error) {
	rows, err := tx.Query("SELECT id,name,description,price,seller,images FROM products WHERE id IN (?)", ids)
	products := []entity.Product{}
	if err != nil {
		return nil, fmt.Errorf("error querying: %w", entity.ErrDB)
	}

	for rows.Next() {
		product := entity.Product{}
		images := ""
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Seller, &images)
		product.Images, err = entity.JSON_To_Image([]byte(images))
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error scanning product: %w", entity.ErrDB)
		}
		products = append(products, product)
	}

	return products, nil
}
