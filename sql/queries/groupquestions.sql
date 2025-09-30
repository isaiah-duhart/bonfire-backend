-- name: CreateGroupQuestions :many
WITH available_questions AS (
    SELECT q.id AS question_id
    FROM questions q
    WHERE NOT EXISTS (
        SELECT 1
        FROM group_questions gq
        WHERE gq.group_id = $1 AND gq.question_id = q.id
    )
    ORDER BY RANDOM()
    LIMIT $2
),
inserted AS (
    INSERT INTO group_questions (
        id,
        group_id,
        date,
        question_id,
        created_at,
        updated_at,
        created_by
    )
    SELECT
        gen_random_uuid(),
        $1,
        $3,
        aq.question_id,
        NOW(),
        NOW(),
        $4
    FROM available_questions aq
    RETURNING id, group_id, date, question_id, created_by
)
SELECT
    i.id,
    i.group_id,
    i.date,
    q.text,
    i.created_by
FROM inserted i
JOIN questions q ON i.question_id = q.id;

-- name: DeleteGroupQuestions :exec
DELETE FROM group_questions
WHERE group_id = $1 and date = $2;

-- name: GetGroupQuestions :many
SELECT DISTINCT ON (gq.id)
    gq.id,
    gq.group_id,
    gq.date,
    q.text,
    gq.created_by
FROM group_questions gq
JOIN questions q ON gq.question_id = q.id 
WHERE gq.date = $1 AND gq.group_id = $2
  AND (
    gq.created_by = $3
    OR EXISTS (
      SELECT 1
      FROM group_responses gr
      WHERE gr.group_question_id = gq.id
    )
  );

-- name: CountGroupQuestions :one
SELECT COUNT(*)
FROM group_questions
WHERE date = $1 AND group_id = $2 AND created_by = $3;



