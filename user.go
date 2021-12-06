package avitoTech

type User struct {
	Id      int `json:"-" db:"id"`
	UserId  int `json:"user_id" db:"user_id"`
	Balance float64 `json:"balance" db:"balance"`
}

type Transaction struct {
	Id      int `json:"-" db:"id"`
	UserId  int `json:"user_id" db:"user_id"`
	Operation string `json:"operation" db:"operation"`
}
