package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ProAltro/Amazon-Clone/entity"
	"golang.org/x/crypto/bcrypt"
)

var _ entity.UserService = (*UserService)(nil) //enforces that the service implements the interface

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db}
}

func (us *UserService) CreateUser(ctx context.Context, name string, email string, password string) (*entity.User, error) {
	tx, err := us.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	_, err = getUserByEmail(tx, email)
	if err == nil {
		return nil, fmt.Errorf("user already exists: %w", entity.ErrConflict)
	} else if !errors.Is(err, entity.ErrNotFound) {
		return nil, err
	}

	doj := time.Now().Format("2006-01-02 15:04:05")
	hashedPassword, err := bcrypt.GenerateFromPassword(Pepper_Pass(password), 10)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", entity.ErrServer)
	}
	result, err := tx.Exec("INSERT INTO users (name,email,doj,password) VALUES (?,?,?,?)", name, email, doj, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", entity.ErrDB)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", entity.ErrDB)
	}
	user := &entity.User{
		Id:       int(id),
		Name:     name,
		Email:    email,
		Password: password,
		DOJ:      doj,
	}

	tx.Commit()
	return user, nil

}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	tx, err := us.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	user, err := getUserByEmail(tx, email)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (us *UserService) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	tx, err := us.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	user, err := getUserByID(tx, id)
	if err != nil {
		return nil, err
	}

	tx.Commit()
	return user, nil
}

func (us *UserService) AuthenticateUser(ctx context.Context, email string, password string) (*entity.User, error) {
	tx, err := us.db.BeginTx(nil)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", entity.ErrDB)
	}
	defer tx.Rollback()

	user, err := getUserByEmail(tx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), Pepper_Pass(password))
	if err != nil {
		return nil, fmt.Errorf("password is incorrect: %w", entity.ErrForbidden)
	}

	tx.Commit()
	return user, nil
}

func getUserByEmail(tx *Tx, email string) (*entity.User, error) {
	var user entity.User
	row := tx.QueryRow("SELECT id,name,email,doj,password FROM users WHERE email=?", email)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.DOJ, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user does not exist: %w", entity.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error getting user: %w", entity.ErrDB)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", entity.ErrDB)
	}
	return &user, nil

}

func getUserByID(tx *Tx, id int) (*entity.User, error) {
	var user entity.User
	row := tx.QueryRow("SELECT id,name,email,doj,password FROM users WHERE id=?", id)
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.DOJ, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user does not exist: %w", entity.ErrNotFound)
	} else if err != nil {
		return nil, fmt.Errorf("error getting user: %w", entity.ErrDB)
	}

	return &user, nil

}
