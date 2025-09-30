
-- name: CreateGroup :one
INSERT INTO groups (id, group_id, group_name, user_id, created_at, updated_at)
VALUES(gen_random_uuid(), $1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM users
WHERE id = $1;

-- name: GetGroupsByUserID :many
SELECT * FROM groups
WHERE user_id = $1;