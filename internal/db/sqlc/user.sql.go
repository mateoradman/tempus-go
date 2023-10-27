// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: user.sql

package db

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one

INSERT INTO users ( username, password, email, name, surname, birth_date, gender, language, country, timezone, company_id, manager_id, team_id)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13) RETURNING id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
`

type CreateUserParams struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
	Language  string    `json:"language"`
	Country   *string   `json:"country"`
	Timezone  string    `json:"timezone"`
	CompanyID *int64    `json:"company_id"`
	ManagerID *int64    `json:"manager_id"`
	TeamID    *int64    `json:"team_id"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.Name,
		arg.Surname,
		arg.BirthDate,
		arg.Gender,
		arg.Language,
		arg.Country,
		arg.Timezone,
		arg.CompanyID,
		arg.ManagerID,
		arg.TeamID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Name,
		&i.Surname,
		&i.CompanyID,
		&i.Password,
		&i.Gender,
		&i.BirthDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.Country,
		&i.Timezone,
		&i.ManagerID,
		&i.TeamID,
		&i.Role,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one

DELETE
FROM users
WHERE id = $1 RETURNING id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, deleteUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Name,
		&i.Surname,
		&i.CompanyID,
		&i.Password,
		&i.Gender,
		&i.BirthDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.Country,
		&i.Timezone,
		&i.ManagerID,
		&i.TeamID,
		&i.Role,
	)
	return i, err
}

const getUser = `-- name: GetUser :one

SELECT id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Name,
		&i.Surname,
		&i.CompanyID,
		&i.Password,
		&i.Gender,
		&i.BirthDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.Country,
		&i.Timezone,
		&i.ManagerID,
		&i.TeamID,
		&i.Role,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one

SELECT id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Name,
		&i.Surname,
		&i.CompanyID,
		&i.Password,
		&i.Gender,
		&i.BirthDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.Country,
		&i.Timezone,
		&i.ManagerID,
		&i.TeamID,
		&i.Role,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one

SELECT id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
FROM users
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Name,
		&i.Surname,
		&i.CompanyID,
		&i.Password,
		&i.Gender,
		&i.BirthDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.Country,
		&i.Timezone,
		&i.ManagerID,
		&i.TeamID,
		&i.Role,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many

SELECT id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Name,
			&i.Surname,
			&i.CompanyID,
			&i.Password,
			&i.Gender,
			&i.BirthDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Language,
			&i.Country,
			&i.Timezone,
			&i.ManagerID,
			&i.TeamID,
			&i.Role,
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

const updateUser = `-- name: UpdateUser :one

UPDATE users
SET name = COALESCE($1, name),
    surname = COALESCE($2, surname),
    gender = COALESCE($3, gender),
    birth_date = COALESCE($4::timestamp, birth_date),
    language = COALESCE($5, language),
    country = COALESCE($6, country)
WHERE id = $7 RETURNING id, username, email, name, surname, company_id, password, gender, birth_date, created_at, updated_at, language, country, timezone, manager_id, team_id, role
`

type UpdateUserParams struct {
	Name      *string    `json:"name"`
	Surname   *string    `json:"surname"`
	Gender    *string    `json:"gender"`
	BirthDate *time.Time `json:"birth_date"`
	Language  *string    `json:"language"`
	Country   *string    `json:"country"`
	ID        int64      `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Name,
		arg.Surname,
		arg.Gender,
		arg.BirthDate,
		arg.Language,
		arg.Country,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Name,
		&i.Surname,
		&i.CompanyID,
		&i.Password,
		&i.Gender,
		&i.BirthDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Language,
		&i.Country,
		&i.Timezone,
		&i.ManagerID,
		&i.TeamID,
		&i.Role,
	)
	return i, err
}
