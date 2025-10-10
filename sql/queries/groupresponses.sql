-- name: CreateGroupResponse :one
WITH inserted as (
  INSERT INTO group_responses (id, group_question_id, response, created_at, author_id)
  VALUES (gen_random_uuid(), $1, $2, NOW(), $3)
  RETURNING *
)
SELECT
  i.id,
  i.group_question_id,
  i.response,
  i.created_at,
  u.name,
  u.id
FROM inserted i
JOIN users u ON i.author_id = u.id;


-- name: GetGroupResponses :many
SELECT 
  gr.id,
  gr.group_question_id,
  gr.response,
  gr.created_at,
  u.name,
  u.id
FROM group_responses gr
JOIN group_questions gq ON gr.group_question_id = gq.id
JOIN users u ON gr.author_id = u.id
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

-- name: IsUserInGroupByGroupQuestionID :one
SELECT EXISTS (
  SELECT 1 
  FROM group_questions as gq
  JOIN groups as g ON gq.group_id = g.group_id
  WHERE gq.id = $1 AND g.user_id = $2
);