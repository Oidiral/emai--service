package config

import (
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string `env:"APP_NAME" envDefault:"email-service"`
	GRPCPort        int    `env:"GRPC_PORT" envDefault:"50051"`
	ShutDownTimeout int    `env:"SHUTDOWN_TIMEOUT" envDefault:"10"`
	LogFormat       string `env:"APP_LOG_FORMAT" envDefault:"json"`
	LogLvl          string `env:"APP_LOG_LVL" envDefault:"info"`
	LogQuery        bool   `env:"APP_LOG_QUERY" envDefault:"false"`

	DB struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Pass     string `env:"DB_PASS" envDefault:""`
		Database string `env:"DB_NAME" envDefault:"postgres"`
	}

	RabbitMQ struct {
		Url string `env:"RMQ_URL"`

		Consumer struct {
			Email struct {
				Exchange string `env:"RMQ_EMAIL_EXCHANGE"`
				Send     string `env:"RMQ_EMAIL_ROUTE_SEND"`
			}
		}
	}

	Redis struct {
		Url      string `env:"REDIS_URL"`
		Password string `env:"REDIS_PASSWORD"`
		EmailDB  int    `env:"REDIS_EMAIL_DB" envDefault:"0"`
	}

	Providers struct {
		Email struct {
			Login    string `env:"APP_EMAIL_LOGIN"`
			Pass     string `env:"APP_EMAIL_PASS"`
			SmtpHost string `env:"APP_EMAIL_SMTP_HOST"`
			SmtpPort int    `env:"APP_EMAIL_SMTP_PORT"`
		}
	}

	SendEmailEnabled bool          `env:"APP_SEND_EMAIL_ENABLED" envDefault:"false"`
	MaxWorkers       int           `env:"APP_MAX_WORKERS" envDefault:"10"`
	WorkerTimeout    time.Duration `env:"APP_WORKER_TIMEOUT" envDefault:"10s"`
	MaxRetries       int           `env:"APP_MAX_RETRIES" envDefault:"3"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
