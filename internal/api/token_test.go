package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	mockdb "github.com/mateoradman/tempus/internal/db/mock"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
)

func requireBodyMatchRefreshToken(t *testing.T, body *bytes.Buffer, server *Server) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotResponse refreshTokenResponse
	err = json.Unmarshal(data, &gotResponse)
	require.NoError(t, err)

	require.NotEmpty(t, gotResponse.AccessToken)
	require.WithinDuration(
		t,
		time.Now().Add(server.config.AccessTokenDuration),
		gotResponse.AccessTokenExpiresAt,
		2*time.Second,
	)
}

func TestRefreshTokenAPI(t *testing.T) {
	user := randomUser()
	password := util.RandomString(20)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user.Password = hashedPassword

	session := db.Session{
		ID:        uuid.New(),
		Username:  user.Username,
		UserAgent: util.RandomString(10),
		ClientIp:  util.RandomString(10),
		IsBlocked: false,
		CreatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(session, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchRefreshToken(t, recorder.Body, server)
			},
		},
		{
			name: "MissmatchedSessionToken",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).Return(db.Session{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "IncorrectSessionUser",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{
						ID:           session.ID,
						ExpiresAt:    session.ExpiresAt,
						ClientIp:     session.ClientIp,
						RefreshToken: "blabla",
						CreatedAt:    session.CreatedAt,
						Username:     user.Username,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BlockedSession",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{
						ID:           session.ID,
						ExpiresAt:    session.ExpiresAt,
						ClientIp:     session.ClientIp,
						RefreshToken: session.RefreshToken,
						CreatedAt:    session.CreatedAt,
						Username:     session.Username,
						IsBlocked:    true,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "MissingToken",
			buildStubs: func(store *mockdb.MockStore) {
				session.RefreshToken = ""
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "IncorrectToken",
			buildStubs: func(store *mockdb.MockStore) {
				session.RefreshToken = "fake_token"
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionSuddenlyExpired",
			buildStubs: func(store *mockdb.MockStore) {
				session.ExpiresAt = time.Now()
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "SessionDoesNotExist",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, server *Server) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// create a fake token
			token, payload, err := server.tokenMaker.CreateToken(user.Username, 24*time.Hour)
			require.NoError(t, err)
			// set the refresh token to the newly created token
			session.RefreshToken = token
			session.ExpiresAt = payload.ExpiredAt

			tc.buildStubs(store)

			jsonPayload := fmt.Sprintf(`{"refresh_token": "%s"}`, session.RefreshToken)
			jsonData := []byte(jsonPayload)
			request, err := http.NewRequest(http.MethodPost, "/tokens/refresh", bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder, server)
		})
	}
}
