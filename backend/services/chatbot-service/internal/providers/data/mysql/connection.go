package mysql

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Connection struct { DB *sql.DB }

func NewConnection(host, port, name, user, password string) (*Connection, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
    db, err := sql.Open("mysql", dsn)
    if err != nil { return nil, err }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    if err := db.Ping(); err != nil { return nil, err }
    return &Connection{DB: db}, nil
}