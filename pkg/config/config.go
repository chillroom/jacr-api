package config

type Config struct {
	SlackURL      string
	SlackToken    string
	SlackChannels string

	JWTSecret string
	Address   string
}
