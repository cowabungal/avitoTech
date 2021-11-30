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

func (s *UserService) ConvertBalance(user *avitoTech.User, currency string) (*avitoTech.User, error) {
	/*var cur struct {
		Cur      float64
	}


		endpoint := fmt.Sprintf("https://free.currconv.com/api/v7/convert?q=%s_RUB&compact=ultra&apiKey=02359d71f5b53497b479", currency)
		resp, err := http.Get(endpoint)
		if err != nil {
		return nil, err
	}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&cur.Cur)
		user.Balance = user.Balance * int(cur.Cur)
*/
		return user, nil
	}
