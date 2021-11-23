package service

import (
	"avitoTech"
	"avitoTech/pkg/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Balance(userId int) (*avitoTech.User, error) {
	return s.repo.User.Balance(userId)
}

func (s *UserService) TopUp(userId int, amount int) (*avitoTech.User, error) {
	return s.repo.User.TopUp(userId, amount)
}

func (s *UserService) Debit(userId int, amount int) (*avitoTech.User, error) {
	return s.repo.User.Debit(userId, amount)
}

func (s *UserService) Transfer(userId int, toId int, amount int) (*avitoTech.User, error) {
	return s.repo.User.Transfer(userId, toId, amount)
}
