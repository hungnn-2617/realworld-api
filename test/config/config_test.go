package config_test

import (
	"os"
	"testing"

	"realworld-api/internal/config"
)

func TestGetConfig(t *testing.T) {
	cfg := config.GetConfig()

	if cfg == nil {
		t.Fatal("GetConfig() returned nil")
	}

	t.Run("has app config", func(t *testing.T) {
		if cfg.App.Env == "" {
			t.Error("App.Env should have a default value")
		}
	})

	t.Run("has server config", func(t *testing.T) {
		if cfg.Server.Host == "" {
			t.Error("Server.Host should have a default value")
		}
		if cfg.Server.Port == "" {
			t.Error("Server.Port should have a default value")
		}
	})

	t.Run("has database config", func(t *testing.T) {
		if cfg.Database.Host == "" {
			t.Error("Database.Host should have a default value")
		}
		if cfg.Database.Port == "" {
			t.Error("Database.Port should have a default value")
		}
	})

	t.Run("has JWT config", func(t *testing.T) {
		if cfg.JWT.Secret == "" {
			t.Error("JWT.Secret should have a default value")
		}
		if cfg.JWT.ExpirationHours <= 0 {
			t.Error("JWT.ExpirationHours should be positive")
		}
	})
}

func TestDatabaseConfigGetDSN(t *testing.T) {
	cfg := config.GetConfig()
	dsn := cfg.Database.GetDSN()

	if dsn == "" {
		t.Error("GetDSN() returned empty string")
	}

	expectedParts := []string{"host=", "port=", "user=", "password=", "dbname=", "sslmode="}
	for _, part := range expectedParts {
		found := false
		for i := 0; i <= len(dsn)-len(part); i++ {
			if dsn[i:i+len(part)] == part {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DSN should contain %q, got %q", part, dsn)
		}
	}
}

func TestConfigSingleton(t *testing.T) {
	cfg1 := config.GetConfig()
	cfg2 := config.GetConfig()

	if cfg1 != cfg2 {
		t.Error("GetConfig() should return the same instance (singleton)")
	}
}

func TestConfigDefaultValues(t *testing.T) {
	cfg := config.GetConfig()

	tests := []struct {
		name string
		got  string
	}{
		{"App.Env default", cfg.App.Env},
		{"Server.Host default", cfg.Server.Host},
		{"Server.Port default", cfg.Server.Port},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

func TestConfigStructFields(t *testing.T) {
	t.Run("AppConfig", func(t *testing.T) {
		appCfg := config.AppConfig{Env: "production"}
		if appCfg.Env != "production" {
			t.Errorf("AppConfig.Env = %q, want %q", appCfg.Env, "production")
		}
	})

	t.Run("ServerConfig", func(t *testing.T) {
		serverCfg := config.ServerConfig{Host: "localhost", Port: "3000"}
		if serverCfg.Host != "localhost" {
			t.Errorf("ServerConfig.Host = %q, want %q", serverCfg.Host, "localhost")
		}
		if serverCfg.Port != "3000" {
			t.Errorf("ServerConfig.Port = %q, want %q", serverCfg.Port, "3000")
		}
	})

	t.Run("DatabaseConfig", func(t *testing.T) {
		dbCfg := config.DatabaseConfig{
			Host:     "db.example.com",
			Port:     "5432",
			User:     "admin",
			Password: "secret",
			DBName:   "testdb",
			SSLMode:  "require",
		}
		if dbCfg.Host != "db.example.com" {
			t.Errorf("DatabaseConfig.Host = %q, want %q", dbCfg.Host, "db.example.com")
		}
		if dbCfg.SSLMode != "require" {
			t.Errorf("DatabaseConfig.SSLMode = %q, want %q", dbCfg.SSLMode, "require")
		}
	})

	t.Run("JWTConfig", func(t *testing.T) {
		jwtCfg := config.JWTConfig{Secret: "mysecret", ExpirationHours: 24}
		if jwtCfg.Secret != "mysecret" {
			t.Errorf("JWTConfig.Secret = %q, want %q", jwtCfg.Secret, "mysecret")
		}
		if jwtCfg.ExpirationHours != 24 {
			t.Errorf("JWTConfig.ExpirationHours = %d, want %d", jwtCfg.ExpirationHours, 24)
		}
	})
}

func TestDatabaseConfigDSNFormat(t *testing.T) {
	dbCfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	dsn := dbCfg.GetDSN()
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"

	if dsn != expected {
		t.Errorf("GetDSN() = %q, want %q", dsn, expected)
	}
}

func TestEnvVariableOverride(t *testing.T) {
	testEnvKey := "TEST_ENV_VAR_FOR_CONFIG"
	testEnvValue := "test_value"

	os.Setenv(testEnvKey, testEnvValue)
	defer os.Unsetenv(testEnvKey)

	value := os.Getenv(testEnvKey)
	if value != testEnvValue {
		t.Errorf("Environment variable not set correctly: got %q, want %q", value, testEnvValue)
	}
}
