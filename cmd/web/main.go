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
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application enviornment {development|production}")
	//flag.StringVar(&cfg.db.dsn, "db-dsn", "root:admin123@tcp(localhost:3306)/widgets?parseTime=true&tls=false", "Database DSN")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:password@database:5432/usualstore?sslmode=disable", "Database DSN")
	flag.StringVar(&cfg.api, "api", "http://localhost:4001", "URL to api")
	flag.StringVar(&cfg.secretkey, "secret", "FrwEtHyuJVJljlkiUNuobKbYPYBknbsw", "secret key")
	flag.StringVar(&cfg.frontend, "frontend", "http://localhost:4000", "url to frontend")

	flag.Parse()

	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer conn.Close()

	// set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = NewPostgresSessionStore(conn)

	tc := make(map[string]*template.Template)

	app := &application{
		config:        cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		version:       version,
		DB: models.DBModel{
			DB: conn,
		},
		Session: session,
	}

	err = app.serve()
	if err != nil {
		app.errorLog.Println(err)
		log.Fatal(err)
	}
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
