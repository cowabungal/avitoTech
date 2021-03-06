package repository

import (
	"avitoTech"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"time"
)

func TestUserRepository_Balance(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type mockBehavior func(userId int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   int
		want    *avitoTech.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, userId, 10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input: 1,
			want: &avitoTech.User{
				Id:      1,
				UserId:  1,
				Balance: 10,
			},
		},
		{
			name: "User has no balance",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(0, 0, 0).RowError(0, errors.New("some error"))
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Balance(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_TopUp(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		userId int
		amount float64
	}

	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    *avitoTech.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(2, args.userId, 100)
				got, err := r.Balance(args.userId)
				if err != nil {
					mock.ExpectQuery(fmt.Sprintf("INSERT INTO %s", usersTable)).
						WithArgs(args.userId, args.amount).WillReturnRows(rows)
				} else {
					mock.ExpectQuery(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", usersTable)).
						WithArgs(got.Balance+args.amount, args.userId).WillReturnRows(rows)
				}
				date := time.Now().Format("01-02-2006 15:04:05")
				operation := fmt.Sprintf("Top-up by some by %fRUB", args.amount)
				result := sqlmock.NewResult(0,0)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).WithArgs(args.userId, args.amount, operation, date).WillReturnResult(result)
			},
			input: args{
				userId: 2,
				amount: 100,
			},
			want: &avitoTech.User{
				Id:      2,
				UserId:  2,
				Balance: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.TopUp(tt.input.userId, tt.input.amount, "some by")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Debit(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		userId int
		amount float64
	}

	type mockBehavior func(args args)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    *avitoTech.User
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, args.userId, 10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(args.userId).WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(1, args.userId, 5)
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", usersTable)).
					WithArgs(args.amount, args.userId).WillReturnRows(rows)

				date := time.Now().Format("01-02-2006 15:04:05")
				operation := fmt.Sprintf("Debit by some by %fRUB", args.amount)
				result := sqlmock.NewResult(0,0)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).WithArgs(args.userId, args.amount, operation, date).WillReturnResult(result)

			},
			input: args{
				userId: 1,
				amount: 5,
			},
			want: &avitoTech.User{
				Id:      1,
				UserId:  1,
				Balance: 5,
			},
		},
		{
			name: "User has no balance",
			mock: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(0, 0, 0).RowError(0, errors.New("some error"))
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(args.userId).WillReturnRows(rows)
			},
			input:   args{
				userId: 0,
				amount: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Debit(tt.input.userId, tt.input.amount, "some by")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Transaction(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type mockBehavior func(userId int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   int
		want    *[]avitoTech.Transaction
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).AddRow(1, userId, 1000, "Top-up by bank_card", "date")
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", transactionsTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input: 1,
			want: &[]avitoTech.Transaction{
				{Id: 1,
					UserId:    1,
					Amount:    1000,
					Operation: "Top-up by bank_card",
					Date:      "date",
				}},
		},
		{
			name: "User has no transaction",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).AddRow(0, 0, 0,"","").RowError(0, errors.New("some error"))
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", transactionsTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Transaction(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_OrderByDateTransaction(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type mockBehavior func(userId int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   int
		want    *[]avitoTech.Transaction
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).AddRow(1, userId, 1000, "Top-up by bank_card", "date")
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", transactionsTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input: 1,
			want: &[]avitoTech.Transaction{
				{Id: 1,
					UserId:    1,
					Amount:    1000,
					Operation: "Top-up by bank_card",
					Date:      "date",
				}},
		},
		{
			name: "User has no transaction",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).AddRow(0, 0, 0,"","").RowError(0, errors.New("some error"))
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", transactionsTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.OrderByDateTransaction(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_OrderByAmountTransaction(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type mockBehavior func(userId int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   int
		want    *[]avitoTech.Transaction
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).AddRow(1, userId, 1000, "Top-up by bank_card", "date")
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", transactionsTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input: 1,
			want: &[]avitoTech.Transaction{
				{Id: 1,
					UserId:    1,
					Amount:    1000,
					Operation: "Top-up by bank_card",
					Date:      "date",
				}},
		},
		{
			name: "User has no transaction",
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).AddRow(0, 0, 0,"","").RowError(0, errors.New("some error"))
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", transactionsTable)).
					WithArgs(userId).WillReturnRows(rows)
			},
			input:   0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.OrderByAmountTransaction(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
