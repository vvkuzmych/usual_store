package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	frontend string
}

// application holds all the dependencies for the application
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	version  string
}

func main() {
	// `cfg` is used throughout the application to access these values in a structured manner.
	var cfg config

	// Load environment variables
	smtpPassword := mustGetEnv("SMTP_PASSWORD")
	frontUrl := mustGetEnv("FRONT_URL")
	smtpUser := mustGetEnv("SMTP_USER")
	smtpPort := mustParseEnvInt("SMTP_PORT")
	invoicePort := mustParseEnvInt("INVOICE_PORT")

	// Bind command-line flags to fields in the `cfg` struct.
	// Each flag maps to a specific configuration value, with default values provided if not specified.

	// Define the server port and SMTP configuration
	flag.IntVar(&cfg.port, "port", invoicePort, "Server port to listen on")
	flag.IntVar(&cfg.smtp.port, "smtpport", smtpPort, "SMTP port to listen on")

	// SMTP host and credentials
	flag.StringVar(&cfg.smtp.host, "smtphost", "smtp.mailtrap.io", "SMTP host")
	flag.StringVar(&cfg.smtp.password, "smtppassword", smtpPassword, "SMTP password")
	flag.StringVar(&cfg.smtp.username, "smtpuser", smtpUser, "SMTP user")

	// Application secrets and frontend URL
	flag.StringVar(&cfg.frontend, "frontend", frontUrl, "Frontend URL")

	// Parse the command-line flags and apply their values.
	// This step processes all the flags defined above, overriding default values
	// with those provided in the command line.
	flag.Parse()

	// Set up loggers for info and error logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize application with the repository
	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		version:  version,
	}

	err := app.CreateDirIfNotExist("./invoices")
	if err != nil {
		app.errorLog.Fatal(err)
		return
	}

	err = app.serve()
	if err != nil {
		errorLog.Fatal(err)
	}
}

// Helper to ensure environment variables are set
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is not set in .env file", key)
	}
	return value
}

// mustParseEnvInt retrieves an environment variable and parses it as an integer.
func mustParseEnvInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Environment variable %s is not set\n", key)
		os.Exit(1)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		fmt.Printf("Error converting %s to int: %v\n", key, err)
		os.Exit(1)
	}

	return intValue
}

// serve starts the HTTP server and handles graceful shutdown
func (app *application) serve() error {
	// Create the HTTP server with specific configurations
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	// Log that the server is starting
	app.infoLog.Printf("Starting invoice microservice end server in on port %d", app.config.port)

	// Graceful shutdown handling of invoice
	go func() {
		// Listen for interrupt or termination signals
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		// Wait for a signal
		<-ch
		app.infoLog.Println("Shutting down server gracefully...")

		// Create a context with a timeout for the shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			app.errorLog.Printf("Error during server shutdown: %v", err)
		}
	}()

	// Start the server
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err // Only return unexpected errors
	}

	// Log successful shutdown
	app.infoLog.Println("Server stopped successfully.")
	return nil
}
