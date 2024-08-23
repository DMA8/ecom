package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ConfigContextKey struct{}

var ConfigKey = ConfigContextKey{}

type Config struct {
	LogLevel                 string `envconfig:"LOG_LEVEL"           default:"debug"`
	ServiceURL               string `envconfig:"SERVICE_URL"         default:":8082"`
	ServiceURLGrpc           string `envconfig:"SERVICE_URL_GRPC"    default:":8084"`
	GrpcNetwork              string `envconfig:"GRPC_NETWORK"        default:"tcp"`
	LomsURLGrpc              string `envconfig:"LOMS_URL_GRPC"       default:"localhost:8083"`
	ProductServiceURL        string `envconfig:"PRODUCT_SERVICE_URL" required:"true"`
	ProductServiceGrpcURL    string `envconfig:"PRODUCT_SERVICE_GRPC_URL" required:"true"`
	ProductServiceRateLimRPS int64  `envconfig:"PRODUCT_SERVICE_RATE_LIM_RPS" default:"10"`
	PSQLConnStr              string `envconfig:"PSQL_CONN_STR" default:"postgres://postgres:our_password@localhost:5432/postgres?sslmode=enable"`
	JaegerURL                string `envconfig:"JAEGER_URL" default:"localhost:6831"`
	JaegerServiceName        string `envconfig:"JAEGER_SERVICE_NAME" default:"cart"`
}

func New() (*Config, error) {
	var cfg Config
	const op = `config.New`

	envFiles := []string{".env"}
	if val := os.Getenv("IS_DOCKER"); val == "" {
		envFiles = append(envFiles, ".env.override")
	}

	if err := godotenv.Overload(envFiles...); err != nil {
		log.Println(err)
	}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("%s err: %w", op, err)
	}

	return &cfg, nil
}
