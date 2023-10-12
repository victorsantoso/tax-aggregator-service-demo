package dbconn

import (
	"database/sql"
	"fmt"
	"log"
	"tax-aggregator-service-demo/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

const (
	maxOpenConns    = 60
	connMaxLifetime = 120 * time.Second
	maxIdleConns    = 30
	connMaxIdleTime = 20 * time.Second
	DefaultTimeout  = 2000
)

func NewMySQLDBConn(dbConfig *config.Database) (*sql.DB, error) {
	dbConn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.DBUsername,
		dbConfig.DBPassword,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBName,
	))
	if err != nil {
		log.Fatal("[dbconn.NewMySQLDBConn]:: error opening connection into source mysql")
	}
	setSQLDBConn(dbConn)
	if err = dbConn.Ping(); err != nil {
		log.Fatal("[dbconn.NewMySQLDBConn]:: error pinging connection")
	}
	return dbConn, nil
}

func NewPostgreSQLDBConn(dbConfig *config.Database) (*sql.DB, error) {
	dbConn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBUsername,
		dbConfig.DBPassword,
		dbConfig.DBName,
	))
	setSQLDBConn(dbConn)
	if err != nil {
		log.Println("[dbconn.NewPostgreSQLDBConn]:: error opening connection into source postgresql")
		return nil, err
	}
	return dbConn, nil
}

func setSQLDBConn(sqlConn *sql.DB) {
	sqlConn.SetMaxOpenConns(maxOpenConns)
	sqlConn.SetConnMaxLifetime(connMaxLifetime)
	sqlConn.SetMaxIdleConns(maxIdleConns)
	sqlConn.SetConnMaxIdleTime(connMaxIdleTime)
}

type TransactionOption func(*sql.Tx) error

func WithTransaction(db *sql.DB, fns ...TransactionOption) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println("[dbConn.WithTransaction]:: error starting transaction")
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if err = tx.Rollback(); err != nil {
				log.Println("[dbConn.WithTransaction]:: error recovering, rollback transaction")
			}
		} else if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Println("[dbConn.WithTransaction]:: error on rollback transaction")
			}
		} else {
			if err = tx.Commit(); err != nil {
				log.Println("[dbConn.WithTransaction]:: error commit transaction")
			}
		}
	}()

	for _, fn := range fns {
		if err := fn(tx); err != nil {
			log.Println("[dbConn.WithTransaction]:: error")
			return err
		}
	}
	return nil
}
