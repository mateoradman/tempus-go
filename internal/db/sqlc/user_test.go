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

	// test if default values were correctly set
	require.False(t, user.CompanyID.Valid)
	require.False(t, user.Timezone.Valid)
	require.False(t, user.ManagerID.Valid)
	require.False(t, user.TeamID.Valid)

	return user
}
