package domain

import "context"

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type UserRepo interface {
	Create(ctx context.Context, user User) (User, error)
	GetById(ctx context.Context, id int64) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	UpdateUser(ctx context.Context, oldUser User) (User, error)
}
