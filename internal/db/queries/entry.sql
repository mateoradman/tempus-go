-- name: CreateEntry :one
INSERT INTO entries (
user_id, start_time, end_time, date
) VALUES (
$1,
$2,
$3,
$4
)
RETURNING *;

-- name: GetEntry :one
SELECT * 
FROM entries
WHERE id = $1
LIMIT 1;

-- name: ListEntries :many
SELECT *
FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListUserEntries :many
SELECT *
FROM entries
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateEntry :one
UPDATE entries
SET 
user_id = $2, 
start_time = $3, 
end_time = $4,
date = $5
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :one
DELETE
FROM entries
WHERE id = $1
RETURNING *;