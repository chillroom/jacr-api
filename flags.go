package main

import "github.com/namsral/flag"

func getFlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	// Basic options
	fs.String("config", "", "path to the config file")
	fs.String("jwt_secret", "", "JWT Secret")
	// fs.String("log_level", "info", "lowest level of log messages to print")

	// postgres connection
	fs.String("postgres_addr", "localhost:5432", "Address to the postgres server")
	fs.String("postgres_database", "", "Name of the database to use")
	fs.String("postgres_user", "", "Username of the database user to use")
	fs.String("postgres_password", "", "Password of the database user to use")

	// Slack
	fs.String("slack_url", "", "Slack Room URL")
	fs.String("slack_token", "", "Slack Token")
	fs.String("slack_channels", "", "Slack channels for the user to join")

	// Address
	fs.String("http_address", "127.0.0.1:80", "HTTP Address to bind to")

	return fs
}
