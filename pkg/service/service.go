package service

import (
	"avitoTech"
	"avitoTech/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{User: NewUserService(repo)}
}

type User interface {
	Balance(userId int) (*avitoTech.User, error)
	TopUp(userId int, amount float64) (*avitoTech.User, error)
	Debit(userId int, amount float64) (*avitoTech.User, error)
	Transfer(userId int, toId int, amount float64) (*avitoTech.User, error)
	ConvertBalance(user *avitoTech.User, currency string) (*avitoTech.User, error)
	Transaction(userId int) (*[]avitoTech.Transaction, error)
}
