package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env string
}

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

var (
	config *Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		config = &Config{
			App: AppConfig{
				Env: getEnv("APP_ENV", "development"),
			},
			Server: ServerConfig{
				Host: getEnv("SERVER_HOST", "0.0.0.0"),
				Port: getEnv("SERVER_PORT", "8080"),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "postgres"),
				DBName:   getEnv("DB_NAME", "realworld"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},
			JWT: JWTConfig{
				Secret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
				ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 72),
			},
		}
	})
	return config
}

func GetConfig() *Config {
	if config == nil {
		return Load()
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (d *DatabaseConfig) GetDSN() string {
	return "host=" + d.Host +
		" port=" + d.Port +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.DBName +
		" sslmode=" + d.SSLMode
}
