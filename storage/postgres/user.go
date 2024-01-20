package postgres

import (
	"context"
	"database/sql"

	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) repo.UserStorageI {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*repo.User, error) {
	query := `
		SELECT 
			id,
			email,
			password_hash,
			created_at
		FROM users WHERE email = $1
	`
	var user repo.User
	err := u.db.QueryRow(
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) GetByUserID(ctx context.Context, userID int64) (*repo.User, error) {
	query := `
		SELECT 
			id,
			email,
			created_at
		FROM users WHERE id = $1
	`
	var user repo.User
	err := u.db.QueryRow(
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepo) CreateUser(ctx context.Context, user *repo.User) (*repo.User, error) {
	query := `
		INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at
	`
	if err := u.db.QueryRow(
		query,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, userID int64) error {
	res, err := u.db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
