package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var Default = map[string]Config{
	"APP_VERSION":          "23.01.311924",
	"APP_URL":              "localhost",
	"APP_ENV":              "production",
	"APP_PORT":             "4000",
	"DB_DRIVER":            "mysql",
	"DB_HOST":              "127.0.0.1",
	"DB_PORT":              "3306",
	"DB_NAME":              "pitra",
	"DB_USER":              "root",
	"DB_PASSWORD":          "",
	"DB_MAX_OPEN_CONNS":    "100",
	"DB_MAX_IDLE_CONNS":    "2",
	"DB_CONN_MAX_LIFETIME": "0ms",
	"DB_IS_DEBUG":          "false",
	"REDIS_HOST":           "127.0.0.1",
	"REDIS_PORT":           "6379",
	"REDIS_PASSWORD":       "",
	"REDIS_DB":             "0",
	"API_IS_DEBUG":         "false",
	"SMTP":                 "smtp-relay.sendinblue.com",
	"SMTP_USER":            "umarpedefi@gmail.com",
	"SMTP_PASSWORD":        "A6kEXjDcf904ZGCd",
	"SLACK_ERROR_URL":      "https://hooks.slack.com/services/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",
}

type Config string

func init() {
	godotenv.Load()
}

func Get(key string) Config {
	value := Config(os.Getenv(key))
	if value == "" {
		value = Default[key]
	}
	return value
}

func (c Config) String() string {
	return string(c)
}

func (c Config) Int() int {
	v, err := strconv.Atoi(c.String())
	if err != nil {
		return 0
	}
	return v
}

func (c Config) Bool() bool {
	return strings.ToLower(c.String()) == "true"
}

func (c Config) Duration() time.Duration {
	v, err := time.ParseDuration(c.String())
	if err != nil {
		return 0
	}
	return v
}

var JWT_KEY = []byte("kjgjksdfkjsdhfkjhsdkjfh23434324")

type JWTClaim struct {
	ID       string
	Email    string
	UserName string
	Company  string
	jwt.RegisteredClaims
}
