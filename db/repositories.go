package db

import "context"

type UserRepository interface {
	GetUserByUsernameOrEmail(ctx context.Context, username, email string) (User, error)
	CreateUser(ctx context.Context, user User) (User, error)
}
