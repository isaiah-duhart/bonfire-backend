
-- name: CreateGroup :one
WITH users as (
    SELECT * FROM users WHERE users.email = $3
)
INSERT INTO groups (id, group_id, group_name, user_id, created_at, updated_at)
SELECT gen_random_uuid(), $1, $2, users.id, NOW(), NOW()
FROM users
RETURNING *;

-- name: CreateGroupWithUserID :one
INSERT INTO groups (id, group_id, group_name, user_id, created_at, updated_at)
VALUES(gen_random_uuid(), $1, $2, $3, NOW(), NOW())
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM users
WHERE id = $1;

-- name: GetGroupsByUserID :many
SELECT * FROM groups
WHERE user_id = $1;

-- name: IsUserInGroup :one
SELECT EXISTS (
    SELECT 1
    FROM groups
    WHERE group_id = $1 AND user_id = $2
);