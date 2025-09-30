-- name: CreateGroupResponse :one
INSERT INTO group_responses (id, group_question_id, response, created_at, author_id)
VALUES (gen_random_uuid(), $1, $2, NOW(), $3)
RETURNING *;

-- name: GetGroupResponses :many
SELECT gr.*
FROM group_responses gr
JOIN group_questions gq ON gr.group_question_id = gq.id
WHERE gr.group_question_id = $1
  AND (
    gr.author_id = $2
    OR (
      SELECT COUNT(DISTINCT author_id)
      FROM group_responses
      WHERE group_question_id = $1
    ) >= (
      SELECT COUNT(*)
      FROM groups
      WHERE group_id = gq.group_id
    )
  );