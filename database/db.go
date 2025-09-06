package database

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
    var err error
    db, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
        log.Fatalf("无法连接数据库: %v", err)
        return err
    }
    return nil
}