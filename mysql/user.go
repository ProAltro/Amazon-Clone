package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/ProAltro/Amazon-Clone/entity"
	"golang.org/x/crypto/bcrypt"
)

var _ entity.UserService = (*UserService)(nil)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (service *UserService) CreateUser(user *entity.User) (*entity.User, error) {
	//check if user email is unique
	if user.Email == "" {
		return nil, errors.New("empty email")
	}
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		fmt.Println("WTF", err)
		return nil, err
	}
	defer tx.Rollback()

	users, err := getAllUsers(tx)

	if err != nil {
		return nil, err
	}

	for _, u := range users {
		if u.Email == user.Email {
			return nil, errors.New("email already exists")
		}
	}
	//pepper is a secret key
	hash, _ := bcrypt.GenerateFromPassword(pepper_pass(user.Password), 10)
	user.Password = string(hash)
	user.DOJ = time.Now()
	err = createUser(tx, user)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit()
}

func (service *UserService) FindAllUsers() ([]entity.User, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users, err := getAllUsers(tx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (service *UserService) FindUserByEmail(email string) (*entity.User, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := getUserByEmail(tx, "email", email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) FindUserByID(id int) (*entity.User, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := getUserByID(tx, "id", id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) AuthenticateUser(email string, password string) (*entity.User, error) {
	tx, err := service.db.BeginTx(nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := getUserByEmail(tx, "email", email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), pepper_pass(password))
	if err != nil {
		return nil, err
	}
	return user, nil

}

func createUser(tx *Tx, user *entity.User) error {
	result, err := tx.Exec("INSERT INTO users (name,email,password,doj,isPrime) VALUES (?,?,?,?,?)", user.Name, user.Email, user.Password, user.DOJ, user.IsPrime)

	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	user.Id = int(id)
	return nil
}

func getAllUsers(tx *Tx) ([]entity.User, error) {
	rows, err := tx.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.DOJ, &user.IsPrime)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getUserByID(tx *Tx, param string, value int) (*entity.User, error) {
	var user entity.User
	var doj entity.NullTime
	err := tx.QueryRow("SELECT id,name,email,password,doj,isPrime FROM users WHERE "+param+"=?", value).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &doj, &user.IsPrime)
	user.DOJ = doj.Time
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func getUserByEmail(tx *Tx, param string, value string) (*entity.User, error) {
	var user entity.User
	var doj entity.NullTime
	err := tx.QueryRow("SELECT id,name,email,password,doj,isPrime FROM users WHERE "+param+"=?", value).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &doj, &user.IsPrime)
	fmt.Println("doj", doj)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func pepper_pass(password string) []byte {
	return []byte(password + "pepper")
}
