package service

import (
	"avitoTech"
	"avitoTech/pkg/repository"
	"encoding/json"
	"errors"
	"net/http"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

var (
	purchase = "purchase"
	bankCard = "bank_card"
	transfer = "transfer"
)

func (s *UserService) Balance(userId int) (*avitoTech.User, error) {
	return s.repo.User.Balance(userId)
}

func (s *UserService) TopUp(userId int, amount float64) (*avitoTech.User, error) {
	return s.repo.User.TopUp(userId, amount, bankCard)
}

func (s *UserService) Debit(userId int, amount float64) (*avitoTech.User, error) {
	return s.repo.User.Debit(userId, amount, purchase)
}

func (s *UserService) Transfer(userId int, toId int, amount float64) (*avitoTech.User, error) {
	_, err := s.repo.User.Balance(toId)
	if err != nil {
		return nil, errors.New("the recipient has no balance")
	}

	_, err = s.repo.User.Debit(userId, amount, transfer)
	if err != nil {
		return nil, err
	}

	ans, err := s.repo.User.TopUp(toId, amount, transfer)

	return ans, err
}

func (s *UserService) ConvertBalance(user *avitoTech.User, currency string) (*avitoTech.User, error) {
	type Currency struct {
		Base  string `json:"base"` //base = EUR
		Rates struct {
			Rub float64 `json:"RUB"`
			Usd float64 `json:"USD"`
		} `json:"rates"`
	}

	var cur Currency

	endpoint := "http://api.exchangeratesapi.io/v1/latest?access_key=e532701035ed3f4040b2660e6b7a8a3d"
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&cur)

	//balance in eur, because it is the base in api.exchangeratesapi.io
	beur := user.Balance / cur.Rates.Rub

	switch currency {
	case "USD":
		busd := beur * cur.Rates.Usd
		user.Balance = busd
	case "EUR":
		user.Balance = beur
	}

	return user, err
}

func (s *UserService) Transaction(userId int) (*[]avitoTech.Transaction, error) {
	return s.repo.User.Transaction(userId)
}
