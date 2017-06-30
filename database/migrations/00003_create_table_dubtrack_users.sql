-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists dubtrack_users
(
	id serial not null
		constraint dubtrack_users_pkey
			primary key,
	karma integer default 0 not null,
	dub_id char(24) not null,
	username text not null,
	seen_time timestamp default now() not null,
	seen_message text default ''::text not null,
	seen_type last_seen,
	rethink_id varchar(36)
)
;

create unique index if not exists dubtrack_users_id_uindex
	on dubtrack_users (id)
;

create unique index if not exists dubtrack_users_dub_id_uindex
	on dubtrack_users (dub_id)
;

create unique index if not exists dubtrack_users_rethinkid_uindex
	on dubtrack_users (rethink_id)
;

create index if not exists dubtrack_users_username_index
	on dubtrack_users (username)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists dubtrack_users;
