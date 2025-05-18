package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Server struct {
		Port         string        `env:"SERVER_PORT" envDefault:":8080"`
		ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"5s"`
		WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"10s"`
	}

	Database struct {
		Host        string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port        string `env:"POSTGRES_PORT" envDefault:"5432"`
		User        string `env:"POSTGRES_USER" envDefault:"postgres"`
		Password    string `env:"POSTGRES_PASSWORD" envDefault:"123456"`
		Name        string `env:"POSTGRES_DB" envDefault:"barcode_checker"`
		SSLMode     string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
		MongoDBURI  string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
		MongoDBName string `env:"MONGO_DB" envDefault:"barcode_checker"`
	}

	Auth struct {
		JWTSecret   string `env:"JWT_SECRET" envDefault:"very_secret_key"`
		JWTDuration int    `env:"JWT_DURATION_HOURS" envDefault:"72"`
	}

	BarcodeAPI struct {
		APIKey string `env:"BARCODE_API_KEY"`
	}

	RateLimit struct {
		AuthLimit int `env:"AUTH_RATE_LIMIT" envDefault:"5"`
		AuthBurst int `env:"AUTH_RATE_BURST" envDefault:"5"`
		APILimit  int `env:"API_RATE_LIMIT" envDefault:"30"`
		APIBurst  int `env:"API_RATE_BURST" envDefault:"30"`
	}

	Migration struct {
		Dir string `env:"MIGRATION_DIR" envDefault:"./migrations"`
	}
}

func (c *Config) GetPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
