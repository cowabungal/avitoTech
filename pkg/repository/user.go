package repository

import (
	"avitoTech"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Balance(userId int) (*avitoTech.User, error) {
	var ans avitoTech.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", usersTable)
	err := r.db.Get(&ans, query, userId)

	return &ans, err
}

func (r *UserRepository) TopUp(userId int, amount float64, by string) (*avitoTech.User, error) {
	var query string
	var ans avitoTech.User

	user, err := r.Balance(userId)
	if err != nil {
		query = fmt.Sprintf("INSERT INTO %s (user_id, balance) VALUES ($1, $2) RETURNING *", usersTable)
		err = r.db.Get(&ans, query, userId, amount)
	} else {
		newBalance := user.Balance + amount
		query = fmt.Sprintf("UPDATE %s SET balance=$1 WHERE user_id=$2 RETURNING *", usersTable)
		err = r.db.Get(&ans, query, newBalance, userId)
	}

	if err != nil {
		return &ans, err
	}

	//record transaction
	time := time.Now()
	//Format MM-DD-YYYY hh:mm:ss
	date := time.Format("01-02-2006 15:04:05")
	query = fmt.Sprintf("INSERT INTO %s (user_id, operation, date) VALUES ($1, $2, $3)", transactionsTable)
	operation := fmt.Sprintf("Top-up by %s %fRUB", by, amount)
	_, err = r.db.Exec(query, userId, operation, date)

	return &ans, err
}

func (r *UserRepository) Debit(userId int, amount float64, by string) (*avitoTech.User, error) {
	var ans avitoTech.User

	user, err := r.Balance(userId)
	if err != nil {
		return nil, errors.New("user has no balance")
	}

	newBalance := user.Balance - amount
	if newBalance < 0 {
		return nil, errors.New("insufficient funds")
	}

	query := fmt.Sprintf("UPDATE %s SET balance=$1 WHERE user_id=$2 RETURNING *", usersTable)
	err = r.db.Get(&ans, query, newBalance, userId)

	if err != nil {
		return &ans, err
	}

	//record transaction
	time := time.Now()
	//Format MM-DD-YYYY hh:mm:ss
	date := time.Format("01-02-2006 15:04:05")
	query = fmt.Sprintf("INSERT INTO %s (user_id, operation, date) VALUES ($1, $2, $3)", transactionsTable)
	operation := fmt.Sprintf("Debit by %s %fRUB", by, amount)
	_, err = r.db.Exec(query, userId, operation, date)

	return &ans, err
}

func (r *UserRepository) Transaction(userId int) (*[]avitoTech.Transaction, error) {
	var ans []avitoTech.Transaction

	//sorting from new to old
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 ORDER BY date DESC", transactionsTable)
	err := r.db.Select(&ans, query, userId)

	return &ans, err
}

