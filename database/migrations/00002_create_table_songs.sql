-- +goose Up
-- SQL in this section is executed when the migration is applied.
create table if not exists songs
(
	id serial not null
		constraint songs_pkey
			primary key,
	fkid varchar(32) not null,
	name text not null,
	last_play timestamp not null,
	skip_reason skip_reason,
	recent_plays integer default 0 not null,
	total_plays integer default 0 not null,
	rethink_id varchar(36),
	type song_type not null,
	retagged boolean default false not null,
	autoretagged boolean default false,
	constraint songs_type_fkid_pk
		unique (type, fkid)
)
;

create unique index if not exists songs_id_uindex
	on songs (id)
;

create unique index if not exists songs_fkid_uindex
	on songs (fkid)
;

create unique index if not exists songs_rethink_id_uindex
	on songs (rethink_id)
;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table if exists songs;
