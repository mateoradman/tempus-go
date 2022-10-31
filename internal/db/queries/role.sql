-- name: GetRole :one
SELECT *
FROM roles
WHERE role = $1
LIMIT 1;

-- name: ListRoles :many
SELECT *
FROM roles
ORDER BY id
LIMIT $1
OFFSET $2;
