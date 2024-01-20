package models

import "time"

type Blog struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	User      *BlogUser `json:"user_info"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogUser struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type CreateAndUpdateBlog struct {
	Title   string `json:"title" binding:"required,min=5,max=300"`
	Content string `json:"content" binding:"required,min=20"`
}

type Message struct {
	Msg string `json:"message"`
}

type GetAllBlogsResult struct {
	Blogs []*Blog `json:"blogs"`
	Count int64   `json:"count"`
}