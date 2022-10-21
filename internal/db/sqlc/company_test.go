package db

import (
	"context"
	"testing"
	"time"

	"github.com/mateoradman/tempus/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomCompany(t *testing.T) Company {
	name := util.RandomString(200)

	company, err := testQueries.CreateCompany(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, company)
	require.Equal(t, name, company.Name)

	require.NotZero(t, company.ID)
	require.NotNil(t, company.CreatedAt)
	require.Nil(t, company.UpdatedAt)

	return company
}

func TestCreateCompany(t *testing.T) {
	createRandomCompany(t)
}

func TestGetCompany(t *testing.T) {
	company := createRandomCompany(t)
	gotCompany, err := testQueries.GetCompany(context.Background(), company.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotCompany)
	require.Equal(t, company.ID, gotCompany.ID)
	require.Equal(t, company.Name, gotCompany.Name)
	require.Equal(t, company.CreatedAt, gotCompany.CreatedAt)
	require.Equal(t, company.UpdatedAt, gotCompany.UpdatedAt)
	require.Equal(t, company.Name, gotCompany.Name)
}

func TestUpdateCompany(t *testing.T) {
	company := createRandomCompany(t)
	expectedLen := 25
	name := util.RandomString(expectedLen)
	arg := UpdateCompanyParams{
		ID:   company.ID,
		Name: name,
	}

	updatedCompany, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedCompany)
	require.Equal(t, company.ID, updatedCompany.ID)
	require.Equal(t, name, updatedCompany.Name)
	require.NotNil(t, updatedCompany.UpdatedAt)
	require.WithinDuration(t, time.Now(), *updatedCompany.UpdatedAt, time.Second)
	require.Equal(t, company.CreatedAt, updatedCompany.CreatedAt)
}

func TestDeleteCompany(t *testing.T) {
	company := createRandomCompany(t)
	deletedCompany, err := testQueries.DeleteCompany(context.Background(), company.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedCompany)
	require.Equal(t, company.ID, deletedCompany.ID)
	require.Equal(t, company.Name, deletedCompany.Name)
	require.Nil(t, deletedCompany.UpdatedAt)
	require.Equal(t, company.CreatedAt, deletedCompany.CreatedAt)
}

func TestListCompanies(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 companies
		createRandomCompany(t)
	}

	arg := ListCompaniesParams{
		Limit:  5,
		Offset: 0,
	}

	companies, err := testQueries.ListCompanies(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, companies, int(arg.Limit))

	for _, company := range companies {
		require.NotEmpty(t, company)
	}
}

func TestListEmployee(t *testing.T) {
	users := []User{}
	company := createRandomCompany(t)
	for i := 0; i < 10; i++ {
		user := createRandomUser(t, &company.ID, nil)
		arg := UpdateUserParams{
			ID:        user.ID,
			Name:      &user.Name,
			Surname:   &user.Surname,
			Gender:    &user.Gender,
			BirthDate: &user.BirthDate,
			Language:  &user.Language,
			Country:   &user.Country,
		}
		updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
		require.NoError(t, err)
		users = append(users, updatedUser)
	}
	arg := ListCompanyEmployeesParams{
		ID:     company.ID,
		Limit:  100,
		Offset: 0,
	}
	employees, err := testQueries.ListCompanyEmployees(context.Background(), arg)
	require.NoError(t, err)
	require.Subset(t, employees, users)
}
