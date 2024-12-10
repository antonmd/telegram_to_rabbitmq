package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gt.devopspoint.io/devopspoint/telegram-bot/internal/config"
	"gt.devopspoint.io/devopspoint/telegram-bot/internal/queue"
	"gt.devopspoint.io/devopspoint/telegram-bot/internal/telegram"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}
	fmt.Println("Configuration loaded successfully!")

	// Initialize Telegram client (no changes here)
	client := telegram.NewClient(cfg.TelegramBotToken, cfg.TelegramAPIURL)
	fmt.Println("Telegram client initialized!")

	// Set up RabbitMQ connection (no changes to connection logic, just no queue pre-declaration)
	qClient, err := queue.NewQueueClient(cfg.RabbitMQConn)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer qClient.Close()
	log.Println("Connected to RabbitMQ")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = runPollingLoop(ctx, client, qClient)
	if err != nil {
		log.Printf("Polling loop ended with error: %v", err)
	} else {
		log.Println("Polling loop ended gracefully.")
	}

	log.Println("Bot shut down.")
}

func runPollingLoop(ctx context.Context, client *telegram.Client, qClient *queue.QueueClient) error {
	offset := 0
	const longPollingTimeout = 30

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		updates, err := client.GetUpdates(offset, longPollingTimeout)
		if err != nil {
			log.Printf("Error fetching updates: %v", err)
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(5 * time.Second):
			}
			continue
		}

		if len(updates) > 0 {
			for _, update := range updates {
				if update.Message != nil {
					log.Printf("Received message from chat %d: %s",
						update.Message.Chat.ID, update.Message.Text)

					// Publish message to RabbitMQ.
					// CHANGED: We do NOT ensure the queue here; if it doesn't exist, we'll get a return and handle it asynchronously.
					msgPayload := queue.MessagePayload{
						UpdateID: update.UpdateID,
						ChatID:   update.Message.Chat.ID,
						Text:     update.Message.Text,
					}

					err := qClient.PublishMessage(msgPayload)
					if err != nil {
						log.Printf("Failed to publish message: %v", err)
					} else {
						log.Printf("Message publish to chat_%d",
							update.Message.Chat.ID)
					}
				}

				if update.UpdateID >= offset {
					offset = update.UpdateID + 1
				}
			}
		}
	}
}
