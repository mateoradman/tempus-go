package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAbsence(t *testing.T) Absence {
	user := createRandomUser(t)

	today := time.Now().UTC()
	arg := CreateAbsenceParams{
		UserID: user.ID,
		Paid:   false,
		Date:   today,
		Length: 5.5,
	}

	absence, err := testQueries.CreateAbsence(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, absence)
	require.NotZero(t, absence.ID)
	require.Equal(t, arg.UserID, absence.UserID)
	require.Equal(t, arg.Paid, absence.Paid)
	require.Equal(t, arg.Length, absence.Length)
	require.Equal(t, today.Day(), absence.Date.Day())
	require.Equal(t, today.Month(), absence.Date.Month())
	require.Equal(t, today.Year(), absence.Date.Year())
	require.Nil(t, absence.ApprovedByID)
	require.WithinDuration(t, time.Now(), absence.CreatedAt, 2*time.Second)
	require.Nil(t, absence.UpdatedAt)

	return absence
}

func TestCreateAbsence(t *testing.T) {
	createRandomAbsence(t)
}

func TestGetAbsence(t *testing.T) {
	absence := createRandomAbsence(t)
	gotAbsence, err := testQueries.GetAbsence(context.Background(), absence.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotAbsence)
	require.Equal(t, absence.ID, gotAbsence.ID)
	require.Equal(t, absence.UserID, gotAbsence.UserID)
	require.Equal(t, absence.ApprovedByID, gotAbsence.ApprovedByID)
	require.Equal(t, absence.Paid, gotAbsence.Paid)
	require.Equal(t, absence.Length, gotAbsence.Length)
	require.Equal(t, absence.Date, gotAbsence.Date)
	require.Equal(t, absence.CreatedAt, gotAbsence.CreatedAt)
	require.Equal(t, absence.UpdatedAt, gotAbsence.UpdatedAt)
}

func TestUpdateAbsence(t *testing.T) {
	absence := createRandomAbsence(t)
	user := createRandomUser(t)
	today := time.Now().UTC()
	arg := UpdateAbsenceParams{
		ID:     absence.ID,
		UserID: user.ID,
		Paid:   true,
		Date:   today,
		Length: 10,
	}

	updatedAbsence, err := testQueries.UpdateAbsence(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAbsence)
	require.Equal(t, absence.ID, updatedAbsence.ID)
	require.Equal(t, arg.UserID, updatedAbsence.UserID)
	require.Equal(t, arg.Paid, updatedAbsence.Paid)
	require.Equal(t, today.Day(), absence.Date.Day())
	require.Equal(t, today.Month(), absence.Date.Month())
	require.Equal(t, today.Year(), absence.Date.Year())
	require.Equal(t, arg.Length, updatedAbsence.Length)
	require.NotNil(t, updatedAbsence.UpdatedAt)
	require.WithinDuration(t, time.Now(), *updatedAbsence.UpdatedAt, time.Second)
	require.Equal(t, absence.CreatedAt, updatedAbsence.CreatedAt)
}

func TestDeleteAbsence(t *testing.T) {
	absence := createRandomAbsence(t)
	deletedAbsence, err := testQueries.DeleteAbsence(context.Background(), absence.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedAbsence)
	require.Equal(t, absence.ID, deletedAbsence.ID)
	require.Equal(t, absence.UserID, deletedAbsence.UserID)
	require.Equal(t, absence.Paid, deletedAbsence.Paid)
	require.Equal(t, absence.Date, deletedAbsence.Date)
	require.Equal(t, absence.Length, deletedAbsence.Length)
	require.Equal(t, absence.CreatedAt, deletedAbsence.CreatedAt)
	require.Equal(t, absence.UpdatedAt, deletedAbsence.UpdatedAt)
}

func TestListAbsences(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 absences
		createRandomAbsence(t)
	}

	arg := ListAbsencesParams{
		Limit:  5,
		Offset: 0,
	}

	absences, err := testQueries.ListAbsences(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, absences, int(arg.Limit))

	for _, absence := range absences {
		require.NotEmpty(t, absence)
	}
}

func TestListUserAbsences(t *testing.T) {
	user := createRandomUser(t)
	var absences []Absence
	for i := 0; i < 10; i++ {
		today := time.Now().UTC()
		arg := CreateAbsenceParams{
			UserID: user.ID,
			Paid:   false,
			Date:   today,
			Length: 5.5,
		}
		absence, err := testQueries.CreateAbsence(context.Background(), arg)
		require.NoError(t, err)
		absences = append(absences, absence)
	}
	arg := ListUserAbsencesParams{
		UserID: user.ID,
		Limit:  100,
		Offset: 0,
	}
	userAbsences, err := testQueries.ListUserAbsences(context.Background(), arg)
	require.NoError(t, err)
	require.Subset(t, userAbsences, absences)
}
