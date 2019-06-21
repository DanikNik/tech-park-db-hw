package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

var dbObj *pgx.Conn

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
	connConfig := pgx.ConnConfig{
		User:              "postgres",
		Password:          "postgres",
		Host:              "localhost",
		Port:              32768,
		Database:          "tech-db-forum",
		TLSConfig:         nil,
		UseFallbackTLS:    false,
		FallbackTLSConfig: nil,
	}

	if dbObj != nil {
		return AlreadyInit
	}
	dbObj, err = pgx.Connect(connConfig)
	if err != nil {
		return fmt.Errorf("Unable to establish connection: %v\n", err)
	}

	return
}

func Close() error {
	return dbObj.Close()
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

func isExists(dbName string, tableName string, key string, where string, args ...interface{}) (id interface{}, err error) {
	row, err := findRowBy(dbName, tableName, key, where, args...)
	if err != nil {
		return
	}
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return
	}
	return id, nil
}

func insert(dbName string, tableName string, cols string, values string, args ...interface{}) error {
	_, err := QueryRow("INSERT INTO "+dbName+"."+tableName+" ("+cols+") VALUES ("+values+")", args...)
	return err
}

func findRowBy(dbName string, tableName string, cols string, where string, args ...interface{}) (*pgx.Row, error) {
	if where == "" {
		where = "1"
	}
	return QueryRow("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...)
}

// For future use
//
// func findRowsBy(dbName string, tableName string, cols string, where string, args ...interface{}) (*sql.Rows, error) {
// 	if dbObj == nil {
// 		return nil, NotInit
// 	}
//
// 	if where == "" {
// 		where = "1"
// 	}
// 	return Query("SELECT "+cols+" FROM "+dbName+"."+tableName+" WHERE "+where, args...)
// }

func updateBy(dbName string, tableName string, set string, where string, args ...interface{}) (int64, error) {
	if dbObj == nil {
		return 0, NotInit
	}

	if where == "" {
		where = "1"
	}
	result, err := Exec("UPDATE "+dbName+"."+tableName+" SET "+set+" WHERE "+where, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func removeBy(dbName string, tableName string, where string, args ...interface{}) (int64, error) {
	if dbObj == nil {
		return 0, NotInit
	}

	if where == "" {
		where = "1"
	}
	result, err := Exec("DELETE FROM "+dbName+"."+tableName+" WHERE "+where, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func truncate(dbName string, tableName string) error {
	if dbObj == nil {
		return NotInit
	}

	_, err := Exec("TRUNCATE TABLE " + dbName + "." + tableName)
	if err != nil {
		return err
	}
	return nil
}
