-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists response_groups
(
	id serial not null
		constraint response_groups_pkey
			primary key,
	messages text[] not null
)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists response_groups;
