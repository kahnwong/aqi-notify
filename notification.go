package main

import (
	"context"

	"github.com/carlmjohnson/requests"
	_ "github.com/joho/godotenv/autoload"
)

type discordWebhook struct {
	Content string `json:"content"`
}

func notify(outputMessage string) error {

	body := discordWebhook{
		Content: outputMessage,
	}

	err := requests.
		URL(discordWebhookUrl).
		BodyJSON(&body).
		Fetch(context.Background())

	return err
}
