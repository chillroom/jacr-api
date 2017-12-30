# jacr-api

## Local database configuration

The development server uses the `postgres` user without a password, and the `jacr_dev` database.

## Importing migrations

We use [goose](https://github.com/pressly/goose). [Install goose](https://github.com/pressly/goose#install) and
then run the following command to get your local database up to speed.

```
goose -dir ./database/migrations postgres "user=postgres dbname=jacr_dev sslmode=disable" up
```


## Making changes to the database

Run this command, switching out `add_response_message_fn` with another name to describe your change.
Do not include an extension.
The `sql` param at the end of this command will automatically add it on for you.

```
goose -dir ./database/migrations postgres "user=postgres dbname=jacr_dev sslmode=disable" create add_response_message_fn sql
```
