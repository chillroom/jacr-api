-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists users
(
	id serial not null
		constraint users_id_pkey
			primary key,
	username varchar(255) not null
		constraint users_username_key
			unique,
	password char(60) not null,
	email varchar(254) not null
		constraint users_email_key
			unique,
	created_at timestamp default now() not null,
	updated_at timestamp default now() not null,
	activated boolean default false not null,
	banned boolean default false not null,
	slug varchar(255) not null
		constraint users_slug_key
			unique,
	level integer default 1 not null
)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists users;
