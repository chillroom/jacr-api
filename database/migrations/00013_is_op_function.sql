-- +goose Up
-- SQL in this section is executed when the migration is applied.
create or replace function is_op(last_play timestamp without time zone, recent_plays integer) RETURNS boolean
    AS $$ begin return ((now() - last_play) < interval '2 months') and (recent_plays > 10); end; $$
    LANGUAGE plpgsql
    IMMUTABLE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP FUNCTION if exists is_op(timestamp without time zone, INTEGER);