-- name: CreateUser :one
INSERT INTO users (
username, email, name, surname, company_id, password, gender, birth_date, language, country, timezone, manager_id, team_id
) VALUES (
$1,
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
$13
)
RETURNING *;

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
SET 
username = $2,
email = $3,
name = $4,
surname = $5,
company_id = $6,
gender = $7,
birth_date = $8,
language = $9,
country = $10,
timezone = $11,
manager_id = $12,
team_id = $13
WHERE id = $1
RETURNING *;

-- name: DeleteUser :one
DELETE
FROM users
WHERE id = $1
RETURNING *;