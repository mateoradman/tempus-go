package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Session {
	user := createRandomUser(t, nil, nil)
	arg := CreateSessionParams{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: util.RandomString(32),
		UserAgent:    util.RandomString(20),
		ClientIp:     util.RandomString(20),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, arg.ID, session.ID)
	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.WithinDuration(t, time.Now(), session.CreatedAt, 2*time.Second)

	return session
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session := createRandomSession(t)
	gotSession, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, session.ID, gotSession.ID)
	require.Equal(t, session.Username, gotSession.Username)
	require.Equal(t, session.RefreshToken, gotSession.RefreshToken)
	require.Equal(t, session.UserAgent, gotSession.UserAgent)
	require.Equal(t, session.IsBlocked, gotSession.IsBlocked)
	require.WithinDuration(t, time.Now(), gotSession.CreatedAt, 2*time.Second)
}
