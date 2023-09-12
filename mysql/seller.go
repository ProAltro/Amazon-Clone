package mysql

import (
	"errors"
	"fmt"

	"github.com/ProAltro/Amazon-Clone/entity"
	"golang.org/x/crypto/bcrypt"
)

var _ entity.SellerService = (*SellerService)(nil)

type SellerService struct {
	db *DB
}

func NewSellerService(db *DB) *SellerService {
	return &SellerService{db: db}
}

func (service *SellerService) CreateSeller(seller *entity.Seller) (*entity.Seller, error) {
	//check if seller email is unique
	if seller.Email == "" {
		return nil, errors.New("empty email")
	}
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()

	sellers, err := getAllSellers(tx)

	if err != nil {
		return nil, err
	}

	for _, s := range sellers {
		if s.Email == seller.Email {
			return nil, errors.New("email already exists")
		}
	}
	//pepper is a secret key
	hash, _ := bcrypt.GenerateFromPassword(pepper_pass(seller.Password), 10)
	seller.Password = string(hash)
	err = createSeller(tx, seller)
	if err != nil {
		return nil, err
	}

	return seller, tx.Commit()
}

func (service *SellerService) FindAllSellers() ([]entity.Seller, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	sellers, err := getAllSellers(tx)
	if err != nil {
		return nil, err
	}
	return sellers, tx.Commit()
}

func (service *SellerService) FindSellerByID(id int) (*entity.Seller, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	seller, err := getSellerByID(tx, id)
	if err != nil {
		return nil, err
	}
	return seller, tx.Commit()
}

func (service *SellerService) FindSellerByEmail(email string) (*entity.Seller, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	seller, err := getSellerByEmail(tx, "email", email)
	if err != nil {
		return nil, err
	}
	return seller, tx.Commit()
}

func (service *SellerService) AuthenticateSeller(email string, password string) (*entity.Seller, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()
	seller, err := getSellerByEmail(tx, "email", email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(seller.Password), pepper_pass(password))
	if err != nil {
		return nil, err
	}
	return seller, tx.Commit()
}

func createSeller(tx *Tx, seller *entity.Seller) error {
	stmt, err := tx.Prepare("INSERT INTO sellers(name, email, address, password) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(seller.Name, seller.Email, seller.Address, seller.Password)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	seller.Id = int(id)
	return err
}

func getSellerByID(tx *Tx, id int) (*entity.Seller, error) {
	row := tx.QueryRow("SELECT * FROM sellers WHERE id=?", id)
	seller := entity.Seller{}
	err := row.Scan(&seller.Id, &seller.Name, &seller.Email, &seller.Address, &seller.Password)
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func getSellerByEmail(tx *Tx, param string, value string) (*entity.Seller, error) {

	row := tx.QueryRow("SELECT * FROM sellers WHERE "+param+"=?", value)
	seller := entity.Seller{}
	err := row.Scan(&seller.Id, &seller.Name, &seller.Email, &seller.Address, &seller.Password)
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

func getAllSellers(tx *Tx) ([]entity.Seller, error) {
	rows, err := tx.Query("SELECT * FROM sellers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sellers []entity.Seller
	for rows.Next() {
		var seller entity.Seller
		err := rows.Scan(&seller.Id, &seller.Name, &seller.Email, &seller.Address, &seller.Password)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, seller)
	}
	return sellers, nil
}
