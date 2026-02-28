package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"usual_store/internal/messaging"
	"usual_store/internal/workerpool"
)

const version = "1.0.0"

type config struct {
	port  int
	env   string
	kafka struct {
		brokers []string
		topic   string
		groupID string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	workers struct {
		count      int
		bufferSize int
	}
}

func main() {
	var cfg config

	// Parse flags
	flag.IntVar(&cfg.port, "port", 8081, "Messaging service port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|production)")

	kafkaBrokers := flag.String("kafka-brokers", getEnv("KAFKA_BROKERS", "localhost:9093"), "Kafka brokers (comma-separated)")
	flag.StringVar(&cfg.kafka.topic, "kafka-topic", getEnv("KAFKA_TOPIC", "email-queue"), "Kafka topic")
	flag.StringVar(&cfg.kafka.groupID, "kafka-group", getEnv("KAFKA_GROUP_ID", "messaging-service-group"), "Kafka consumer group ID")

	flag.StringVar(&cfg.smtp.host, "smtp-host", getEnv("SMTP_HOST", "smtp.mailtrap.io"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", getEnvInt("SMTP_PORT", 2525), "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-user", getEnv("SMTP_USER", ""), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-pass", getEnv("SMTP_PASSWORD", ""), "SMTP password")

	flag.IntVar(&cfg.workers.count, "workers", getEnvInt("EMAIL_WORKER_COUNT", 10), "Number of email workers")
	flag.IntVar(&cfg.workers.bufferSize, "buffer", getEnvInt("EMAIL_WORKER_BUFFER", 100), "Worker job buffer size")

	flag.Parse()

	cfg.kafka.brokers = strings.Split(*kafkaBrokers, ",")

	// Set up loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Starting Messaging Service v%s", version)
	infoLog.Printf("Environment: %s", cfg.env)
	infoLog.Printf("Kafka Brokers: %v", cfg.kafka.brokers)
	infoLog.Printf("Kafka Topic: %s", cfg.kafka.topic)
	infoLog.Printf("Kafka Group ID: %s", cfg.kafka.groupID)
	infoLog.Printf("Worker Pool: %d workers, buffer size: %d", cfg.workers.count, cfg.workers.bufferSize)

	// Create email handler
	emailHandler := NewEmailHandler(cfg.smtp, infoLog, errorLog)

	// Create worker pool
	pool := workerpool.New(cfg.workers.count, cfg.workers.bufferSize, infoLog)
	pool.Start()
	defer pool.Stop()

	// Create Kafka consumer with worker pool
	consumer := messaging.NewPooledConsumer(
		cfg.kafka.brokers,
		cfg.kafka.topic,
		cfg.kafka.groupID,
		emailHandler,
		pool,
		infoLog,
	)
	defer func() {
		if err := consumer.Close(); err != nil {
			errorLog.Printf("Error closing consumer: %v", err)
		}
	}()

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consumer in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			errChan <- err
		}
	}()

	infoLog.Println("Messaging service is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		infoLog.Println("Shutdown signal received, stopping...")
	case err := <-errChan:
		errorLog.Printf("Consumer error: %v", err)
	}

	cancel()

	// Give workers time to finish current jobs
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	infoLog.Println("Waiting for workers to finish...")
	if err := pool.StopWithContext(shutdownCtx); err != nil {
		errorLog.Printf("Worker pool shutdown timeout: %v", err)
	}

	infoLog.Println("Messaging service stopped gracefully")
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an integer environment variable or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}
