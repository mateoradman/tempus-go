package db

import (
	"context"
	"testing"
	"time"

	"github.com/mateoradman/tempus/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	birthDate := time.Date(1999, 2, 3, 0, 0, 0, 0, time.Now().UTC().Location())
	arg := CreateUserParams{
		Username:  util.RandomString(23),
		Email:     util.RandomEmail(),
		Name:      util.RandomString(25),
		Surname:   util.RandomString(25),
		Password:  util.RandomString(25),
		Gender:    util.RandomGender(),
		BirthDate: birthDate,
		Language:  "hr",
		Country:   "HR",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.Surname, user.Surname)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Gender, user.Gender)
	require.Equal(t, arg.BirthDate, user.BirthDate)
	require.Equal(t, arg.Language, user.Language)
	require.Equal(t, arg.Country, user.Country)
	require.WithinDuration(t, time.Now().UTC(), user.CreatedAt, 2*time.Second)

	// test if default values were correctly set
	require.False(t, user.UpdatedAt.Valid)
	require.False(t, user.CompanyID.Valid)
	require.False(t, user.Timezone.Valid)
	require.False(t, user.ManagerID.Valid)
	require.False(t, user.TeamID.Valid)

	return user
}

func validateGetQuery(t *testing.T, user, gotUser User) {
	require.NotEmpty(t, gotUser)
	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Surname, gotUser.Surname)
	require.Equal(t, user.Gender, gotUser.Gender)
	require.Equal(t, user.BirthDate, gotUser.BirthDate)
	require.Equal(t, user.Language, gotUser.Language)
	require.Equal(t, user.Country, gotUser.Country)
	// test if default values were correctly set
	require.Equal(t, user.UpdatedAt, gotUser.UpdatedAt)
	require.False(t, user.CompanyID.Valid)
	require.False(t, user.Timezone.Valid)
	require.False(t, user.ManagerID.Valid)
	require.False(t, user.TeamID.Valid)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	gotUser, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	validateGetQuery(t, user, gotUser)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)
	gotUser, err := testQueries.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	validateGetQuery(t, user, gotUser)
}

func TestGetUserByUsername(t *testing.T) {
	user := createRandomUser(t)
	gotUser, err := testQueries.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	validateGetQuery(t, user, gotUser)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	manager := createRandomUser(t)
	expectedLen := 25
	arg := UpdateUserParams{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Name:      util.RandomString(expectedLen),
		Surname:   util.RandomString(expectedLen),
		CompanyID: user.CompanyID,
		Gender:    user.Gender,
		BirthDate: user.BirthDate,
		Language:  user.Language,
		Country:   user.Country,
		Timezone:  user.Timezone,
		ManagerID: manager.ManagerID,
		TeamID:    user.TeamID,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.Equal(t, arg.Username, updatedUser.Username)
	require.Equal(t, arg.Email, updatedUser.Email)
	require.Equal(t, arg.Name, updatedUser.Name)
	require.Equal(t, arg.Surname, updatedUser.Surname)
	require.Equal(t, arg.CompanyID, updatedUser.CompanyID)
	require.Equal(t, arg.Gender, updatedUser.Gender)
	require.Equal(t, arg.BirthDate, updatedUser.BirthDate)
	require.Equal(t, arg.Language, updatedUser.Language)
	require.Equal(t, arg.Country, updatedUser.Country)
	require.Equal(t, arg.Timezone, updatedUser.Timezone)
	require.Equal(t, arg.ManagerID, updatedUser.ManagerID)
	require.Equal(t, arg.TeamID, updatedUser.TeamID)

	// validate times remain unchanged
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.True(t, updatedUser.UpdatedAt.Valid)
	require.WithinDuration(t, time.Now(), updatedUser.UpdatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)
	deletedUser, err := testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedUser)
	require.Equal(t, user.ID, deletedUser.ID)
	require.Equal(t, user.Username, deletedUser.Username)
	require.Equal(t, user.Email, deletedUser.Email)
	require.Equal(t, user.Name, deletedUser.Name)
	require.False(t, deletedUser.UpdatedAt.Valid)
	require.Equal(t, user.CreatedAt, deletedUser.CreatedAt)
	require.Equal(t, user.ManagerID, deletedUser.ManagerID)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 users
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, int(arg.Limit))

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
