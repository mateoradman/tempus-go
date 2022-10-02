package db

import (
	"context"
	"testing"
	"time"

	"github.com/mateoradman/tempus/util"
	"github.com/stretchr/testify/require"
)

func createRandomPermission(t *testing.T) Permission {
	name := util.RandomString(200)

	permission, err := testQueries.CreatePermission(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, permission)
	require.Equal(t, name, permission.Name)

	require.NotZero(t, permission.ID)
	require.NotNil(t, permission.CreatedAt)
	require.False(t, permission.UpdatedAt.Valid)

	return permission
}

func TestCreatePermission(t *testing.T) {
	createRandomPermission(t)
}

func TestGetPermission(t *testing.T) {
	permission := createRandomPermission(t)
	gotPermission, err := testQueries.GetPermission(context.Background(), permission.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotPermission)
	require.Equal(t, permission.ID, gotPermission.ID)
	require.Equal(t, permission.Name, gotPermission.Name)
	require.Equal(t, permission.CreatedAt, gotPermission.CreatedAt)
	require.Equal(t, permission.UpdatedAt, gotPermission.UpdatedAt)
	require.Equal(t, permission.Name, gotPermission.Name)
}

func TestUpdatePermission(t *testing.T) {
	permission := createRandomPermission(t)
	expectedLen := 25
	name := util.RandomString(expectedLen)
	arg := UpdatePermissionParams{
		ID:   permission.ID,
		Name: name,
	}

	updatedPermission, err := testQueries.UpdatePermission(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedPermission)
	require.Equal(t, permission.ID, updatedPermission.ID)
	require.Equal(t, name, updatedPermission.Name)
	require.True(t, updatedPermission.UpdatedAt.Valid)
	require.WithinDuration(t, time.Now(), updatedPermission.UpdatedAt.Time, time.Second)
	require.Equal(t, permission.CreatedAt, updatedPermission.CreatedAt)
}

func TestDeletePermission(t *testing.T) {
	permission := createRandomPermission(t)
	deletedPermission, err := testQueries.DeletePermission(context.Background(), permission.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedPermission)
	require.Equal(t, permission.ID, deletedPermission.ID)
	require.Equal(t, permission.Name, deletedPermission.Name)
	require.False(t, deletedPermission.UpdatedAt.Valid)
	require.Equal(t, permission.CreatedAt, deletedPermission.CreatedAt)
}

func TestListPermissions(t *testing.T) {
	for i := 0; i < 10; i++ {
		// Create 10 permissions
		createRandomPermission(t)
	}

	arg := ListPermissionsParams{
		Limit:  5,
		Offset: 0,
	}

	permissions, err := testQueries.ListPermissions(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, permissions, int(arg.Limit))

	for _, permission := range permissions {
		require.NotEmpty(t, permission)
	}
}
