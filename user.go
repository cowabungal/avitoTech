package avitoTech

type User struct {
	Id      int `json:"-" db:"id"`
	UserId  int `json:"user_id" db:"user_id"`
	Balance float64 `json:"balance" db:"balance"`
}

type Transaction struct {
	Id      int `json:"transaction_id" db:"id"`
	UserId  int `json:"user_id" db:"user_id"`
	Amount  float64 `json:"amount" db:"amount"`
	Operation string `json:"operation" db:"operation"`
	Date string `json:"date" db:"date"`
}
