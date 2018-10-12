-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD COLUMN deleted BOOL NOT NULL DEFAULT false;
ALTER TABLE poets ADD COLUMN deleted BOOL NOT NULL DEFAULT false;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN deleted;
ALTER TABLE poets DROP COLUMN deleted;
