package repository

import (
	"avitoTech"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Balance(userId int) (*avitoTech.User, error) {
	var ans avitoTech.User

	query := fmt.Sprintf("SELECT * from %s WHERE user_id=$1", usersTable)
	err := r.db.Get(&ans, query, userId)

	return &ans, err
}

func (r *UserRepository) TopUp(userId int, amount int) (*avitoTech.User, error) {
	var query string
	var ans avitoTech.User

	user, err := r.Balance(userId)
	if err != nil {
		query = fmt.Sprintf("INSERT INTO %s (user_id, balance) values ($1, $2) RETURNING *", usersTable)
		err = r.db.Get(&ans, query, userId, amount)
	} else {
		newBalance := user.Balance + amount
		query = fmt.Sprintf("UPDATE %s SET balance=$1 WHERE user_id=$2 RETURNING *", usersTable)
		err = r.db.Get(&ans, query, newBalance, userId)
	}

	return &ans, err
}

func (r *UserRepository) Debit(userId int, amount int) (*avitoTech.User, error) {
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

	return &ans, err
}

func (r *UserRepository) Transfer(userId int, toId int, amount int) (*avitoTech.User, error) {
	_, err := r.Balance(toId)
	if err != nil {
		return nil, errors.New("the recipient has no balance")
	}

	_, err = r.Debit(userId, amount)
	if err != nil {
		return nil, err
	}

	ans, err := r.TopUp(toId, amount)

	return ans, err
}
