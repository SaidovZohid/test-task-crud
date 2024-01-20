package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/jmoiron/sqlx"
)

type blogRepo struct {
	db *sqlx.DB
}

func NewBlog(db *sqlx.DB) repo.BlogStorageI {
	return &blogRepo{
		db: db,
	}
}

func (b *blogRepo) CreateBlog(ctx context.Context, blog *repo.Blog) (*repo.Blog, error) {
	if err := b.db.QueryRow("INSERT INTO blogs(title, content, user_id) VALUES ($1, $2, $3) RETURNING id, created_at", blog.Title, blog.Content, blog.UserID).Scan(&blog.ID, &blog.CreatedAt); err != nil {
		return nil, err
	}
	return blog, nil
}

func (b *blogRepo) GetBlog(ctx context.Context, blogID int64) (*repo.Blog, error) {
	var blog repo.Blog
	if err := b.db.QueryRow("SELECT id, title, content, user_id, created_at FROM blogs WHERE id = $1", blogID).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Content,
		&blog.UserID,
		&blog.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &blog, nil
}

func (b *blogRepo) Update(ctx context.Context, blog *repo.Blog) (*repo.Blog, error) {
	query := `
		UPDATE blogs SET
			title = $1,
			content = $2
	    WHERE id = $3 AND user_id = $4 RETURNING created_at
	`

	err := b.db.QueryRow(
		query,
		blog.Title,
		blog.Content,
		blog.ID,
		blog.UserID,
	).Scan(&blog.CreatedAt)
	if err != nil {
		return nil, err
	}
	return blog, nil
}

func (b *blogRepo) DeleteBlog(ctx context.Context, blogID, userID int64) error {
	row, err := b.db.Exec(
		"DELETE FROM blogs WHERE id = $1 AND user_id = $2",
		blogID,
		userID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (b *blogRepo) GetAll(ctx context.Context, params *repo.GetBlogsParams) (*repo.GetAllBlogsResult, error) {
	result := repo.GetAllBlogsResult{
		Blogs: make([]*repo.Blog, 0),
	}

	offset := (params.Page - 1) * params.Limit

	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)

	filter := " WHERE true "
	if params.Search != "" {
		srch := "%" + params.Search + "%"
		filter += fmt.Sprintf(" AND content ILIKE '%s' AND title ILIKE '%s'", srch, srch)
	}

	if params.UserID != 0 {
		filter += fmt.Sprintf(" AND user_id = %d", params.UserID)
	}

	orderBy := " ORDER BY created_at DESC"
	if params.SortByDate != "" {
		orderBy = fmt.Sprintf(" ORDER BY created_at %s", params.SortByDate)
	}

	query := `
		SELECT
			id,
			title,
			content,
			user_id,
			created_at
		FROM blogs
	` + filter + orderBy + limit

	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var blog repo.Blog
		err := rows.Scan(
			&blog.ID,
			&blog.Title,
			&blog.Content,
			&blog.UserID,
			&blog.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Blogs = append(result.Blogs, &blog)
	}

	queryCount := "SELECT count(1) FROM blogs " + filter

	err = b.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
