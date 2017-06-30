-- +goose Up
-- SQL in this section is executed when the migration is applied.
DROP INDEX public.dubtrack_users_rethinkid_uindex RESTRICT;
ALTER TABLE public.dubtrack_users DROP rethink_id;

DROP INDEX public.songs_rethink_id_uindex RESTRICT;
ALTER TABLE public.songs DROP rethink_id;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
