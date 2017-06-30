--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.3
-- Dumped by pg_dump version 9.6.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


SET search_path = public, pg_catalog;

--
-- Name: last_seen; Type: TYPE; Schema: public; Owner: jacr
--

CREATE TYPE last_seen AS ENUM (
    'join',
    'quit',
    'message',
    'update'
);


ALTER TYPE last_seen OWNER TO jacr;

--
-- Name: skip_reason; Type: TYPE; Schema: public; Owner: jacr
--

CREATE TYPE skip_reason AS ENUM (
    'forbidden',
    'nsfw',
    'op',
    'theme',
    'unavailable'
);


ALTER TYPE skip_reason OWNER TO jacr;

--
-- Name: song_type; Type: TYPE; Schema: public; Owner: jacr
--

CREATE TYPE song_type AS ENUM (
    'youtube',
    'soundcloud'
);


ALTER TYPE song_type OWNER TO jacr;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: dubtrack_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE dubtrack_users (
    id integer NOT NULL,
    karma integer DEFAULT 0 NOT NULL,
    dub_id character(24) NOT NULL,
    username text NOT NULL,
    seen_time timestamp without time zone DEFAULT now() NOT NULL,
    seen_message text DEFAULT ''::text NOT NULL,
    seen_type last_seen,
    rethink_id character varying(36)
);


ALTER TABLE dubtrack_users OWNER TO postgres;

--
-- Name: dubtrack_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE dubtrack_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE dubtrack_users_id_seq OWNER TO postgres;

--
-- Name: dubtrack_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE dubtrack_users_id_seq OWNED BY dubtrack_users.id;


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE goose_db_version OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE goose_db_version_id_seq OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE goose_db_version_id_seq OWNED BY goose_db_version.id;


