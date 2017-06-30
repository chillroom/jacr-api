-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists response_commands
(
	id serial not null
		constraint response_commands_id_pkey
			primary key,
	name varchar(32) not null
		constraint response_commands_name_pk
			unique,
	"group" integer not null
		constraint response_commands_response_groups_id_fk
			references response_groups,
	constraint response_commands_name_group_pk
		unique (name, "group")
)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists response_commands;
