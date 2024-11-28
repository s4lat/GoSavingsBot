package domain

import (
	"context"
	"time"
)

type Expense struct {
	Id       int64     `json:"id"`
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	Amount   float64   `json:"amount"`
	Currency string    `json:"currency"`
	TypeID   int32     `json:"type_id"`
	UserId   int64     `json:"user_id"`
}

type ExpenseType struct {
	Id        int64  `json:"id"`
	TypeName  string `json:"type_name"`
	TypeColor Color  `json:"type_color"`
	UserId    int64  `json:"user_id"`
}

type Color uint8

const (
	RedColor Color = iota
)

type ExpenseRepo interface {
	CreateExpense(ctx context.Context, expense Expense) (Expense, error)
	UpdateExpense(ctx context.Context, expense Expense) (Expense, error)
	DeleteExpense(ctx context.Context, id int64) error

	GetExpensesCountByDateWithTz(ctx context.Context, userId int64, date time.Time) (int64, error)
	GetExpensesByDateWithTz(ctx context.Context, userId int64, date time.Time) ([]Expense, error)

	CreateExpenseType(ctx context.Context, expenseType ExpenseType) (ExpenseType, error)
	UpdateExpenseType(ctx context.Context, expenseType ExpenseType) (ExpenseType, error)
	DeleteExpenseType(ctx context.Context, id int64, userId int64) error
}
