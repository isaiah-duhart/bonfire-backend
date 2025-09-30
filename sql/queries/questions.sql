
-- name: CreateQuestion :one
INSERT INTO questions (id, text, created_at, updated_at)
VALUES (gen_random_uuid(), $1, NOW(), NOW())
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = $1;