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
	TopUp(userId int, amount int) (*avitoTech.User, error)
	Debit(userId int, amount int) (*avitoTech.User, error)
	Transfer(userId int, toId int, amount int) (*avitoTech.User, error)
	ConvertBalance(ans *avitoTech.User, currency string) (*avitoTech.User, error)
}
