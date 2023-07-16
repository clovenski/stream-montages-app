package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	// postgres
	DBHost         string `env:"SM_POSTGRES_HOST,required"`
	DBUserName     string `env:"SM_POSTGRES_USER,required"`
	DBUserPassword string `env:"SM_POSTGRES_PASSWORD,required"`
	DBName         string `env:"SM_POSTGRES_DB,required"`
	DBPort         string `env:"SM_POSTGRES_PORT,required"`
	Schema         string `env:"SM_POSTGRES_SCHEMA,required"`
	ServerPort     string `env:"SM_SERVER_PORT,required"`

	// kafka
	BootstrapServer string `env:"SM_KAFKA_BOOTSTRAP_SERVER,required"`
	ClientId        string `env:"SM_KAFKA_CLIENT_ID,required"`
	JobsTopic       string `env:"SM_KAFKA_JOBS_TOPIC,required"`
}

func LoadDBConfig() (config Config, err error) {
	log.Println("Loading Server Configurations...")

	err = envdecode.Decode(&config)
	return
}
