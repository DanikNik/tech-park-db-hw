package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var dbObj *pgx.ConnPool

var (
	NotInit     = errors.New("db wasn't initialized")
	AlreadyInit = errors.New("db already initialized")
	ErrNotFound = errors.New("n")
	ErrConflict = errors.New("c")
)

const (
	uniqueIntegrityError = "23505"
	foreignKeyError      = "23503"
	notNullError         = "23502"
)

func Open() (err error) {

	port := 5432
	if os.Getenv("TP_DB_DEVELOPMENT") == "true" {
		port = 32768
	}

	connConfig := pgx.ConnConfig{
		User:              "postgres",
		Password:          "postgres",
		Host:              "localhost",
		Port:              uint16(port),
		Database:          "tech-db-forum",
		TLSConfig:         nil,
		UseFallbackTLS:    false,
		FallbackTLSConfig: nil,
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: 50,
		AcquireTimeout: 10 * time.Second,
		AfterConnect:   nil,
	}

	if dbObj != nil {
		return AlreadyInit
	}
	dbObj, err = pgx.NewConnPool(poolConfig)
	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v\n", err)
	}

	return
}

func Close() {
	dbObj.Close()
}

func QueryRow(query string, args ...interface{}) (*pgx.Row, error) {
	if dbObj == nil {
		return nil, NotInit
	}

	return dbObj.QueryRow(query, args...), nil
}

func Query(query string, args ...interface{}) (*pgx.Rows, error) {
	if dbObj == nil {
		return nil, NotInit
	}

	return dbObj.Query(query, args...)
}

func Exec(query string, args ...interface{}) (pgx.CommandTag, error) {
	return dbObj.Exec(query, args...)
}
