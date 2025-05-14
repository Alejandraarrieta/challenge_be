package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var sqlClient *sql.DB

func initSQLClient() error {
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUsername, dbPassword, dbName,
	)

	newSQLClient, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = newSQLClient.Ping(); err != nil {
		return err
	}

	sqlClient = newSQLClient
	return nil
}

func GetSQLClientInstance() (*sql.DB, error) {
	if sqlClient == nil {
		if err := initSQLClient(); err != nil {
			return nil, err
		}
	}

	return sqlClient, nil
}
