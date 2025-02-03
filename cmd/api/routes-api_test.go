package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"usual_store/internal/models"

	"github.com/stretchr/testify/assert"
)

func setupTestApp(t *testing.T) *application {
	// Connect to the main PostgreSQL instance
	dsn := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" // Use a default DB to create new ones
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Drop the test database if it already exists
	_, err = db.Exec(`DROP DATABASE IF EXISTS testdb`)
	if err != nil {
		t.Fatalf("Failed to drop existing test database: %v", err)
	}

	// Create a new temporary test database
	_, err = db.Exec(`CREATE DATABASE testdb`)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Close the main connection, we'll connect to the test database
	db.Close()

	// Now connect to the temporary test database
	testDSN := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	db, err = sql.Open("postgres", testDSN)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the necessary tables for testing
	createWidgets(t, db)

	// Return the application with the DBModel initialized
	return &application{DB: models.DBModel{DB: db}}
}

func createWidgets(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS widgets (id              SERIAL PRIMARY KEY,
     name            VARCHAR(255) NOT NULL,
     description     TEXT,
     inventory_level INTEGER,
     price           INTEGER,
    image TEXT, 
    is_recurring BOOL, plan_id  VARCHAR(255),
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO widgets (
    id, name, description, inventory_level, price, image, 
    is_recurring, plan_id, created_at, updated_at
) VALUES (
    1, 'Test Widget', 'Sample description', 100, 9.99, '', 
    true, 123, NOW(), NOW()
);
`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
}

func teardownTestApp(t *testing.T, db *sql.DB) {
	// Close the test database connection first
	err := db.Close()
	if err != nil {
		t.Fatalf("Failed to close the test database connection: %v", err)
	}

	// Now reconnect to the main database (postgres)
	dsn := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to reconnect to the main database: %v", err)
	}

	// Drop the test database
	_, err = db.Exec(`DROP DATABASE IF EXISTS testdb`)
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
	}{
		{"Valid GET route", "GET", "/api/widgets/1", http.StatusOK},
		{"Invalid Route", "GET", "/api/non-existent", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the app with the DB
			app := setupTestApp(t)
			defer func() {
				// Teardown the database after the test
				// This ensures the test database is dropped after each test run
				teardownTestApp(t, app.DB.DB)
			}()

			// Create a new HTTP request and recorder
			handler := app.routes()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rec := httptest.NewRecorder()

			// Serve the HTTP request and capture the response
			handler.ServeHTTP(rec, req)

			// Assert that the expected status code matches the actual response
			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
