package usecase

type ExpenseUseCase struct {
	*UseCase
}

func NewExpenseUseCase(uc *UseCase) *ExpenseUseCase {
	return &ExpenseUseCase{
		UseCase: uc,
	}
}
