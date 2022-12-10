package configs

import "os"

// getenv returns environment variable by name or default value
func getenv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

var (
	Environment = struct {
		ServerPort        string
		Environment       string
		LogOutput         string
		LogLevel          string
		ServerReadTimeout string
	}{
		ServerPort:        getenv("CLOUD_API_SERVER_PORT", "5000"),
		Environment:       getenv("ENVIRONMENT", "development"),
		LogOutput:         getenv("LOG_OUTPUT", "stdout"),
		LogLevel:          getenv("LOG_LEVEL", "info"),
		ServerReadTimeout: getenv("SERVER_READ_TIMEOUT", "60"),
	}
)
