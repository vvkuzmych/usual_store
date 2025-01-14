package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func OpenDB(dsn string) (*sql.DB, error) {
	// Open the database connection using the postgres driver
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("got error")
		return nil, err
	}

	// Ping the database to ensure the connection is established with
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Return the connection object
	fmt.Println("Successfully connected to database")
	return db, nil
}
