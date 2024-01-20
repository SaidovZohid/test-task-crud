package postgres_test

import (
	"context"
	"testing"

	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createBlog(t *testing.T) *repo.Blog {
	user := createUser(t)
	blog, err := dbManager.Blog().CreateBlog(context.TODO(), &repo.Blog{
		Title:   "Facebook",
		Content: "Facebook is stopped working on Meta Project",
		UserID:  user.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, blog)
	return blog
}

func deleteBlog(t *testing.T, blog_id, user_id int64) {
	err := dbManager.Blog().DeleteBlog(context.TODO(), blog_id, user_id)
	require.NoError(t, err)
}

func TestCreateBlog(t *testing.T) {
	blog := createBlog(t)
	deleteBlog(t, blog.ID, blog.UserID)
	require.NotEmpty(t, blog)
	deleteUser(t, blog.UserID)
}

func TestUpdateBlog(t *testing.T) {
	blog := createBlog(t)
	updatedBlog, err := dbManager.Blog().Update(context.TODO(), &repo.Blog{
		ID:      blog.ID,
		Title:   faker.Sentence(),
		Content: faker.Sentence(),
		UserID:  blog.UserID,
	})
	require.NoError(t, err)
	deleteBlog(t, blog.ID, blog.UserID)
	deleteUser(t, blog.UserID)
	require.NotEmpty(t, blog)
	require.NotEmpty(t, updatedBlog)
}

func TestDeleteBlog(t *testing.T) {
	blog := createBlog(t)
	require.NotEmpty(t, blog)
	deleteBlog(t, blog.ID, blog.UserID)
}

func TestGetBlog(t *testing.T) {
	blog := createBlog(t)
	require.NotEmpty(t, blog)
	b, err := dbManager.Blog().GetBlog(context.TODO(), blog.ID)
	require.NoError(t, err)
	require.NotEmpty(t, b)
	deleteBlog(t, blog.ID, blog.UserID)
}

func TestGetAllPosts(t *testing.T) {
	blog := createBlog(t)
	require.NotEmpty(t, blog)
	blogs, err := dbManager.Blog().GetAll(context.TODO(), &repo.GetBlogsParams{
		Limit:      10,
		Page:       1,
		SortByDate: "ASC",
		UserID:     blog.UserID,
		Search:     blog.Title,
	})
	require.GreaterOrEqual(t, len(blogs.Blogs), 1)
	require.NoError(t, err)
	deleteBlog(t, blog.ID, blog.UserID)
}
