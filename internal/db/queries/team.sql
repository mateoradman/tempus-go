-- name: CreateTeam :one
INSERT INTO teams (
    name,
    manager_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetTeam :one
SELECT * 
FROM teams
WHERE id = $1
LIMIT 1;

-- name: ListTeams :many
SELECT *
FROM teams
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTeam :one
UPDATE teams
SET name = $2, manager_id = $3
WHERE id = $1
RETURNING *;

-- name: ListTeamMembers :many
SELECT *
FROM users
WHERE team_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteTeam :one
DELETE
FROM teams
WHERE id = $1
RETURNING *;