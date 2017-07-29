-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- +goose StatementBegin
create or replace function add_response_message(groupName text, msg text) RETURNS VOID AS $rerg$
declare grp int;
begin
    if (msg = '') or (groupName = '') THEN
      RETURN;
    END IF;

    grp := (select "group"::int from response_commands where name = groupName limit 1);

    if grp is null THEN
        with resgroup as (
					insert into response_groups(messages) values (array[msg]) returning id
				)
				insert into response_commands (name, "group")
				values (groupName, (select id from resgroup));

        RETURN;
    END IF;

    update response_groups set messages = array_append(messages, msg) where id = grp;
end;
$rerg$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop function if exists add_response_message(text, text);