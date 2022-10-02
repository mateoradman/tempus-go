package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	user := createRandomUser(t)

	today := time.Now().UTC()
	arg := CreateEntryParams{
		UserID:    user.ID,
		StartTime: today,
		Date:      today,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, arg.UserID, entry.UserID)
	require.Equal(t, arg.StartTime, entry.StartTime)
	require.Equal(t, today.Day(), entry.Date.Day())
	require.Equal(t, today.Month(), entry.Date.Month())
	require.Equal(t, today.Year(), entry.Date.Year())
	require.False(t, entry.EndTime.Valid)
	require.WithinDuration(t, time.Now(), entry.CreatedAt, 2*time.Second)
	require.False(t, entry.UpdatedAt.Valid)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)
	gotEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotEntry)
	require.Equal(t, entry.ID, gotEntry.ID)
	require.Equal(t, entry.UserID, gotEntry.UserID)
	require.Equal(t, entry.StartTime, gotEntry.StartTime)
	require.Equal(t, entry.EndTime, gotEntry.EndTime)
	require.Equal(t, entry.Date, gotEntry.Date)
	require.Equal(t, entry.CreatedAt, gotEntry.CreatedAt)
	require.Equal(t, entry.UpdatedAt, gotEntry.UpdatedAt)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)
	user := createRandomUser(t)
	arg := UpdateEntryParams{
		ID:        entry.ID,
		UserID:    user.ID,
		StartTime: time.Now().UTC(),
		EndTime:   entry.EndTime,
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)
	require.Equal(t, entry.ID, updatedEntry.ID)
	require.Equal(t, arg.UserID, updatedEntry.UserID)
	require.Equal(t, arg.StartTime, updatedEntry.StartTime)
	require.Equal(t, entry.EndTime, updatedEntry.EndTime)
	require.True(t, updatedEntry.UpdatedAt.Valid)
	require.WithinDuration(t, time.Now(), updatedEntry.UpdatedAt.Time, time.Second)
	require.Equal(t, entry.CreatedAt, updatedEntry.CreatedAt)
}

func TestDeleteEntry(t *testing.T) {
	entry := createRandomEntry(t)
	deletedEntry, err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedEntry)
	require.Equal(t, entry.ID, deletedEntry.ID)
	require.Equal(t, entry.UserID, deletedEntry.UserID)
	require.Equal(t, entry.StartTime, deletedEntry.StartTime)
	require.Equal(t, entry.EndTime, deletedEntry.EndTime)
	require.Equal(t, entry.Date, deletedEntry.Date)
	require.Equal(t, entry.CreatedAt, deletedEntry.CreatedAt)
	require.Equal(t, entry.UpdatedAt, deletedEntry.UpdatedAt)
}

func TestListEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 entries
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 0,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, int(arg.Limit))

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

func TestListUserEntries(t *testing.T) {
	user := createRandomUser(t)
	var entries []Entry
	for i := 0; i < 10; i++ {
		today := time.Now().UTC()
		arg := CreateEntryParams{
			UserID:    user.ID,
			StartTime: today,
			Date:      today,
		}
		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		entries = append(entries, entry)
	}
	arg := ListUserEntriesParams{
		UserID: user.ID,
		Limit:  100,
		Offset: 0,
	}
	userEntries, err := testQueries.ListUserEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Subset(t, userEntries, entries)
}
