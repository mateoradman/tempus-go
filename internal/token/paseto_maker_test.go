package token

import (
	"testing"
	"time"

	"github.com/mateoradman/tempus/util"
	"github.com/stretchr/testify/require"
)

func getPasetoMaker(t *testing.T) Maker {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	return maker
}
func TestPasetoMaker(t *testing.T) {
	maker := getPasetoMaker(t)

	username := util.RandomString(25)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload_, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload_)
	require.NotZero(t, payload_.ID)
	require.Equal(t, username, payload_.Username)
	require.WithinDuration(t, issuedAt, payload_.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload_.ExpiredAt, time.Second)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredPasetoToken(t *testing.T) {
	maker := getPasetoMaker(t)

	token, payload, err := maker.CreateToken(util.RandomString(25), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	maker := getPasetoMaker(t)

	token := "fake_token"

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidSymmetricKey(t *testing.T) {
	maker, err := NewPasetoMaker("fake_key")
	require.Nil(t, maker)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid key size")
}
