package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

//func OpenDB(dsn string) (*sql.DB, error) {
//	db, err := sql.Open("mysql", dsn)
//	if err != nil {
//		return nil, err
//	}
//	if err := db.Ping(); err != nil {
//		return nil, err
//	}
//	return db, nil
//}

func OpenDB(dsn string) (*sql.DB, error) {
	// Open the database connection using the postgres driver
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("got error")
		return nil, err
	}
	fmt.Println("success to open")
	fmt.Println("db+++++", db)
	// Ping the database to ensure the connection is established
	if err := db.Ping(); err != nil {
		fmt.Println("no ping")
		return nil, err
	}
	fmt.Println("success to ping")
	// Return the connection object
	fmt.Println("Successfully connected to database")
	return db, nil
}
