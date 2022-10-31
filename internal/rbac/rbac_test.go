package rbac

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	mockdb "github.com/mateoradman/tempus/internal/db/mock"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/token"
	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
)

func setupUser(role util.AccessRole) db.User {
	return db.User{
		ID:        util.RandomInt(1, 1000),
		Name:      util.RandomString(5),
		Username:  util.RandomString(5),
		Surname:   util.RandomString(5),
		Gender:    util.Pointer(util.RandomGender()),
		Email:     util.RandomEmail(),
		CreatedAt: time.Now().UTC(),
		BirthDate: time.Now().UTC(),
		Role:      int32(role),
	}
}

func getContext(t *testing.T, user db.User) *gin.Context {
	ctx := &gin.Context{}
	maker, err := token.NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	token, payload, err := maker.CreateToken(user.Username, time.Hour)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	require.NoError(t, err)
	ctx.Set(util.AuthPayloadKey, payload)
	return ctx
}

func setupRBAC(t *testing.T) (*mockdb.MockStore, *RBACService) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)
	rbacService := NewRBACService(store)
	return store, rbacService
}

func TestGetUserRBAC(t *testing.T) {
	user := setupUser(util.DefaultRole)
	store, rbacService := setupRBAC(t)
	store.EXPECT().
		GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
		Times(1).
		Return(user, nil)

	ctx := getContext(t, user)
	gotUser, err := rbacService.getUser(ctx)
	require.Equal(t, user, gotUser)
	require.NoError(t, err)
}

func TestIsAdminTrue(t *testing.T) {
	user := setupUser(util.SuperUserRole)
	_, rbacService := setupRBAC(t)
	result := rbacService.isAdmin(user)
	require.True(t, result)
}

func TestIsAdminFalse(t *testing.T) {
	user := setupUser(util.CompanyAdminRole)
	_, rbacService := setupRBAC(t)
	result := rbacService.isAdmin(user)
	require.False(t, result)
}

func TestIsCompanyAdminTrue(t *testing.T) {
	user := setupUser(util.CompanyAdminRole)
	user.CompanyID = util.Pointer(int64(2))
	_, rbacService := setupRBAC(t)
	result := rbacService.isCompanyAdmin(user, user.CompanyID)
	require.True(t, result)
}

func TestIsCompanyAdminTrueUserIsAdmin(t *testing.T) {
	user := setupUser(util.AdminRole)
	_, rbacService := setupRBAC(t)
	result := rbacService.isCompanyAdmin(user, util.Pointer(int64(1)))
	require.True(t, result)
}

func TestIsCompanyAdminFalse(t *testing.T) {
	user := setupUser(util.DefaultRole)
	_, rbacService := setupRBAC(t)
	result := rbacService.isCompanyAdmin(user, util.Pointer(int64(1)))
	require.False(t, result)
}

func TestEnforceRole(t *testing.T) {
	store, rbacService := setupRBAC(t)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore) db.User
		checkResponse func(t *testing.T, result bool)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) db.User {
				user := setupUser(util.AdminRole)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				return user
			},
			checkResponse: func(t *testing.T, result bool) {
				require.True(t, result)
			},
		},
		{
			name: "UserHasLowerRole",
			buildStubs: func(store *mockdb.MockStore) db.User {
				user := setupUser(util.DefaultRole)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				return user
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
		{
			name: "ErrorFindingUser",
			buildStubs: func(store *mockdb.MockStore) db.User {
				user := db.User{}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, pgx.ErrNoRows)
				return user
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			user := tc.buildStubs(store)
			ctx := getContext(t, user)
			result := rbacService.EnforceRole(ctx, util.CompanyAdminRole)
			tc.checkResponse(t, result)
		})
	}
}

func TestEnforceUser(t *testing.T) {
	store, rbacService := setupRBAC(t)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore) (db.User, int64)
		checkResponse func(t *testing.T, result bool)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) (db.User, int64) {
				user := setupUser(util.AdminRole)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				return user, user.ID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.True(t, result)
			},
		},
		{
			name: "ErrorFindingUser",
			buildStubs: func(store *mockdb.MockStore) (db.User, int64) {
				user := db.User{}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, pgx.ErrNoRows)
				return user, user.ID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
		{
			name: "UserIDMatch",
			buildStubs: func(store *mockdb.MockStore) (db.User, int64) {
				user := setupUser(util.DefaultRole)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				return user, user.ID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.True(t, result)
			},
		},
		{
			name: "UserIDDoesNotExist",
			buildStubs: func(store *mockdb.MockStore) (db.User, int64) {
				user := setupUser(util.DefaultRole)
				user2 := db.User{}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user2, pgx.ErrNoRows)
				return user, user2.ID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
		{
			name: "UserIDCanModifyUser",
			buildStubs: func(store *mockdb.MockStore) (db.User, int64) {
				user := setupUser(util.DefaultRole)
				searchForUser := setupUser(util.CompanyAdminRole)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(searchForUser, nil)
				return user, searchForUser.ID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			user, userID := tc.buildStubs(store)
			ctx := getContext(t, user)
			result := rbacService.EnforceUser(ctx, userID)
			tc.checkResponse(t, result)
		})
	}
}

func TestEnforceCompany(t *testing.T) {
	store, rbacService := setupRBAC(t)

	testCases := []struct {
		name          string
		buildStubs    func(store *mockdb.MockStore) (db.User, *int64)
		checkResponse func(t *testing.T, result bool)
	}{
		{
			name: "OK",
			buildStubs: func(store *mockdb.MockStore) (db.User, *int64){
				user := setupUser(util.AdminRole)
				user.CompanyID = util.Pointer(int64(20))
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				return user, user.CompanyID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.True(t, result)
			},
		},
		{
			name: "ErrorFindingUser",
			buildStubs: func(store *mockdb.MockStore) (db.User, *int64){
				user := db.User{}
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, pgx.ErrNoRows)
				return user, user.CompanyID
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
		{
			name: "UserHasNoPermission",
			buildStubs: func(store *mockdb.MockStore) (db.User, *int64) {
				user := setupUser(util.DefaultRole)
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				return user, util.Pointer(int64(10))
			},
			checkResponse: func(t *testing.T, result bool) {
				require.False(t, result)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			user, companyID := tc.buildStubs(store)
			ctx := getContext(t, user)
			result := rbacService.EnforceCompany(ctx, companyID)
			tc.checkResponse(t, result)
		})
	}
}
