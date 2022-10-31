package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	user := createRandomUser(t, nil, nil)

	today := time.Now().UTC()
	arg := CreateEntryParams{
		UserID:    user.ID,
		StartTime: today,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, arg.UserID, entry.UserID)
	require.WithinDuration(t, entry.StartTime, arg.StartTime, 500*time.Millisecond)
	require.Nil(t, entry.EndTime)
	require.WithinDuration(t, time.Now(), entry.CreatedAt, time.Second)
	require.Nil(t, entry.UpdatedAt)

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
	require.Equal(t, entry.CreatedAt, gotEntry.CreatedAt)
	require.Equal(t, entry.UpdatedAt, gotEntry.UpdatedAt)
}

func TestUpdateEntry(t *testing.T) {
	entry := createRandomEntry(t)
	user := createRandomUser(t, nil, nil)
	now := time.Now().UTC()
	arg := UpdateEntryParams{
		ID:        entry.ID,
		UserID:    user.ID,
		StartTime: now,
		EndTime:   entry.EndTime,
	}

	updatedEntry, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedEntry)
	require.Equal(t, entry.ID, updatedEntry.ID)
	require.Equal(t, arg.UserID, updatedEntry.UserID)
	require.WithinDuration(t, updatedEntry.StartTime, arg.StartTime, 500*time.Millisecond)
	require.Equal(t, updatedEntry.EndTime, arg.EndTime)
	require.NotNil(t, updatedEntry.UpdatedAt)
	require.WithinDuration(t, time.Now(), *updatedEntry.UpdatedAt, time.Second)
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
	user := createRandomUser(t, nil, nil)
	var entries []Entry
	for i := 0; i < 10; i++ {
		today := time.Now().UTC()
		arg := CreateEntryParams{
			UserID:    user.ID,
			StartTime: today,
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
