package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort                string
	JwtAccessTokenSecretKey string
	Redis                   string
	AuthHeaderKey           string
	AuthPayloadKey          string
	SignUpAuthicationTime   time.Duration
	AccessTokenTime         time.Duration
	Postgres                PostgresConfig
	Smtp                    Smtp
}

type Smtp struct {
	Sender   string
	Password string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Load(way string) Config {
	godotenv.Load(way + "/.env")

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		HttpPort:                conf.GetString("HTTP_PORT"),
		Redis:                   conf.GetString("REDIS_ADDR"),
		AuthHeaderKey:           conf.GetString("AUTHORIZATION_HEADER_KEY"),
		AuthPayloadKey:          conf.GetString("AUTHORIZATION_PAYLOAD_KEY"),
		JwtAccessTokenSecretKey: conf.GetString("JWT_ACCESS_TOKEN_SECRET_KEY"),
		SignUpAuthicationTime:   conf.GetDuration("SIGNUP_AUTHENTICATION_DURATION"),
		AccessTokenTime:         conf.GetDuration("ACCESS_TOKEN_DURATION"),
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
		Smtp: Smtp{
			Sender:   conf.GetString("SMTP_SENDER"),
			Password: conf.GetString("SMTP_PASSWORD"),
		},
	}
	return cfg
}
