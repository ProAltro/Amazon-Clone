package mysql

import "github.com/ProAltro/Amazon-Clone/entity"

var _ entity.InventoryService = (*InventoryService)(nil)

type InventoryService struct {
	db *DB
}

func NewInventoryService(db *DB) *InventoryService {
	return &InventoryService{db: db}
}

func (service *InventoryService) CreateInventory(inventory *entity.Inventory) (*entity.Inventory, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = createInventory(tx, inventory)
	if err != nil {
		return nil, err
	}

	return inventory, tx.Commit()
}

func (service *InventoryService) FindInventoryBySeller(seller int) (*entity.SellerInventory, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	inventory, err := getInventoryBySeller(tx, seller)
	if err != nil {
		return nil, err
	}

	return inventory, tx.Commit()
}

func (service *InventoryService) FindInventoryByProduct(product int) (*entity.Inventory, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	inventory, err := getInventoryByProduct(tx, product)
	if err != nil {
		return nil, err
	}

	return inventory, tx.Commit()
}

func (service *InventoryService) UpdateInventory(pid, quantity int) (*entity.Inventory, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE inventory SET quantity = ? WHERE product_id = ?", quantity, pid)
	if err != nil {
		return nil, err
	}

	return &entity.Inventory{ProductId: pid, Quantity: quantity}, tx.Commit()
}

func (service *InventoryService) UnlinkInventory(inventory *entity.Inventory) error {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM inventory WHERE product_id = ?", inventory.ProductId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func createInventory(tx *Tx, inventory *entity.Inventory) error {
	_, err := tx.Exec("INSERT INTO inventory (product_id, quantity, seller) VALUES (?, ?, ?)", inventory.ProductId, inventory.Quantity, inventory.Seller)
	return err
}

func getInventoryByProduct(tx *Tx, product int) (*entity.Inventory, error) {
	inventory := entity.Inventory{}
	err := tx.QueryRow("SELECT product_id, quantity, seller FROM inventory WHERE product_id = ?", product).Scan(&inventory.ProductId, &inventory.Quantity, &inventory.Seller)
	if err != nil {
		return nil, err
	}
	return &inventory, nil
}

func getInventoryBySeller(tx *Tx, seller int) (*entity.SellerInventory, error) {
	inventory := entity.SellerInventory{}
	//using a join query to get the inventory of a product, also get product details
	res, err := tx.Query("SELECT inventory.product_id, inventory.quantity, inventory.seller, product.name, product.price, product.description,product.features FROM inventory INNER JOIN product ON inventory.product_id = product.id WHERE inventory.seller_id = ?", seller)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var prod entity.Product
		var quantity int
		var inv entity.Inventory
		err = res.Scan(&inv.ProductId, &quantity, &inv.Seller, &prod.Name, &prod.Price, &prod.Description, &prod.Features)
		if err != nil {
			return nil, err
		}
		inventory.Products = append(inventory.Products, inv)
	}
	return &inventory, nil
}
