-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists notices
(
	id serial not null
		constraint notices_pkey
			primary key,
	message text not null,
	title text not null
)
;

create unique index if not exists notices_id_uindex
	on notices (id)
;

create unique index if not exists notices_title_uindex
	on notices (title)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists notices;
