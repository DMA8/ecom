package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel          string `envconfig:"LOG_LEVEL" default:"debug"`
	ServiceURL        string `envconfig:"SERVICE_URL" default:":8083"`
	ServiceURLGrpc    string `envconfig:"SERVICE_URL_GRPC" default:":8085"`
	GrpcNetwork       string `envconfig:"GRPC_NETWORK" default:"tcp"`
	ProductServiceURL string `envconfig:"PRODUCT_SERVICE_URL" required:"true"`
	PSQLConnStr       string `envconfig:"PSQL_CONN_STR" default:"postgres://postgres:our_password@localhost:5432/postgres?sslmode=enable"`
	KafkaBrokers      string `envconfig:"KAFKA_BROKERS" required:"true"`
}

func New() (*Config, error) {
	var cfg Config
	const op = "config.New"

	envFiles := []string{".env"}
	if val := os.Getenv("IS_DOCKER"); val == "" {
		envFiles = append(envFiles, ".env.override")
	}

	if err := godotenv.Overload(envFiles...); err != nil {
		log.Println(err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &cfg, nil
}
