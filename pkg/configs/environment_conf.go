package configs

var EnvironmentConf = map[string]string{
	"KOTAL_API_SERVER_PORT": ":5000",
	"ENVIRONMENT":           "development",
	"SERVER_READ_TIMEOUT":   "60",
	"LOG_OUTPUT":            "stdout",
	"LOG_LEVEL":             "info",
}
