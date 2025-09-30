-- +goose Up
ALTER TABLE groups ADD COLUMN name TEXT NOT NULL;

-- +goose Down
ALTER TABLE groups DROP COLUMN name;
