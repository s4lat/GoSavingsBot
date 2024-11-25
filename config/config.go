package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/s4lat/gosavingsbot/internal/log"
	"time"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" env-default:"5000"`

	BotToken string `env:"BOT_TOKEN" env-required:"true"`

	PgUrl         string        `env:"PG_URL" env-required:"true"`
	PgMaxConns    int32         `env:"PG_MAX_CONNECTIONS" env-default:"32"`
	PgConnTimeout time.Duration `env:"PG_CONNECTION_TIMEOUT" env-default:"30s"`

	//RedisPassword string `env:"REDIS_PASSWORD" env-required:"true"`
	//RedisHost     string `env:"REDIS_HOST" env-required:"true"`
	//RedisPort     int    `env:"REDIS_PORT" env-required:"true"`

	ServeSwagger bool   `env:"SERVE_SWAGGER" env-default:"false"`
	LogLevel     string `env:"LOG_LEVEL"     env-default:"info"`

	initialized bool
}

var config Config

func GetConfig() Config {
	if config.initialized {
		return config
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Sugar().Fatal(err)
	}
	return config
}
