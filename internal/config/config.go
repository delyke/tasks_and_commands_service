package config

import (
	"github.com/joho/godotenv"

	"github.com/delyke/tasks_and_commands_service/internal/config/env"
)

var appConfig *config

type config struct {
	Logger LoggerConfig
	HTTP   HTTPConfig
	JWT    JWTConfig
	Redis  RedisConfig
	MySQL  MySQLConfig
}

// Load reads environment variables from the specified path and initializes configuration.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}
	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}
	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}
	jwtCfg, err := env.NewJWTConfig()
	if err != nil {
		return err
	}
	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}
	mysqlCfg, err := env.NewMySQLConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		Logger: loggerCfg,
		HTTP:   httpCfg,
		JWT:    jwtCfg,
		Redis:  redisCfg,
		MySQL:  mysqlCfg,
	}
	return nil
}

// AppConfig returns the application configuration singleton.
func AppConfig() *config {
	return appConfig
}
