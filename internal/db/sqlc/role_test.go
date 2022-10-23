package db

import (
	"context"
	"testing"

	"github.com/mateoradman/tempus/internal/rbac"
	"github.com/stretchr/testify/require"
)

func TestGetRole(t *testing.T) {
	gotRole, err := testQueries.GetRole(context.Background(), int32(rbac.DefaultRole))
	require.NoError(t, err)
	require.NotEmpty(t, gotRole)
	require.NotEmpty(t, gotRole.ID)
	require.Equal(t, gotRole.Role, int32(rbac.DefaultRole))
}

func TestListRoles(t *testing.T) {

	arg := ListRolesParams{
		Limit:  5,
		Offset: 0,
	}

	roles, err := testQueries.ListRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, roles, int(arg.Limit))

	for i, role := range roles {
		require.NotEmpty(t, role)
		roleValue := int32(i + 1)
		require.Equal(t, role.Role, roleValue)
	}
}
