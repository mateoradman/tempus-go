package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	mockdb "github.com/mateoradman/tempus/internal/db/mock"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/token"
	"github.com/mateoradman/tempus/internal/types"
	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.Password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func requireBodyMatchUserList(t *testing.T, body *bytes.Buffer, users []db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUsers []db.User
	err = json.Unmarshal(data, &gotUsers)
	require.NoError(t, err)
	require.Equal(t, users, gotUsers)
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	fmt.Sprintln(data)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}

func randomUser() db.User {
	return db.User{
		ID:        util.RandomInt(1, 1000),
		Name:      util.RandomString(5),
		Username:  util.RandomString(5),
		Surname:   util.RandomString(5),
		Gender:    util.RandomGender(),
		Email:     util.RandomEmail(),
		Language:  util.RandomLanguage(),
		CreatedAt: time.Now().UTC(),
		BirthDate: time.Now().UTC(),
		Role:      types.EmployeeRole,
	}
}

func TestCreateUserAPI(t *testing.T) {
	user := randomUser()
	password := util.RandomString(10)

	arg := db.CreateUserParams{
		Username:  user.Username,
		Password:  password,
		Email:     user.Email,
		Name:      user.Name,
		Surname:   user.Surname,
		Gender:    user.Gender,
		CompanyID: user.CompanyID,
		BirthDate: user.BirthDate,
		Language:  user.Language,
		Country:   user.Country,
		Timezone:  user.Timezone,
		ManagerID: user.ManagerID,
		TeamID:    user.TeamID,
	}

	testCases := []struct {
		name          string
		arg           db.CreateUserParams
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalServerError",
			arg:  arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			arg:  db.CreateUserParams{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			jsonData, err := json.Marshal(&tc.arg)
			require.NoError(t, err)

			t.Logf("jsonData: %s", jsonData)
			request, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			server.router.ServeHTTP(recorder, request)
			t.Logf("%v", recorder.Body)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	user := randomUser()
	password := util.RandomString(20)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user.Password = hashedPassword

	testCases := []struct {
		name          string
		arg           gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			arg: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "UserNotFound",
			arg: gin.H{
				"username": "FakeUsername",
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			arg: gin.H{
				"username": "############",
				"password": "FakePassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			arg: gin.H{
				"username": user.Username,
				"password": "FakePassword",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalServerError",
			arg: gin.H{
				"username": user.Username,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
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
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			jsonData, err := json.Marshal(&tc.arg)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d", tc.userID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteUserAPI(t *testing.T) {
	user := randomUser()

	testCases := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%d", tc.userID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateUserAPI(t *testing.T) {
	user := randomUser()
	arg := db.UpdateUserParams{
		ID:        user.ID,
		Name:      &user.Name,
		Surname:   &user.Surname,
		Gender:    &user.Gender,
		BirthDate: &user.BirthDate,
		Language:  &user.Language,
		Country:   user.Country,
	}

	testCases := []struct {
		name          string
		userID        int64
		arg           db.UpdateUserParams
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: arg.ID,
			arg:    arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(user, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: arg.ID,
			arg:    arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			userID: arg.ID,
			arg:    arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidUserID",
			userID: 0,
			arg:    arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			userID: arg.ID,
			arg:    arg,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			jsonData, err := json.Marshal(&tc.arg)
			require.NoError(t, err)
			url := fmt.Sprintf("/users/%d", tc.userID)
			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)

			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListUsersAPI(t *testing.T) {
	user := randomUser()
	arg := db.ListUsersParams{
		Offset: int32(util.RandomInt(0, 100000)),
		Limit:  int32(util.RandomInt(1, 100)),
	}
	returnVal := []db.User{user}

	testCases := []struct {
		name          string
		offset        int32
		limit         int32
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			offset: arg.Offset,
			limit:  arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnVal, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserList(t, recorder.Body, returnVal)
			},
		},
		{
			name:   "InternalServerError",
			offset: arg.Offset,
			limit:  arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnVal, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidOffset",
			limit:  arg.Limit,
			offset: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidLimit",
			offset: arg.Offset,
			limit:  1000,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Unauthorized",
			offset: arg.Offset,
			limit:  arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/users", nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			q := request.URL.Query()
			q.Set("offset", fmt.Sprintf("%d", tc.offset))
			q.Set("limit", fmt.Sprintf("%d", tc.limit))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
