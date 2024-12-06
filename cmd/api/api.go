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
	"syscall"
	"time"
	"usual_store/internal/driver"
	"usual_store/internal/models"
	"usual_store/internal/pkg/repository"
	"usual_store/internal/pkg/service"
)

const version = "1.0.0"

// config holds the configuration values for the application
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	secretkey string
	frontend  string
}

// application holds all the dependencies for the application
type application struct {
	config       config
	infoLog      *log.Logger
	errorLog     *log.Logger
	version      string
	DB           models.DBModel
	tokenService service.TokenService
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
	app.infoLog.Println(fmt.Sprintf("Starting Back end server in %s mode on port %d", app.config.env, app.config.port))

	// Graceful shutdown handling
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

func main() {
	var cfg config

	// Parse command-line flags
	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	//flag.StringVar(&cfg.db.dsn, "db-dsn", "root:admin123@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "Database DSN")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:password@database:5432/usualstore?sslmode=disable", "Database DSN")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production|maintenance}")
	flag.StringVar(&cfg.smtp.host, "smtphost", "smtp.mailtrap.io", "smtp host")
	flag.StringVar(&cfg.smtp.password, "smtppassword", "8d80f34d4bbe3d", "smtp password")
	flag.StringVar(&cfg.smtp.username, "smtpuser", "d43e9bfd6010ba", "smtp user")
	flag.IntVar(&cfg.smtp.port, "smtpport", 587, "smtp port to listen on")
	flag.StringVar(&cfg.secretkey, "secret", "FrwEtHyuJVJljlkiUNuobKbYPYBknbsw", "secret key")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to frontend")

	flag.Parse()

	// Load Stripe keys from environment variables
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	// Set up loggers for info and error logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Channel to handle the database connection and errors
	dbConnectedCh := make(chan *models.DBModel)
	dbErrorCh := make(chan error)

	// Connect to the database in a separate goroutine
	go func() {
		conn, err := driver.OpenDB(cfg.db.dsn)
		if err != nil {
			dbErrorCh <- err
			return
		}
		dbConnectedCh <- &models.DBModel{DB: conn}
	}()

	// Wait for database connection or an error
	select {
	case dbModel := <-dbConnectedCh:
		// Initialize the repository with the database connection
		repo := repository.NewDBModel(dbModel.DB)

		// Initialize application with the repository
		app := &application{
			config:   cfg,
			infoLog:  infoLog,
			errorLog: errorLog,
			version:  version,
			DB: models.DBModel{
				DB: dbModel.DB,
			},
			tokenService: *service.NewTokenService(repo),
		}

		// Start the server
		err := app.serve()
		if err != nil {
			errorLog.Fatal(err)
		}

	case err := <-dbErrorCh:
		// Handle database connection error
		errorLog.Fatal("Error connecting to the database:", err)
	}
}
