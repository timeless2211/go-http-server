-- +goose Up
ALTER TABLE users ADD COLUMN hashed_password TEXT NOT NULL;

-- +goose Down
DROP COLUMN hashed_password;