--
-- Name: history; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE history (
    id integer NOT NULL,
    dub_id character(24) NOT NULL,
    score_down integer DEFAULT 0 NOT NULL,
    score_grab integer DEFAULT 0 NOT NULL,
    score_up integer DEFAULT 0 NOT NULL,
    song integer NOT NULL,
    "user" integer NOT NULL,
    "time" timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE history OWNER TO postgres;

--
-- Name: history_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE history_id_seq OWNER TO postgres;

--
-- Name: history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE history_id_seq OWNED BY history.id;


--
-- Name: notices; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE notices (
    id integer NOT NULL,
    message text NOT NULL,
    title text NOT NULL
);


ALTER TABLE notices OWNER TO postgres;

--
-- Name: notices_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE notices_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE notices_id_seq OWNER TO postgres;

--
-- Name: notices_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE notices_id_seq OWNED BY notices.id;


--
-- Name: response_commands; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE response_commands (
    id integer NOT NULL,
    name character varying(32) NOT NULL,
    "group" integer NOT NULL
);


ALTER TABLE response_commands OWNER TO postgres;

--
-- Name: response_commands_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE response_commands_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE response_commands_id_seq OWNER TO postgres;

--
-- Name: response_commands_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE response_commands_id_seq OWNED BY response_commands.id;


--
-- Name: response_groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE response_groups (
    id integer NOT NULL,
    messages text[] NOT NULL
);


ALTER TABLE response_groups OWNER TO postgres;

--
-- Name: response_groups_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE response_groups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE response_groups_id_seq OWNER TO postgres;

--
-- Name: response_groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE response_groups_id_seq OWNED BY response_groups.id;


--
-- Name: settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE settings (
    name text NOT NULL,
    value jsonb NOT NULL
);


ALTER TABLE settings OWNER TO postgres;

--
-- Name: songs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE songs (
    id integer NOT NULL,
    fkid character varying(32) NOT NULL,
    name text NOT NULL,
    last_play timestamp without time zone NOT NULL,
    skip_reason skip_reason,
    recent_plays integer DEFAULT 0 NOT NULL,
    total_plays integer DEFAULT 0 NOT NULL,
    rethink_id character varying(36),
    type song_type NOT NULL,
    retagged boolean DEFAULT false NOT NULL,
    autoretagged boolean DEFAULT false
);


ALTER TABLE songs OWNER TO postgres;

--
-- Name: songs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE songs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE songs_id_seq OWNER TO postgres;

--
-- Name: songs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE songs_id_seq OWNED BY songs.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE users (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    password character(60) NOT NULL,
    email character varying(254) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    activated boolean DEFAULT false NOT NULL,
    banned boolean DEFAULT false NOT NULL,
    slug character varying(255) NOT NULL,
    level integer DEFAULT 1 NOT NULL
);


ALTER TABLE users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: dubtrack_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY dubtrack_users ALTER COLUMN id SET DEFAULT nextval('dubtrack_users_id_seq'::regclass);


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY goose_db_version ALTER COLUMN id SET DEFAULT nextval('goose_db_version_id_seq'::regclass);


--
-- Name: history id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history ALTER COLUMN id SET DEFAULT nextval('history_id_seq'::regclass);


--
-- Name: notices id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY notices ALTER COLUMN id SET DEFAULT nextval('notices_id_seq'::regclass);


--
-- Name: response_commands id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands ALTER COLUMN id SET DEFAULT nextval('response_commands_id_seq'::regclass);


--
-- Name: response_groups id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_groups ALTER COLUMN id SET DEFAULT nextval('response_groups_id_seq'::regclass);


--
-- Name: songs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY songs ALTER COLUMN id SET DEFAULT nextval('songs_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: dubtrack_users dubtrack_users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY dubtrack_users
    ADD CONSTRAINT dubtrack_users_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: history history_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history
    ADD CONSTRAINT history_pkey PRIMARY KEY (id);


--
-- Name: notices notices_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY notices
    ADD CONSTRAINT notices_pkey PRIMARY KEY (id);


--
-- Name: response_commands response_commands_id_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands
    ADD CONSTRAINT response_commands_id_pkey PRIMARY KEY (id);


--
-- Name: response_commands response_commands_name_group_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands
    ADD CONSTRAINT response_commands_name_group_pk UNIQUE (name, "group");


--
-- Name: response_commands response_commands_name_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands
    ADD CONSTRAINT response_commands_name_pk UNIQUE (name);


--
-- Name: response_groups response_groups_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_groups
    ADD CONSTRAINT response_groups_pkey PRIMARY KEY (id);


--
-- Name: settings settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (name);


--
-- Name: songs songs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY songs
    ADD CONSTRAINT songs_pkey PRIMARY KEY (id);


--
-- Name: songs songs_type_fkid_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY songs
    ADD CONSTRAINT songs_type_fkid_pk UNIQUE (type, fkid);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_id_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_id_pkey PRIMARY KEY (id);


--
-- Name: users users_slug_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_slug_key UNIQUE (slug);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: dubtrack_users_dub_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX dubtrack_users_dub_id_uindex ON dubtrack_users USING btree (dub_id);


--
-- Name: dubtrack_users_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX dubtrack_users_id_uindex ON dubtrack_users USING btree (id);


--
-- Name: dubtrack_users_rethinkid_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX dubtrack_users_rethinkid_uindex ON dubtrack_users USING btree (rethink_id);


--
-- Name: dubtrack_users_username_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX dubtrack_users_username_index ON dubtrack_users USING btree (username);


--
-- Name: history_dub_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX history_dub_id_uindex ON history USING btree (dub_id);


--
-- Name: history_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX history_id_uindex ON history USING btree (id);


--
-- Name: notices_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX notices_id_uindex ON notices USING btree (id);


--
-- Name: notices_title_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX notices_title_uindex ON notices USING btree (title);


--
-- Name: songs_fkid_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX songs_fkid_uindex ON songs USING btree (fkid);


--
-- Name: songs_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX songs_id_uindex ON songs USING btree (id);


--
-- Name: songs_rethink_id_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX songs_rethink_id_uindex ON songs USING btree (rethink_id);


--
-- Name: history history_dubtrack_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history
    ADD CONSTRAINT history_dubtrack_users_id_fk FOREIGN KEY ("user") REFERENCES dubtrack_users(id);


--
-- Name: history history_songs_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY history
    ADD CONSTRAINT history_songs_id_fk FOREIGN KEY (song) REFERENCES songs(id);


--
-- Name: response_commands response_commands_response_groups_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY response_commands
    ADD CONSTRAINT response_commands_response_groups_id_fk FOREIGN KEY ("group") REFERENCES response_groups(id);


--
-- PostgreSQL database dump complete
--

