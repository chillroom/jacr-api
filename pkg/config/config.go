package config

type Config struct {
	LogLevel string `default:"debug"`

	Database *Database `default:"postgres"`
	Postgres PostgresConfig

	SlackURL      string `required:"true"`
	SlackToken    string `required:"true"`
	SlackChannels string `required:"true"`

	JWTSecret string `required:"true"`
	Address   string `default:"0.0.0.0:8080"`
}

// Database is a database type enum
type Database string

// String implements flag.Value
func (d *Database) String() string {
	return string(*d)
}

// Set implements flag.Value
func (d *Database) Set(value string) error {
	*d = Database(value)
	return nil
}

// Available databases
const (
	Postgres Database = "postgres"
)

// PostgresConfig contains all configuration data for a PostgreSQL connection
type PostgresConfig struct {
	ConnectionString string `required:"true"`
}
