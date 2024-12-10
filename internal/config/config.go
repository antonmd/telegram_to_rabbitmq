package config

import (
	"fmt"
	"os"
)

type Config struct {
	TelegramBotToken string
	TelegramAPIURL   string
	RabbitMQConn     string
}

func LoadConfig() (Config, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		return Config{}, fmt.Errorf("missing TELEGRAM_BOT_TOKEN environment variable")
	}

	apiURL := os.Getenv("TELEGRAM_API_URL")
	if apiURL == "" {
		apiURL = "https://api.telegram.org"
	}

	rabbitMQConn := os.Getenv("RABBITMQ_CONN")
	if rabbitMQConn == "" {
		return Config{}, fmt.Errorf("missing RABBITMQ_CONN environment variable")
	}

	return Config{
		TelegramBotToken: token,
		TelegramAPIURL:   apiURL,
		RabbitMQConn:     rabbitMQConn,
	}, nil
}
