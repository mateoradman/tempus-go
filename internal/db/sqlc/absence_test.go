package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAbsence(t *testing.T) Absence {
	user := createRandomUser(t, nil, nil)

	today := time.Now().UTC()
	end := today.Add(24 * time.Hour)
	arg := CreateAbsenceParams{
		UserID:    user.ID,
		Paid:      false,
		Reason:    "no reason",
		StartTime: today,
		EndTime:   &end,
	}

	absence, err := testQueries.CreateAbsence(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, absence)
	require.NotZero(t, absence.ID)
	require.Equal(t, arg.UserID, absence.UserID)
	require.Equal(t, arg.Paid, absence.Paid)
	require.WithinDuration(t, arg.StartTime, absence.StartTime, 500*time.Millisecond)
	require.WithinDuration(t, *arg.EndTime, *absence.EndTime, 500*time.Millisecond)
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
	require.Equal(t, absence.StartTime, gotAbsence.StartTime)
	require.Equal(t, absence.EndTime, gotAbsence.EndTime)
	require.Equal(t, absence.CreatedAt, gotAbsence.CreatedAt)
	require.Equal(t, absence.UpdatedAt, gotAbsence.UpdatedAt)
}

func TestUpdateAbsence(t *testing.T) {
	absence := createRandomAbsence(t)
	user := createRandomUser(t, nil, nil)
	today := time.Now().UTC()
	end := today.Add(10 * time.Hour)
	arg := UpdateAbsenceParams{
		ID:        absence.ID,
		UserID:    user.ID,
		Paid:      true,
		StartTime: today,
		EndTime:   &end,
	}

	updatedAbsence, err := testQueries.UpdateAbsence(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAbsence)
	require.Equal(t, absence.ID, updatedAbsence.ID)
	require.Equal(t, arg.UserID, updatedAbsence.UserID)
	require.Equal(t, arg.Paid, updatedAbsence.Paid)
	require.WithinDuration(t, arg.StartTime, updatedAbsence.StartTime, time.Second)
	require.WithinDuration(t, *arg.EndTime, *updatedAbsence.EndTime, time.Second)
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
	require.Equal(t, absence.StartTime, deletedAbsence.StartTime)
	require.Equal(t, absence.EndTime, deletedAbsence.EndTime)
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
	user := createRandomUser(t, nil, nil)
	var absences []Absence
	for i := 0; i < 10; i++ {
		today := time.Now().UTC()
		arg := CreateAbsenceParams{
			UserID:    user.ID,
			Paid:      false,
			StartTime: today,
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
