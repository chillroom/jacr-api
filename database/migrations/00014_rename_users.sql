-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE public.users RENAME TO accounts;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE public.accounts RENAME TO users;