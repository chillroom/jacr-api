-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
CREATE EXTENSION IF EXISTS "uuid-ossp";
