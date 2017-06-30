-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists sessions
(
	id uuid default uuid_generate_v4() not null
		constraint sessions_pkey
			primary key,
	dub_id integer not null
		constraint sessions_dubtrack_users_dub_id_fk
			references dubtrack_users
				on delete cascade,
	expired_at timestamp default (now() + '1 day'::interval) not null,
	created_at timestamp default now() not null,
	login_id uuid default uuid_generate_v4() not null,
	activated boolean default false not null
)
;

create unique index if not exists sessions_id_uindex
	on sessions (id)
;

create unique index if not exists sessions_login_id_uindex
	on sessions (login_id)
;



-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists sessions;