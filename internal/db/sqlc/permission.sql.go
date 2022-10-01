// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: permission.sql

package db

import (
	"context"
)

const createPermission = `-- name: CreatePermission :one
INSERT INTO permissions (
name
) VALUES (
    $1
)
RETURNING id, name, created_at, updated_at
`

func (q *Queries) CreatePermission(ctx context.Context, name string) (Permission, error) {
	row := q.db.QueryRow(ctx, createPermission, name)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePermission = `-- name: DeletePermission :one
DELETE
FROM permissions
WHERE id = $1 RETURNING id, name, created_at, updated_at
`

func (q *Queries) DeletePermission(ctx context.Context, id int64) (Permission, error) {
	row := q.db.QueryRow(ctx, deletePermission, id)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPermission = `-- name: GetPermission :one
SELECT id, name, created_at, updated_at
FROM permissions
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetPermission(ctx context.Context, id int64) (Permission, error) {
	row := q.db.QueryRow(ctx, getPermission, id)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPermissions = `-- name: ListPermissions :many
SELECT id, name, created_at, updated_at
FROM permissions
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListPermissionsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]Permission, error) {
	rows, err := q.db.Query(ctx, listPermissions, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Permission{}
	for rows.Next() {
		var i Permission
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePermission = `-- name: UpdatePermission :one
UPDATE permissions 
SET name = $2
WHERE id = $1 
RETURNING id, name, created_at, updated_at
`

type UpdatePermissionParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdatePermission(ctx context.Context, arg UpdatePermissionParams) (Permission, error) {
	row := q.db.QueryRow(ctx, updatePermission, arg.ID, arg.Name)
	var i Permission
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}