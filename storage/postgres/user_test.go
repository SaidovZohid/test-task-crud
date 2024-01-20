package postgres_test

import (
	"context"
	"testing"

	"github.com/SaidovZohid/test-task-crud/pkg/utils"
	"github.com/SaidovZohid/test-task-crud/storage/repo"
	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T) *repo.User {
	hashedPassword, err := utils.HashPassword("1234567890")
	require.NoError(t, err)
	code, err := utils.GenerateRandomCode(6)
	require.NoError(t, err)
	user, err := dbManager.User().CreateUser(context.Background(), &repo.User{
		Email:    "johndoe" + code + "@gmail.com",
		Password: hashedPassword,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
	return user
}

func deleteUser(t *testing.T, user_id int64) {
	err := dbManager.User().DeleteUser(context.Background(), user_id)
	require.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	user := createUser(t)
	deleteUser(t, user.ID)
}

func TestDeleteUser(t *testing.T) {
	user := createUser(t)
	deleteUser(t, user.ID)
}

func TestGetUserByID(t *testing.T) {
	user := createUser(t)
	u, err := dbManager.User().GetByUserID(context.Background(), user.ID)
	deleteUser(t, u.ID)
	require.NoError(t, err)
	require.NotEmpty(t, u)
}

func TestGetUserByEmail(t *testing.T) {
	user := createUser(t)
	u, err := dbManager.User().GetByEmail(context.Background(), user.Email)
	deleteUser(t, u.ID)
	require.NoError(t, err)
	require.NotEmpty(t, u)
}
