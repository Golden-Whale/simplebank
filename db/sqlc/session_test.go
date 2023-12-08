package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomSession(t *testing.T, user User, refreshToken string) Session {
	arg := CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: [16]byte([]byte(uuid.New().String())),
			Valid: true,
		},
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		IsBlocked:    false,
		ClientIp:     "",
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	session, err := testStore.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)

	require.NotZero(t, user.CreatedAt)
	return session
}

func TestCreateSession(t *testing.T) {
	user := createRandomUser(t)
	createRandomSession(t, user, "refresh_token")
}

func TestGetSession(t *testing.T) {
	user := createRandomUser(t)
	session1 := createRandomSession(t, user, "refresh_token")

	session2, err := testStore.GetSession(context.Background(), session1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)

	require.WithinDuration(t, session1.CreatedAt.Time, session2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, session1.ExpiresAt.Time, session2.ExpiresAt.Time, time.Second)
}
