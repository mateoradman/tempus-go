// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: entry.sql

package db

import (
	"context"
	"time"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
user_id, start_time, end_time, date
) VALUES (
$1,
$2,
$3,
$4
)
RETURNING id, user_id, start_time, end_time, created_at, updated_at, date
`

type CreateEntryParams struct {
	UserID    int64      `json:"user_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Date      time.Time  `json:"date"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRow(ctx, createEntry,
		arg.UserID,
		arg.StartTime,
		arg.EndTime,
		arg.Date,
	)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StartTime,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Date,
	)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :one
DELETE
FROM entries
WHERE id = $1
RETURNING id, user_id, start_time, end_time, created_at, updated_at, date
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRow(ctx, deleteEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StartTime,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Date,
	)
	return i, err
}

const getEntry = `-- name: GetEntry :one
SELECT id, user_id, start_time, end_time, created_at, updated_at, date 
FROM entries
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRow(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StartTime,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Date,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, user_id, start_time, end_time, created_at, updated_at, date
FROM entries
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.Query(ctx, listEntries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StartTime,
			&i.EndTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Date,
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

const listUserEntries = `-- name: ListUserEntries :many
SELECT id, user_id, start_time, end_time, created_at, updated_at, date
FROM entries
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListUserEntriesParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUserEntries(ctx context.Context, arg ListUserEntriesParams) ([]Entry, error) {
	rows, err := q.db.Query(ctx, listUserEntries, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StartTime,
			&i.EndTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Date,
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

const updateEntry = `-- name: UpdateEntry :one
UPDATE entries
SET 
user_id = $2, 
start_time = $3, 
end_time = $4,
date = $5
WHERE id = $1
RETURNING id, user_id, start_time, end_time, created_at, updated_at, date
`

type UpdateEntryParams struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Date      time.Time  `json:"date"`
}

func (q *Queries) UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error) {
	row := q.db.QueryRow(ctx, updateEntry,
		arg.ID,
		arg.UserID,
		arg.StartTime,
		arg.EndTime,
		arg.Date,
	)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StartTime,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Date,
	)
	return i, err
}
