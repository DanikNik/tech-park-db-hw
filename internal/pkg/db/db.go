package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
	"sync/atomic"
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

func Truncate() error {
	if dbObj == nil {
		return NotInit
	}

	_, err := Exec("TRUNCATE TABLE tp_forum.users, tp_forum.forum, tp_forum.thread, tp_forum.post, tp_forum.vote CASCADE;")
	if err != nil {
		return err
	}
	Exec("INSERT INTO tp_forum.post (id) VALUES (0)")
	atomic.SwapInt32(&forumCount, 0)
	atomic.SwapInt32(&threadCount, 0)
	atomic.SwapInt32(&postCount, 0)
	atomic.SwapInt32(&userCount, 0)
	return nil
}

var forumCount int32
var threadCount int32
var postCount int32
var userCount int32

func increaseForumCount() {
	atomic.AddInt32(&forumCount, 1)
}

func increaseThreadCount() {
	atomic.AddInt32(&threadCount, 1)
}

func increaseUserCount() {
	atomic.AddInt32(&userCount, 1)
}

func increasePostCount(count int32) {
	atomic.AddInt32(&postCount, count)
}
