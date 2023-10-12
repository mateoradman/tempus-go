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
	"github.com/jackc/pgx/v5"
	mockdb "github.com/mateoradman/tempus/internal/db/mock"
	db "github.com/mateoradman/tempus/internal/db/sqlc"
	"github.com/mateoradman/tempus/internal/token"
	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
)

func requireBodyMatchCompany(t *testing.T, body *bytes.Buffer, company db.Company) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCompany db.Company
	err = json.Unmarshal(data, &gotCompany)
	require.NoError(t, err)
	require.Equal(t, company, gotCompany)
}

func requireBodyMatchCompanyList(t *testing.T, body *bytes.Buffer, companies []db.Company) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCompanies []db.Company
	err = json.Unmarshal(data, &gotCompanies)
	require.NoError(t, err)
	require.Equal(t, companies, gotCompanies)

}

func randomCompany() db.Company {
	return db.Company{
		ID:        util.RandomInt(1, 1000),
		Name:      util.RandomString(100),
		CreatedAt: time.Now().UTC(),
	}
}

func TestCreateCompanyAPI(t *testing.T) {
	user := randomUser()
	company := randomCompany()

	testCases := []struct {
		name          string
		companyName   string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			companyName: company.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Eq(company.Name)).
					Times(1).
					Return(company, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchCompany(t, recorder.Body, company)
			},
		},
		{
			name:        "InternalServerError",
			companyName: company.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Eq(company.Name)).
					Times(1).
					Return(db.Company{}, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "BadRequest",
			companyName: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Unauthorized",
			companyName: company.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateCompany(gomock.Any(), gomock.Any()).
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

			jsonData := []byte(fmt.Sprintf(`{"name": "%s"}`, tc.companyName))
			request, err := http.NewRequest(http.MethodPost, "/companies", bytes.NewBuffer(jsonData))
			tc.setupAuth(t, request, server.tokenMaker)
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			tc.setupAuth(t, request, server.tokenMaker)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetCompanyAPI(t *testing.T) {
	user := randomUser()
	company := randomCompany()

	testCases := []struct {
		name          string
		companyID     int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(company, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCompany(t, recorder.Body, company)
			},
		},
		{
			name:      "NotFound",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(db.Company{}, pgx.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(db.Company{}, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			companyID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/companies/%d", tc.companyID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			tc.setupAuth(t, request, server.tokenMaker)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestDeleteCompanyAPI(t *testing.T) {
	user := randomUser()
	company := randomCompany()

	testCases := []struct {
		name          string
		companyID     int64
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(company, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCompany(t, recorder.Body, company)
			},
		},
		{
			name:      "NotFound",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(db.Company{}, pgx.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(db.Company{}, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			companyID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteCompany(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/companies/%d", tc.companyID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			tc.setupAuth(t, request, server.tokenMaker)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestUpdateCompanyAPI(t *testing.T) {
	user := randomUser()
	company := randomCompany()
	arg := db.UpdateCompanyParams{
		Name: company.Name,
		ID:   company.ID,
	}

	testCases := []struct {
		name          string
		companyID     int64
		companyName   string
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			companyID:   arg.ID,
			companyName: arg.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(company, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCompany(t, recorder.Body, company)
			},
		},
		{
			name:        "NotFound",
			companyID:   arg.ID,
			companyName: arg.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Company{}, pgx.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:        "InternalServerError",
			companyID:   arg.ID,
			companyName: arg.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Company{}, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "InvalidCompanyID",
			companyID:   0,
			companyName: arg.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "InvalidCompanyName",
			companyID:   arg.ID,
			companyName: "",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Unauthorized",
			companyID:   arg.ID,
			companyName: arg.Name,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateCompany(gomock.Any(), gomock.Any()).
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

			jsonData := []byte(fmt.Sprintf(`{"name": "%s"}`, tc.companyName))
			url := fmt.Sprintf("/companies/%d", tc.companyID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
			tc.setupAuth(t, request, server.tokenMaker)
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListCompaniesAPI(t *testing.T) {
	user := randomUser()
	company := randomCompany()
	arg := db.ListCompaniesParams{
		Offset: int32(util.RandomInt(0, 100000)),
		Limit:  int32(util.RandomInt(1, 100)),
	}
	returnVal := []db.Company{company}

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
					ListCompanies(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnVal, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchCompanyList(t, recorder.Body, returnVal)
			},
		},
		{
			name:   "InternalServerError",
			offset: arg.Offset,
			limit:  arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanies(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnVal, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
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
					ListCompanies(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
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
					ListCompanies(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
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
					ListCompanies(gomock.Any(), gomock.Any()).
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

			request, err := http.NewRequest(http.MethodGet, "/companies", nil)
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

func TestListCompanyEmployeesAPI(t *testing.T) {
	company := randomCompany()
	arg := db.ListCompanyEmployeesParams{
		ID:     company.ID,
		Offset: int32(util.RandomInt(0, 100000)),
		Limit:  int32(util.RandomInt(1, 100)),
	}
	user := randomUser()
	returnVal := []db.User{user}

	testCases := []struct {
		name          string
		companyID     int64
		offset        int32
		limit         int32
		buildStubs    func(store *mockdb.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			companyID: company.ID,
			offset:    arg.Offset,
			limit:     arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanyEmployees(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnVal, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUserList(t, recorder.Body, returnVal)
			},
		},
		{
			name:      "InternalServerError",
			companyID: company.ID,
			offset:    arg.Offset,
			limit:     arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanyEmployees(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(returnVal, pgx.ErrTxClosed)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidOffset",
			companyID: company.ID,
			limit:     arg.Limit,
			offset:    -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanyEmployees(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidLimit",
			companyID: company.ID,
			offset:    arg.Offset,
			limit:     1000,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanyEmployees(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InvalidCompanyID",
			companyID: 0,
			offset:    arg.Offset,
			limit:     arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanyEmployees(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, util.AuthTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Unauthorized",
			companyID: arg.ID,
			offset:    arg.Offset,
			limit:     arg.Limit,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListCompanyEmployees(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/companies/%d/employees", tc.companyID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
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
