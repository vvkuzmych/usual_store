package main

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
	"usual_store/internal/driver"
	"usual_store/internal/models"
)

const version = "1.0.0"
const cssVersion = "1"

var session *scs.SessionManager

type config struct {
	port int
	env  string
	api  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	secretkey string
	frontend  string
}

type application struct {
	config        config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	version       string
	DB            models.DBModel
	Session       *scs.SessionManager
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", app.config.port),
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	app.infoLog.Println(fmt.Sprintf("Starting HTTP server in %s mode on port %d", app.config.env, app.config.port))

	return srv.ListenAndServe()
}

func main() {
	gob.Register(TransactionData{})

	// Load environment variables
	secretKeyForFront := mustGetEnv("SECRET_FOR_FRONT")
	defaultDSN := mustGetEnv("DATABASE_DSN")
	apiUrl := mustGetEnv("API_URL")
	frontUrl := mustGetEnv("FRONT_URL")

	// Parse command-line flags
	cfg := parseFlags(secretKeyForFront, defaultDSN, apiUrl, frontUrl)

	// Setup loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Connect to the database
	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Close()

	// Setup session management
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = NewPostgresSessionStore(conn)

	// Initialize application
	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: make(map[string]*template.Template),
		version:       version,
		DB: models.DBModel{
			DB: conn,
		},
		Session: session,
	}

	// Start the server
	if err := app.serve(); err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
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

// Parse flags and populate the configuration
func parseFlags(secretKey, dsn, apiUrl, frontUrl string) config {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment {development|production}")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dsn, "Database DSN")
	flag.StringVar(&cfg.secretkey, "secret", secretKey, "Secret key")
	flag.StringVar(&cfg.api, "api", apiUrl, "URL to API")
	flag.StringVar(&cfg.frontend, "frontend", frontUrl, "URL to frontend")
	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	return cfg
}

// PostgresSessionStore is a custom PostgreSQL session store implementation.
type PostgresSessionStore struct {
	db *sql.DB
}

// Find retrieves a session from the PostgreSQL database.
func (store *PostgresSessionStore) Find(token string) (b []byte, found bool, err error) {
	var data []byte
	var expiry time.Time

	// Query the database for the session based on the session token
	err = store.db.QueryRow("SELECT data, expiry FROM sessions WHERE token = $1", token).Scan(&data, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil // Session not found
		}
		return nil, false, err // Database error
	}

	// Check if the session has expired
	if expiry.Before(time.Now()) {
		return nil, false, nil // Session expired
	}

	// Return session data
	return data, true, nil
}

// Commit stores or updates a session in the PostgreSQL database.
func (store *PostgresSessionStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	// Insert or update the session in the database
	_, err = store.db.Exec(
		"INSERT INTO sessions (token, data, expiry) VALUES ($1, $2, $3) "+
			"ON CONFLICT (token) DO UPDATE SET data = $2, expiry = $3",
		token, b, expiry,
	)
	return err
}

// Get retrieves a session from the PostgreSQL database.
func (store *PostgresSessionStore) Get(sessionID string) (map[string]interface{}, error) {
	var data []byte
	var expiry time.Time

	err := store.db.QueryRow("SELECT data, expiry FROM sessions WHERE token = $1", sessionID).Scan(&data, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Session not found
		}
		return nil, err
	}

	if expiry.Before(time.Now()) {
		return nil, nil // Session expired
	}

	// Deserialize session data (you could use JSON, gob, or other formats)
	// Here we assume data is in JSON format
	var sessionData map[string]interface{}
	err = json.Unmarshal(data, &sessionData)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

// Set saves a session to the PostgreSQL database.
func (store *PostgresSessionStore) Set(sessionID string, sessionData map[string]interface{}, expiry time.Time) error {
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	_, err = store.db.Exec("INSERT INTO sessions (token, data, expiry) VALUES ($1, $2, $3) ON CONFLICT (token) DO UPDATE SET data = $2, expiry = $3",
		sessionID, data, expiry)
	return err
}

// Delete removes a session from the PostgreSQL database.
func (store *PostgresSessionStore) Delete(sessionID string) error {
	_, err := store.db.Exec("DELETE FROM sessions WHERE token = $1", sessionID)
	return err
}

// NewPostgresSessionStore creates a new instance of PostgresSessionStore.
func NewPostgresSessionStore(db *sql.DB) *PostgresSessionStore {
	return &PostgresSessionStore{db: db}
}
