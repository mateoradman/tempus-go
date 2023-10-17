package db

import (
	"context"
	"testing"
	"time"

	"github.com/mateoradman/tempus/internal/types"
	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T, company_id, team_id *int64) User {
	birthDate := time.Date(1999, 2, 3, 0, 0, 0, 0, time.Now().UTC().Location())
	hashedPassword, err := util.HashPassword(util.RandomString(20))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:  util.RandomString(23),
		Email:     util.RandomEmail(),
		Name:      util.RandomString(25),
		Surname:   util.RandomString(25),
		Password:  hashedPassword,
		Gender:    util.RandomGender(),
		Language:  types.German,
		BirthDate: birthDate,
		CompanyID: company_id,
		TeamID:    team_id,
		Timezone:  "Europe/Zagreb",
	}

	user, err := testStore.CreateUser(context.Background(), arg)
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
	require.Nil(t, user.Country)
	require.WithinDuration(t, time.Now().UTC(), user.CreatedAt, 2*time.Second)
	require.Equal(t, arg.Timezone, user.Timezone)

	// test if default values were correctly set
	require.Nil(t, user.UpdatedAt)
	require.Nil(t, user.ManagerID)
	require.Equal(t, user.CompanyID, company_id)
	require.Equal(t, user.TeamID, team_id)

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
	require.Equal(t, user.Timezone, gotUser.Timezone)
	// test if default values were correctly set
	require.Equal(t, user.UpdatedAt, gotUser.UpdatedAt)
	require.Nil(t, user.CompanyID)
	require.Nil(t, user.ManagerID)
	require.Nil(t, user.TeamID)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t, nil, nil)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t, nil, nil)
	gotUser, err := testStore.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	validateGetQuery(t, user, gotUser)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t, nil, nil)
	gotUser, err := testStore.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	validateGetQuery(t, user, gotUser)
}

func TestGetUserByUsername(t *testing.T) {
	user := createRandomUser(t, nil, nil)
	gotUser, err := testStore.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	validateGetQuery(t, user, gotUser)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t, nil, nil)
	expectedLen := 25
	randomString := util.RandomString(expectedLen)
	arg := UpdateUserParams{
		ID:        user.ID,
		Name:      &randomString,
		Surname:   &randomString,
		Gender:    &user.Gender,
		BirthDate: &user.BirthDate,
		Language:  util.Pointer(types.English),
		Country:   user.Country,
	}

	updatedUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, arg.ID, updatedUser.ID)
	require.Equal(t, *arg.Name, updatedUser.Name)
	require.Equal(t, *arg.Surname, updatedUser.Surname)
	require.Equal(t, *arg.Gender, updatedUser.Gender)
	require.Equal(t, *arg.BirthDate, updatedUser.BirthDate)
	require.Equal(t, *arg.Language, updatedUser.Language)
	require.Equal(t, arg.Country, updatedUser.Country)

	// validate times remain unchanged
	require.Equal(t, user.CreatedAt, updatedUser.CreatedAt)
	require.NotNil(t, updatedUser.UpdatedAt)
	require.WithinDuration(t, time.Now(), *updatedUser.UpdatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t, nil, nil)
	deletedUser, err := testStore.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedUser)
	require.Equal(t, user.ID, deletedUser.ID)
	require.Equal(t, user.Username, deletedUser.Username)
	require.Equal(t, user.Email, deletedUser.Email)
	require.Equal(t, user.Name, deletedUser.Name)
	require.Nil(t, deletedUser.UpdatedAt)
	require.Equal(t, user.CreatedAt, deletedUser.CreatedAt)
	require.Equal(t, user.ManagerID, deletedUser.ManagerID)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 users
		createRandomUser(t, nil, nil)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, int(arg.Limit))

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
