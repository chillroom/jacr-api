-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE public.response_commands DROP CONSTRAINT response_commands_response_groups_id_fk;

ALTER TABLE public.response_commands
			ADD CONSTRAINT response_commands_response_groups_id_fk
			FOREIGN KEY ("group") REFERENCES response_groups (id) ON DELETE CASCADE;

-- +goose StatementBegin
create or replace function remove_if_empty_response_group() RETURNS trigger AS $rerg$
    begin
        -- Check that empname and salary are given
        if array_length(new.messages, 1) is null then
            delete from response_groups where id = old.id;
        end if;

				return new;
    end;
$rerg$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER empty_group_trigger AFTER UPDATE or insert ON response_groups
    FOR EACH ROW EXECUTE PROCEDURE remove_if_empty_response_group();

-- +goose StatementBegin
create or replace function remove_orphan_response_group() RETURNS trigger AS $rerg$
    DECLARE
        count integer := 0;
    begin
        select count(*) into count from response_commands where "group" = old.group;
        if count = 0 THEN
          delete from response_groups where id = old.group;
        END IF;

				return old;
    end;
$rerg$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER orphaned_group AFTER delete ON response_commands
    FOR EACH ROW EXECUTE PROCEDURE remove_orphan_response_group();

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TRIGGER IF EXISTS orphaned_group on response_commands;

DROP FUNCTION if exists remove_orphan_response_group();

DROP TRIGGER IF EXISTS empty_group_trigger on response_groups;

DROP FUNCTION if exists remove_if_empty_response_group();

ALTER TABLE public.response_commands DROP CONSTRAINT response_commands_response_groups_id_fk;

ALTER TABLE public.response_commands
			ADD CONSTRAINT response_commands_response_groups_id_fk
			FOREIGN KEY ("group") REFERENCES response_groups (id);
