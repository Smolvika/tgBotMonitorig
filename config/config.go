package config

import "os"

type ConfigDB struct {
	Host     string
	Port     string
	Username string
	DBName   string
	SSLMode  string
	Password string
}

type Config struct {
	ConfigDB      ConfigDB
	TelegramToken string
}

func New() *Config {
	return &Config{
		ConfigDB: ConfigDB{
			Host:     getEnv("HOST", ""),
			Port:     getEnv("PORT", ""),
			Username: getEnv("USER", ""),
			DBName:   getEnv("NAME", ""),
			SSLMode:  getEnv("SSLMODE", ""),
			Password: getEnv("DB_PASSWORD", ""),
		},
		TelegramToken: getEnv("TOKEN", ""),
	}
}
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
