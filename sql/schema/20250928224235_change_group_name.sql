-- +goose Up
ALTER TABLE groups RENAME COLUMN name TO group_name;

-- +goose Down
ALTER TABLE groups RENAME COLUMN group_name TO name;
