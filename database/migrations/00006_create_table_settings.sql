-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists settings
(
	name text not null
		constraint settings_pkey
			primary key,
	value jsonb not null
)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists settings;
