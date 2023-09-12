package mysql

import "github.com/ProAltro/Amazon-Clone/entity"

var _ entity.ProductService = (*ProductService)(nil)

type ProductService struct {
	db *DB
}

func NewProductService(db *DB) *ProductService {
	return &ProductService{db: db}
}

func (service *ProductService) CreateProduct(product *entity.Product) (*entity.Product, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = createProduct(tx, product)
	if err != nil {
		return nil, err
	}

	return product, tx.Commit()
}

func (service *ProductService) FindAllProducts() ([]entity.Product, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	products, err := getAllProducts(tx)
	if err != nil {
		return nil, err
	}

	return products, tx.Commit()
}

func (service *ProductService) FindProductByID(id int) (*entity.Product, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	product, err := getProductByID(tx, id)
	if err != nil {
		return nil, err
	}

	return product, tx.Commit()
}

func (service *ProductService) FindProductByName(name string) (*entity.Product, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	product, err := getProductByName(tx, name)
	if err != nil {
		return nil, err
	}

	return product, tx.Commit()
}

func (service *ProductService) FindProductsBySeller(seller int) ([]entity.Product, error) {

	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	products, err := getProductsBySeller(tx, seller)
	if err != nil {
		return nil, err
	}

	return products, tx.Commit()
}

func (service *ProductService) FindProductsByFilter(filter *entity.ProductFilter) ([]entity.Product, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	products, err := getProductsByFilter(tx, filter)
	if err != nil {
		return nil, err
	}

	return products, tx.Commit()
}

func (service *ProductService) UpdateProduct(product *entity.Product) (*entity.Product, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = updateProduct(tx, product)
	if err != nil {
		return nil, err
	}

	return product, tx.Commit()
}

func createProduct(tx *Tx, product *entity.Product) error {
	stmt, err := tx.Prepare("INSERT INTO products(name, description, price, features, seller) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(product.Name, product.Description, product.Price, product.Features, product.Seller)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	product.Id = int(id)

	return nil
}

func getAllProducts(tx *Tx) ([]entity.Product, error) {
	rows, err := tx.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Features, &product.Seller)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func getProductByID(tx *Tx, id int) (*entity.Product, error) {
	row := tx.QueryRow("SELECT * FROM products WHERE id=?", id)
	product := entity.Product{}
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Features, &product.Seller)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func getProductByName(tx *Tx, name string) (*entity.Product, error) {
	row := tx.QueryRow("SELECT * FROM products WHERE name=?", name)
	product := entity.Product{}
	err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Features, &product.Seller)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func getProductsBySeller(tx *Tx, seller int) ([]entity.Product, error) {
	rows, err := tx.Query("SELECT * FROM products WHERE seller=?", seller)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Features, &product.Seller)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func getProductsByFilter(tx *Tx, filter *entity.ProductFilter) ([]entity.Product, error) {
	query := "SELECT id,name,description, price,features,seller FROM products WHERE price BETWEEN ? AND ?"
	if len(filter.Sellers) != 0 {
		query += " AND seller IN ("

		for i := 0; i < len(filter.Sellers); i++ {
			query += "?"
			if i != len(filter.Sellers)-1 {
				query += ","
			}
		}
		query += ")"
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := []interface{}{filter.MinPrice, filter.MaxPrice}
	for _, seller := range filter.Sellers {
		args = append(args, seller)
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Features, &product.Seller)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func updateProduct(tx *Tx, product *entity.Product) error {
	stmt, err := tx.Prepare("UPDATE products SET name=?, description=?, price=?, features=?, seller=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Description, product.Price, product.Features, product.Seller, product.Id)
	if err != nil {
		return err
	}

	return nil
}
