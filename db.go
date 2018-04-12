package gdb

import (
    "database/sql"
    _ "github.com/lib/pq"
    "errors"
    "sync"
    "fmt"
)

const (
    DRIVER_NAME              = "postgres"
    DATASOURCE_NAME_TEMPLATE = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
)

type dbOptions struct {
    Host     string
    Port     int
    User     string
    Password string
    DataBase string
    SslMode  string
}

func (do *dbOptions) IsNil() bool {
    return do == nil
}

func SetDataSource(host string, port int, user, password, database, sslmode string) {
    client.options = &dbOptions{
        Host:     host,
        Port:     port,
        User:     user,
        Password: password,
        DataBase: database,
        SslMode:  sslmode,
    }
}

type Client struct {
    db      *sql.DB
    err     error
    once    sync.Once
    options *dbOptions
}

var client = &Client{}

var DBOptionsNilErr = errors.New("the options of db is nil")

func loadDB() {
    if client.options.IsNil() {
        client.err = DBOptionsNilErr
        return
    }
    client.db, client.err = sql.Open(DRIVER_NAME, fmt.Sprintf(DATASOURCE_NAME_TEMPLATE, client.options.Host,
        client.options.Port, client.options.User, client.options.Password, client.options.DataBase,
        client.options.SslMode))
}

func DB() (*sql.DB, error) {
    if client.db != nil || client.err != nil {
        return client.db, client.err
    }
    client.once.Do(loadDB)
    return client.db, client.err
}

func query(sqlStr string, args ...interface{}) (*sql.Rows, error) {
    db, err := DB()
    if err != nil {
        return nil, err
    }
    return db.Query(sqlStr, args...)
}
