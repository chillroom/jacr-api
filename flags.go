package main

import "github.com/namsral/flag"

func getFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	// Basic options
	fs.String("config", "", "path to the config file")
	// fs.String("log_level", "info", "lowest level of log messages to print")

	// RethinkDB connection
	fs.String("rethinkdb_address", "localhost:28015", "Address to the RethinkDB server")
	fs.String("rethinkdb_database", "", "Name of the database to use")
	fs.String("rethinkdb_username", "", "Username of the database user to use")
	fs.String("rethinkdb_password", "", "Password of the database user to use")

	// Slack
	fs.String("slack_url", "", "Slack Room URL")
	fs.String("slack_token", "", "Slack Token")
	fs.String("slack_channels", "", "Slack channels for the user to join")

	// Address
	fs.String("http_address", "127.0.0.1:80", "HTTP Address to bind to")

	return fs
}
