-- name: GetRole :one
SELECT *
FROM roles
WHERE id = $1
LIMIT 1;

-- name: CreateRole :one
INSERT INTO roles (
name,
description
) VALUES (
    $1, $2
) 
RETURNING *;

-- name: UpdateRole :one
UPDATE roles 
SET name = $2, description = $3 
WHERE id = $1 
RETURNING *;

-- name: DeleteRole :one
DELETE
FROM roles
WHERE id = $1 RETURNING *;

-- name: ListRoles :many
SELECT *
FROM roles
ORDER BY id
LIMIT $1
OFFSET $2;
