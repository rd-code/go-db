package gdb

import (
    "database/sql"
    "errors"
    "sync"
    "fmt"
)

const (
    DRIVER_NAME              = "postgres"
    DATASOURCE_NAME_TEMPLATE = "host=%s port=%d user=%s password=%d dbname=%s sslmode=%s"
)

type DBOptions struct {
    Host     string
    Port     int
    User     string
    Password string
    DataBase string
    SslMode  string
}

func (do *DBOptions) IsNil() bool {
    return do == nil
}

func SetOptions(data DBOptions) {
    options = &data
}

var options *DBOptions
var db *sql.DB
var mutex sync.Mutex

var DBOptionsNilErr = errors.New("the options of db is nil")

func DB() (*sql.DB, error) {
    if options.IsNil() {
        return nil, DBOptionsNilErr
    }
    if db != nil {
        return db, nil
    }
    mutex.Lock()
    defer mutex.Unlock()
    if db != nil {
        return db, nil
    }
    var err error
    if db, err = sql.Open(DRIVER_NAME, fmt.Sprintf(DATASOURCE_NAME_TEMPLATE,
        options.Host, options.Port, options.User, options.Password, options.DataBase, options.SslMode)); err != nil {
        return nil, err
    }
    return db, err
}
