package repository

import (
	"avitoTech"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{User: NewUserRepository(db)}
}

type User interface {
	Balance(userId int) (*avitoTech.User, error)
	TopUp(userId int, amount float64) (*avitoTech.User, error)
	Debit(userId int, amount float64) (*avitoTech.User, error)
	Transfer(userId int, toId int, amount float64) (*avitoTech.User, error)
}
