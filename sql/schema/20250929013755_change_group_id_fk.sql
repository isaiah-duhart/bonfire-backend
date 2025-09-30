-- +goose Up
ALTER TABLE group_questions
DROP CONSTRAINT group_questions_group_id_fkey;


-- +goose Down
ALTER TABLE group_questions
ADD CONSTRAINT group_questions_group_id_fkey
FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE;
