package repo

import (
	"context"
	"time"
)

type UserStorageI interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetByUserID(ctx context.Context, userID int64) (*User, error)
	DeleteUser(ctx context.Context, userID int64) error
}

type User struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt time.Time
}
