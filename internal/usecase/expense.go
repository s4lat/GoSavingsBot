package usecase

import (
	"context"
	"github.com/s4lat/gosavingsbot/internal/domain"
	"time"
)

type ExpenseUseCase struct {
	*UseCase
}

func NewExpenseUseCase(uc *UseCase) *ExpenseUseCase {
	return &ExpenseUseCase{
		UseCase: uc,
	}
}

func (u *ExpenseUseCase) CreateExpense(ctx context.Context, e domain.Expense) (domain.Expense, error) {
	return u.expenseRepo.CreateExpense(ctx, e)
}

func (u *ExpenseUseCase) GetExpensesByDateAndTimeZone(
	ctx context.Context,
	userId int64,
	date time.Time,
) ([]domain.Expense, error) {
	expenses, err := u.expenseRepo.GetExpensesByDateWithTz(ctx, userId, date)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
