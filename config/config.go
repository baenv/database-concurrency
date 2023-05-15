package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config is the struct to hold the global configuration
type Config struct {
	// DATABBASE
	DBNAME string `envconfig:"POSTGRES_DB"`
	DBUSER string `envconfig:"POSTGRES_USER"`
	DBPASS string `envconfig:"POSTGRES_PASSWORD"`
	DBPORT string `envconfig:"POSTGRES_PORT"`
	DBHOST string `envconfig:"POSTGRES_HOST"`

	// BLOCKCHAIN
	ETHEREUM_RPC string `envconfig:"ETHEREUM_RPC"`

	// QUEUE
	REDIS_HOST string `envconfig:"REDIS_HOST"`
	REDIS_PORT string `envconfig:"REDIS_PORT"`

	// SERVER
	SERVER_PORT       string `envconfig:"SERVER_PORT"`
	ID_GEN_SERVER_URL string `envconfig:"ID_GEN_SERVER_URL"`

	// STREAMING
	STREAM_NAME       string `envconfig:"STREAM_NAME"`
	STREAM_GROUP_NAME string `envconfig:"STREAM_GROUP_NAME"`

	// STREAMING MASTER
	MASTER_PORT          string `envconfig:"MASTER_PORT"`
	MASTER_CONSUMER_NAME string `envconfig:"MASTER_CONSUMER_NAME"`

	// STREAMING CONSUMER
	CONSUMER_NAME string `envconfig:"CONSUMER_NAME"`
	CONSUMER_PORT string `envconfig:"CONSUMER_PORT"`
	MASTER_URL    string `envconfig:"MASTER_URL" required:"true"`

	// ID_GEN_SERVER
	ID_GEN_SERVER_PORT string `envconfig:"ID_GEN_SERVER_PORT"`
}

// NewConfig loads and returns a new Config instance
func NewConfig() (Config, error) {
	var cfg Config

	godotenv.Load(".env")

	return cfg, envconfig.Process("", &cfg)
}
