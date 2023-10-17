-- name: CreateUser :one

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
        $13) RETURNING *;

-- name: GetUser :one

SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByUsername :one

SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByEmail :one

SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: ListUsers :many

SELECT *
FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one

UPDATE users
SET name = COALESCE(sqlc.narg(name), name),
    surname = COALESCE(sqlc.narg(surname), surname),
    gender = COALESCE(sqlc.narg(gender), gender),
    birth_date = COALESCE(sqlc.narg(birth_date)::timestamp, birth_date),
    language = COALESCE(sqlc.narg(language), language),
    country = COALESCE(sqlc.narg(country), country)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: DeleteUser :one

DELETE
FROM users
WHERE id = $1 RETURNING *;
