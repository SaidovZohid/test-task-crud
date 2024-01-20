package repo

import (
	"context"
	"time"
)

type BlogStorageI interface {
	CreateBlog(ctx context.Context, blog *Blog) (*Blog, error)
	GetBlog(ctx context.Context, blogID int64) (*Blog, error)
	DeleteBlog(ctx context.Context, blogID, userID int64) error
	GetAll(ctx context.Context, params *GetBlogsParams) (*GetAllBlogsResult, error)
	Update(ctx context.Context, blog *Blog) (*Blog, error)
}

type Blog struct {
	ID        int64
	Title     string
	Content   string
	UserID    int64
	CreatedAt time.Time
}

type GetBlogsParams struct {
	Limit      int64
	Page       int64
	Search     string
	UserID     int64
	SortByDate string
}

type GetAllBlogsResult struct {
	Blogs []*Blog
	Count int64
}
