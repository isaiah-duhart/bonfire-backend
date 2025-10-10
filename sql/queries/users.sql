-- name: CreateUser :one
INSERT INTO users (id, email, password, name, birthday, created_at, updated_at)
VALUES(gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;