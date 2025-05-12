package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string

	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
}

func LoadConfig() Config {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Невозможно преобразовать SMTP_PORT: %v", err)
	}

	return Config{
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),

		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: port,
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
	}
}
