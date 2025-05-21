package config

type Config struct {
	DBConnStr string
}

func LoadConfig() *Config {
	return &Config{
		DBConnStr: "postgres://postgres:postgres@localhost:5432/myproject?sslmode=disable",
	}
}
