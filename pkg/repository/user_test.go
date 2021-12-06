package repository

import (
	"avitoTech"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
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

			got, err := r.TopUp(tt.input.userId, tt.input.amount)
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

/*
func TestUserRepository_Debit(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		userId int
		amount int
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
				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(3, args.userId, 100)
				mock.ExpectQuery(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+)", usersTable)).
					WithArgs(args.amount, args.userId).WillReturnRows(rows)
			},
			input: args{
				userId: 3,
				amount: 50,
			},
			want: &avitoTech.User{
				Id:      3,
				UserId:  3,
				Balance: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Debit(tt.input.userId, tt.input.amount)
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
*/
