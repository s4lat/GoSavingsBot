package usecase

import "github.com/s4lat/gosavingsbot/internal/domain"

type UseCase struct {
	userRepo    domain.UserRepo
	expenseRepo domain.ExpenseRepo

	Expense *ExpenseUseCase
	User    *UserUseCase
}

func NewUseCase(userRepo domain.UserRepo, expenseRepo domain.ExpenseRepo) *UseCase {
	uc := UseCase{
		userRepo:    userRepo,
		expenseRepo: expenseRepo,
	}

	uc.Expense = NewExpenseUseCase(&uc)
	uc.User = NewUserUseCase(&uc)

	return &uc
}
