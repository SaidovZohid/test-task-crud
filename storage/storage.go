package storage

import (
	"github.com/SaidovZohid/test-task-crud/storage/postgres"
	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/jmoiron/sqlx"
)

type StorageI interface {
	User() repo.UserStorageI
	Blog() repo.BlogStorageI
}

type Storage struct {
	userRepo repo.UserStorageI
	blogRepo repo.BlogStorageI
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &Storage{
		userRepo: postgres.NewUser(db),
		blogRepo: postgres.NewBlog(db),
	}
}

func (s *Storage) User() repo.UserStorageI {
	return s.userRepo
}

func (s *Storage) Blog() repo.BlogStorageI {
	return s.blogRepo
}
