package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/manuzokas/subscription-api/internal/adapters/database"
	"github.com/manuzokas/subscription-api/internal/core/subscription"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("Starting Worker...")

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	amqpURL := os.Getenv("RABBITMQ_URL")
	databaseUrl := os.Getenv("DATABASE_URL")
	if amqpURL == "" || databaseUrl == "" {
		log.Fatal("RABBITMQ_URL and DATABASE_URL must be set")
	}

	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()
	subRepo := database.NewPostgresRepository(pool)

	const queueName = "subscription_created_events"
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var event subscription.SubscriptionCreatedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error decoding message: %s", err)
				continue
			}

			log.Printf("Sending welcome email to %s for subscription %s...", event.Email, event.SubscriptionID)
			time.Sleep(3 * time.Second)
			log.Println("Email sent!")

			sub, err := subRepo.FindByID(context.Background(), event.SubscriptionID)
			if err != nil {
				log.Printf("Error finding subscription %s: %s", event.SubscriptionID, err)
				continue
			}

			now := time.Now().UTC()
			sub.Status = "TRIAL"
			sub.TrialEndsAt = &[]time.Time{now.Add(14 * 24 * time.Hour)}[0]
			sub.UpdatedAt = now

			if err := subRepo.Save(context.Background(), sub); err != nil {
				log.Printf("Error updating subscription %s: %s", event.SubscriptionID, err)
				continue
			}
			log.Printf("Subscription %s status updated to TRIAL.", sub.ID)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
