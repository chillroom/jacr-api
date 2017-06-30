-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists history
(
	id serial not null
		constraint history_pkey
			primary key,
	dub_id char(24) not null,
	score_down integer default 0 not null,
	score_grab integer default 0 not null,
	score_up integer default 0 not null,
	song integer not null
		constraint history_songs_id_fk
			references songs,
	"user" integer not null
		constraint history_dubtrack_users_id_fk
			references dubtrack_users,
	time timestamp default now() not null
)
;

create unique index if not exists history_id_uindex
	on history (id)
;

create unique index if not exists history_dub_id_uindex
	on history (dub_id)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists history;
