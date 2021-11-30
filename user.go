package avitoTech

type User struct {
	Id      int `json:"-" db:"id"`
	UserId  int `json:"user_id" db:"user_id"`
	Balance int `json:"balance" db:"balance"`
}
