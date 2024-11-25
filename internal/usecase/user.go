package usecase

import (
	"context"
	"github.com/s4lat/gosavingsbot/internal/domain"
)

type UserUseCase struct {
	*UseCase
}

func NewUserUseCase(uc *UseCase) *UserUseCase {
	return &UserUseCase{
		UseCase: uc,
	}
}

func (u *UserUseCase) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	user, err := u.userRepo.GetById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserUseCase) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	user, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	user, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
