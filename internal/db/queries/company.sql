-- name: CreateCompany :one
INSERT INTO companies (
    name
) VALUES (
    $1
)
RETURNING *;

-- name: GetCompany :one
SELECT * 
FROM companies
WHERE id = $1
LIMIT 1;

-- name: ListCompanies :many
SELECT *
FROM companies
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCompany :one
UPDATE companies
SET name = $2
WHERE id = $1
RETURNING *;

-- name: ListCompanyEmployees :many
SELECT *
FROM users
WHERE users.company_id =
    (SELECT companies.id
    FROM companies
    WHERE companies.id = $1)
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteCompany :one
DELETE
FROM companies
WHERE id = $1
RETURNING *;