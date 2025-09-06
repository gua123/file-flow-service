package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type ConnectionPool struct {
	maxConnections int
	connections    chan *sql.DB
}

var (
	db     *sql.DB
	pool   *ConnectionPool
)

func NewPool(maxConnections int) (*ConnectionPool, error) {
	pool := &ConnectionPool{
		maxConnections: maxConnections,
		connections:    make(chan *sql.DB, maxConnections),
	}

	for i := 0; i < maxConnections; i++ {
		db, err := sql.Open("sqlite3", "./database.db")
		if err != nil {
			return nil, err
		}
		pool.connections <- db
	}

	return pool, nil
}

func InitDB() error {
	var err error
	pool, err = NewPool(10)
	if err != nil {
		return err
	}
	db = <-pool.connections
	return nil
}

func GetConnection() *sql.DB {
	return <-pool.connections
}

func ReleaseConnection(dbConn *sql.DB) {
	pool.connections <- dbConn
}