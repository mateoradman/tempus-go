package db

import (
	"context"
	"testing"
	"time"

	"github.com/mateoradman/tempus/util"
	"github.com/stretchr/testify/require"
)

func createRandomRole(t *testing.T) Role {
	arg := CreateRoleParams{
		Name:        util.RandomString(150),
		Description: util.RandomString(150),
	}

	role, err := testQueries.CreateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, role)
	require.Equal(t, arg.Name, role.Name)
	require.Equal(t, arg.Description, role.Description)
	require.NotZero(t, role.ID)
	require.NotNil(t, role.CreatedAt)
	require.False(t, role.UpdatedAt.Valid)

	return role
}

func TestCreateRole(t *testing.T) {
	createRandomRole(t)
}

func TestGetRole(t *testing.T) {
	role := createRandomRole(t)
	gotRole, err := testQueries.GetRole(context.Background(), role.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotRole)
	require.Equal(t, role.ID, gotRole.ID)
	require.Equal(t, role.Name, gotRole.Name)
	require.Equal(t, role.Description, gotRole.Description)
	require.Equal(t, role.CreatedAt, gotRole.CreatedAt)
	require.Equal(t, role.UpdatedAt, gotRole.UpdatedAt)
	require.Equal(t, role.Name, gotRole.Name)
}

func TestUpdateRole(t *testing.T) {
	role := createRandomRole(t)
	expectedLen := 25
	name := util.RandomString(expectedLen)
	arg := UpdateRoleParams{
		ID:          role.ID,
		Name:        name,
		Description: name,
	}

	updatedRole, err := testQueries.UpdateRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedRole)
	require.Equal(t, role.ID, updatedRole.ID)
	require.Equal(t, name, updatedRole.Name)
	require.Equal(t, updatedRole.Name, updatedRole.Description)
	require.True(t, updatedRole.UpdatedAt.Valid)
	require.WithinDuration(t, time.Now(), updatedRole.UpdatedAt.Time, time.Second)
	require.Equal(t, role.CreatedAt, updatedRole.CreatedAt)
}

func TestDeleteRole(t *testing.T) {
	role := createRandomRole(t)
	deletedRole, err := testQueries.DeleteRole(context.Background(), role.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedRole)
	require.Equal(t, role.ID, deletedRole.ID)
	require.Equal(t, role.Name, deletedRole.Name)
	require.Equal(t, role.Description, deletedRole.Description)
	require.False(t, deletedRole.UpdatedAt.Valid)
	require.Equal(t, role.CreatedAt, deletedRole.CreatedAt)
}

func TestListRoles(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 roles
		createRandomRole(t)
	}

	arg := ListRolesParams{
		Limit:  5,
		Offset: 0,
	}

	roles, err := testQueries.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, int(arg.Limit))

	for _, role := range roles {
		require.NotEmpty(t, role)
	}
}
